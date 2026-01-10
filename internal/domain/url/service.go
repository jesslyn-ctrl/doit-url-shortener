package url

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"
	"time"

	_pkgClk "github.com/jesslyn-ctrl/doit-url-shortener/pkg/clock"
	_logger "github.com/jesslyn-ctrl/doit-url-shortener/pkg/logger"
)

var logs = _logger.GetContextLoggerf(nil)

const (
	shortCodeLength = 7
	maxCodeRetries  = 3
)

type Service struct {
	store Store
	clock _pkgClk.Clock
}

func NewService(store Store, clock _pkgClk.Clock) *Service {
	return &Service{
		store: store,
		clock: clock,
	}
}

// CreateShortURL validates input, generates a unique short code,
// also applies TTL and persists the record
func (s *Service) CreateShortURL(
	ctx context.Context,
	longURL string,
	ttl time.Duration,
) (*ShortURL, error) {
	if err := validateURL(longURL); err != nil {
		return nil, ErrInvalidURL
	}

	now := s.clock.Now()
	expiresAt := now.Add(ttl)

	var (
		code  string
		found bool
		err   error
	)

	// Retry to avoid collisions
	for i := 0; i < maxCodeRetries; i++ {
		var genErr error
		code, genErr = generateShortCode(shortCodeLength)
		if genErr != nil {
			logs.Errorf("error generateShortCode(%d): %v", shortCodeLength, genErr)
			return nil, genErr
		}

		_, err = s.store.Get(ctx, code)
		if errors.Is(err, ErrNotFound) {
			found = true
			break // Got the unique one
		}

		if err != nil {
			logs.Errorf("error get storage: %v", err)
			// Unexpected storage error
			return nil, err
		}
	}

	if !found {
		return nil, ErrCodeCollision
	}

	u := &ShortURL{
		Code:           code,
		LongURL:        longURL,
		CreatedAt:      now,
		ExpiresAt:      expiresAt,
		LastAccessedAt: time.Time{},
	}

	// Store the short URL
	if saveErr := s.store.Save(ctx, u); saveErr != nil {
		logs.Errorf("error saving short url: %v", saveErr)
		return nil, saveErr
	}

	return u, nil
}

// Resolve returns the long URL for a short code
// and update click statistics
func (s *Service) Resolve(
	ctx context.Context,
	code string,
) (*ShortURL, error) {
	u, err := s.store.Get(ctx, code)
	if err != nil {
		logs.Errorf("error resolving short url: %v", err)
		return nil, ErrNotFound
	}

	// Check if already expired
	now := s.clock.Now()
	if u.IsExpired(now) {
		logs.Errorf("short url expired")
		return nil, ErrExpired
	}

	// Exec the click count
	if incErr := s.store.IncrementClick(ctx, code, now); incErr != nil {
		logs.Errorf("error incrementing click count: %v", incErr)
		return nil, incErr
	}

	return u, nil
}

// GetStats returns metadata and statistics for a short URL
func (s *Service) GetStats(
	ctx context.Context,
	code string,
) (*ShortURL, error) {
	u, err := s.store.Get(ctx, code)
	if err != nil {
		logs.Errorf("error getting short url: %v", err)
		return nil, ErrNotFound
	}

	if u.IsExpired(s.clock.Now()) {
		logs.Errorf("short url expired")
		return nil, ErrExpired
	}

	return u, nil
}

/**
+================ HELPERS ================+
*/
// validateURL ensures the URL is HTTP or HTTPS
func validateURL(rawURL string) error {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		logs.Errorf("error parsing URL: %v", err)
		return err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		logs.Errorf("invalid URL scheme: %s", parsedURL.Scheme)
		return ErrInvalidURL
	}

	return nil
}

// Alphabet excludes visually ambigous characters: 0, 0, I, l, and 1
const shortCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"

// generateShortCode generates the short code and avoid specific characters
func generateShortCode(length int) (string, error) {
	b := make([]byte, length)

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(shortCodeAlphabet))))
		if err != nil {
			return "", err
		}
		b[i] = shortCodeAlphabet[n.Int64()]
	}

	return string(b), nil
}
