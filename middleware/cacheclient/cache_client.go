package serviceCache

import (
	"context"
	"time"
)

type CacheClient interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value string, cost int64, ttl time.Duration) error
}
