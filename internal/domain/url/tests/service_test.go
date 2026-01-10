package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
	_pkgClk "github.com/jesslyn-ctrl/doit-url-shortener/pkg/clock"
)

func TestService_CreateShortURL(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		longURL  string
		ttl      time.Duration
		store    _domainUrl.Store
		wantErr  error
		assertFn func(t *testing.T, u *_domainUrl.ShortURL)
	}{
		{
			name:    "invalid url",
			longURL: "not-a-url",
			ttl:     time.Hour,
			store:   newFakeStore(),
			wantErr: _domainUrl.ErrInvalidURL,
		},
		{
			name:    "success",
			longURL: "https://example.com",
			ttl:     time.Hour,
			store:   newFakeStore(),
			wantErr: nil,
			assertFn: func(t *testing.T, u *_domainUrl.ShortURL) {
				if u.Code == "" {
					t.Fatal("expected code to be generated")
				}
				if u.LongURL != "https://example.com" {
					t.Fatalf("unexpected long url: %s", u.LongURL)
				}
				if !u.ExpiresAt.Equal(start.Add(time.Hour)) {
					t.Fatalf("unexpected expiresAt: %v", u.ExpiresAt)
				}
			},
		},
		{
			name:    "storage get unexpected error",
			longURL: "https://example.com",
			ttl:     time.Hour,
			store:   &errorGetStore{},
			wantErr: errBoom,
		},
		{
			name:    "save error",
			longURL: "https://example.com",
			ttl:     time.Hour,
			store:   &errorSaveStore{},
			wantErr: errBoom,
		},
		{
			name:    "code collision exhausted",
			longURL: "https://example.com",
			ttl:     time.Hour,
			store:   &alwaysExistsStore{},
			wantErr: _domainUrl.ErrCodeCollision,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clk := _pkgClk.NewFakeClock(start)
			svc := _domainUrl.NewService(tt.store, clk)

			u, err := svc.CreateShortURL(context.Background(), tt.longURL, tt.ttl)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.assertFn != nil {
				tt.assertFn(t, u)
			}
		})
	}
}
