package redisstore_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/eugene-ruby/xconnect/redisstore"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestIntegration_RealRedis_SetAndGet(t *testing.T) {
	ctx := context.Background()

	redisURL := os.Getenv("REDIS_URL")
	require.NotEmpty(t, redisURL, "REDIS_URL must be set")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	defer rdb.Close()

	err := rdb.Ping(ctx).Err()
	require.NoError(t, err)

	store := redisstore.New(redisstore.NewGoRedisAdapter(rdb))

	key := "integration:test:key"
	value := `{"version":"v1","payload":"abc123"}`

	err = store.Set(ctx, key, value, 5*time.Minute)
	require.NoError(t, err)

	got, err := store.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, got)
}
