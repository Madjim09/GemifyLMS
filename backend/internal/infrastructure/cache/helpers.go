package cache

// NewDisabledRedisClient returns a no-op RedisClient for testing.
func NewDisabledRedisClient() *RedisClient {
	return &RedisClient{enabled: false}
}
