package redisstore

import (
	"context"
	"time"
)

type Store struct {
	client Client
}

func New(client Client) *Store {
	return &Store{client: client}
}

// Get returns the string value for the given key.
func (s *Store) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key)
}

// Set stores a string value with a TTL.
func (s *Store) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl)
}
