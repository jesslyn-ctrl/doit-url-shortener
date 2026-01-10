package url

import (
	"context"
	"time"
)

type Store interface {
	Save(ctx context.Context, u *ShortURL) error
	Get(ctx context.Context, code string) (*ShortURL, error)
	IncrementClick(ctx context.Context, code string, now time.Time) error
}
