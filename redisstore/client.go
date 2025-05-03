package redisstore

import (
	"context"
	"time"
)

// Client is a minimal interface for Redis operations.
type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}
