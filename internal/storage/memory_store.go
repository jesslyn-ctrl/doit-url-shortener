package storage

import (
	"context"
	"sync"
	"time"

	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
)

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]*_domainUrl.ShortURL
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]*_domainUrl.ShortURL),
	}
}

// Save stores a ShortURL into memory,
// if the code already exists, it will be overwritten.
func (m *MemoryStore) Save(ctx context.Context, u *_domainUrl.ShortURL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Overwrite will be allowed after the uniqueness is checked in domain level
	m.data[u.Code] = u
	return nil
}

// Get retrieves a short URL from memory by its code
func (m *MemoryStore) Get(
	ctx context.Context,
	code string,
) (*_domainUrl.ShortURL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u, ok := m.data[code]
	if !ok {
		return nil, _domainUrl.ErrNotFound
	}

	return u, nil
}

// IncrementClick increments the click count of a short URL
// and updates its last_accessed timestamp
func (m *MemoryStore) IncrementClick(
	ctx context.Context,
	code string,
	now time.Time,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, ok := m.data[code]
	if !ok {
		return _domainUrl.ErrNotFound
	}

	u.ClickCount++
	u.LastAccessedAt = now

	return nil
}
