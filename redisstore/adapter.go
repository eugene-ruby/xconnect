package redisstore

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// goRedisAdapter wraps *redis.Client to match Client interface.
type goRedisAdapter struct {
	client *redis.Client
}

func NewGoRedisAdapter(rdb *redis.Client) Client {
	return &goRedisAdapter{client: rdb}
}

func (g *goRedisAdapter) Get(ctx context.Context, key string) (string, error) {
	val, err := g.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (g *goRedisAdapter) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return g.client.Set(ctx, key, value, ttl).Err()
}
