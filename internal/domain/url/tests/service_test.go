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

func TestService_Resolve(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		setupStore  func() _domainUrl.Store
		advanceTime time.Duration
		wantErr     error
		wantClicks  int64
	}{
		{
			name: "success",
			setupStore: func() _domainUrl.Store {
				fs := newFakeStore()
				fs.data["abc"] = &_domainUrl.ShortURL{
					Code:      "abc",
					LongURL:   "https://example.com",
					CreatedAt: start,
					ExpiresAt: start.Add(time.Hour),
				}
				return fs
			},
			wantErr:    nil,
			wantClicks: 1,
		},
		{
			name: "not found",
			setupStore: func() _domainUrl.Store {
				return newFakeStore()
			},
			wantErr: _domainUrl.ErrNotFound,
		},
		{
			name: "expired",
			setupStore: func() _domainUrl.Store {
				fs := newFakeStore()
				fs.data["abc"] = &_domainUrl.ShortURL{
					Code:      "abc",
					LongURL:   "https://example.com",
					CreatedAt: start,
					ExpiresAt: start.Add(time.Minute),
				}
				return fs
			},
			advanceTime: 2 * time.Minute,
			wantErr:     _domainUrl.ErrExpired,
		},
		{
			name: "increment click error",
			setupStore: func() _domainUrl.Store {
				fs := newFakeStore()
				fs.data["abc"] = &_domainUrl.ShortURL{
					Code:      "abc",
					LongURL:   "https://example.com",
					CreatedAt: start,
					ExpiresAt: start.Add(time.Hour),
				}
				return &errorIncrementStore{fs}
			},
			wantErr: errBoom,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clk := _pkgClk.NewFakeClock(start)
			store := tt.setupStore()
			svc := _domainUrl.NewService(store, clk)

			if tt.advanceTime > 0 {
				clk.Advance(tt.advanceTime)
			}

			u, err := svc.Resolve(context.Background(), "abc")

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if u.ClickCount != tt.wantClicks {
				t.Fatalf(
					"expected click_count=%d, got=%d",
					tt.wantClicks,
					u.ClickCount,
				)
			}
		})
	}
}

func TestService_GetStats(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		setupStore  func() _domainUrl.Store
		advanceTime time.Duration
		wantErr     error
		wantClicks  int64
	}{
		{
			name: "success",
			setupStore: func() _domainUrl.Store {
				fs := newFakeStore()
				fs.data["abc"] = &_domainUrl.ShortURL{
					Code:           "abc",
					LongURL:        "https://example.com",
					CreatedAt:      start,
					ExpiresAt:      start.Add(time.Hour),
					ClickCount:     5,
					LastAccessedAt: start.Add(10 * time.Minute),
				}
				return fs
			},
			wantErr:    nil,
			wantClicks: 5,
		},
		{
			name: "not found",
			setupStore: func() _domainUrl.Store {
				return newFakeStore()
			},
			wantErr: _domainUrl.ErrNotFound,
		},
		{
			name: "expired",
			setupStore: func() _domainUrl.Store {
				fs := newFakeStore()
				fs.data["abc"] = &_domainUrl.ShortURL{
					Code:      "abc",
					LongURL:   "https://example.com",
					CreatedAt: start,
					ExpiresAt: start.Add(time.Minute),
				}
				return fs
			},
			advanceTime: 2 * time.Minute,
			wantErr:     _domainUrl.ErrExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clk := _pkgClk.NewFakeClock(start)
			store := tt.setupStore()
			svc := _domainUrl.NewService(store, clk)

			if tt.advanceTime > 0 {
				clk.Advance(tt.advanceTime)
			}

			stats, err := svc.GetStats(context.Background(), "abc")

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if stats.ClickCount != tt.wantClicks {
				t.Fatalf(
					"expected click_count=%d, got=%d",
					tt.wantClicks,
					stats.ClickCount,
				)
			}
		})
	}
}
