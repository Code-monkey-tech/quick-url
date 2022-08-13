package cache

import (
	"context"
	"time"
)

type Сaches interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, expiration time.Duration) error
}
