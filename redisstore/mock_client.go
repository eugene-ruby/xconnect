package redisstore

import (
	"context"
	"sync"
	"time"
)

type MockClient struct {
	mu    sync.Mutex
	store map[string]string
}

func NewMockClient() *MockClient {
	return &MockClient{store: make(map[string]string)}
}

func (m *MockClient) Get(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.store[key], nil
}

func (m *MockClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = value
	return nil
}
