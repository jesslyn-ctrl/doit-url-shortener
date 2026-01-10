package tests

import (
	"context"
	"errors"
	"time"

	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
)

var errBoom = errors.New("boom")

// Base fake (mock) store (success path)
type fakeStore struct {
	data map[string]*_domainUrl.ShortURL
}

func newFakeStore() *fakeStore {
	return &fakeStore{data: make(map[string]*_domainUrl.ShortURL)}
}

func (f *fakeStore) Get(_ context.Context, code string) (*_domainUrl.ShortURL, error) {
	u, ok := f.data[code]
	if !ok {
		return nil, _domainUrl.ErrNotFound
	}
	return u, nil
}

func (f *fakeStore) Save(_ context.Context, u *_domainUrl.ShortURL) error {
	f.data[u.Code] = u
	return nil
}

func (f *fakeStore) IncrementClick(
	_ context.Context,
	code string,
	at time.Time,
) error {
	u, ok := f.data[code]
	if !ok {
		return _domainUrl.ErrNotFound
	}
	u.ClickCount++
	u.LastAccessedAt = at
	return nil
}

// Store that always errors on Get
type errorGetStore struct{}

func (e *errorGetStore) Get(context.Context, string) (*_domainUrl.ShortURL, error) {
	return nil, errBoom
}
func (e *errorGetStore) Save(context.Context, *_domainUrl.ShortURL) error { return nil }
func (e *errorGetStore) IncrementClick(context.Context, string, time.Time) error {
	return nil
}

// Store that errors on Save
type errorSaveStore struct{ fakeStore }

func (e *errorSaveStore) Save(context.Context, *_domainUrl.ShortURL) error {
	return errBoom
}

// Store that always causes collisions
type alwaysExistsStore struct{}

func (a *alwaysExistsStore) Get(context.Context, string) (*_domainUrl.ShortURL, error) {
	return &_domainUrl.ShortURL{}, nil
}
func (a *alwaysExistsStore) Save(context.Context, *_domainUrl.ShortURL) error { return nil }
func (a *alwaysExistsStore) IncrementClick(context.Context, string, time.Time) error {
	return nil
}
