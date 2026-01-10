package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
	"github.com/jesslyn-ctrl/doit-url-shortener/internal/storage"
)

func TestMemoryStore_ConcurrentClickIncrement(t *testing.T) {
	store := storage.NewMemoryStore()

	now := time.Now()
	u := &url.ShortURL{
		Code:      "abc123",
		LongURL:   "https://example.com",
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}

	if err := store.Save(context.Background(), u); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	const concurrency = 200
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_ = store.IncrementClick(context.Background(), "abc123", time.Now())
		}()
	}

	wg.Wait()

	result, err := store.Get(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if result.ClickCount != concurrency {
		t.Fatalf("expected click_count=%d, got=%d", concurrency, result.ClickCount)
	}
}
