package redisstore_test

import (
	"context"
	"testing"
	"time"

	"github.com/eugene-ruby/xconnect/redisstore"
	"github.com/stretchr/testify/require"
)

func TestStore_SetAndGet(t *testing.T) {
	ctx := context.Background()
	mock := redisstore.NewMockClient()
	store := redisstore.New(mock)

	key := "some:test:key"
	value := `{"version":"v1","payload":"abc123"}`

	err := store.Set(ctx, key, value, 5*time.Minute)
	require.NoError(t, err)

	got, err := store.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, got)
}
