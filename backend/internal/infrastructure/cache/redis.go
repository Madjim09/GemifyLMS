package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client  *redis.Client
	enabled bool
}

func NewRedisClient() *RedisClient {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Redis not available, caching disabled: %v", err)
		return &RedisClient{client: client, enabled: false}
	}

	log.Println("Redis connected")
	return &RedisClient{client: client, enabled: true}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	if !r.enabled {
		return "", redis.Nil
	}
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if !r.enabled {
		return nil
	}
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	if !r.enabled {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisClient) InvalidateLeaderboard(ctx context.Context) {
	r.Del(ctx, "leaderboard:top")
}
