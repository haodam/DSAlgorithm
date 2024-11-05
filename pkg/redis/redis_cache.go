package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(c *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: c,
		ttl:    ttl,
	}
}

func (c *RedisCache) Get(key int) (string, bool) {
	ctx := context.Background()
	intStr := strconv.Itoa(key)
	val, err := c.client.Get(ctx, intStr).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func (c *RedisCache) Set(key int, val string) error {
	ctx := context.Background()
	intStr := strconv.Itoa(key)
	_, err := c.client.Set(ctx, intStr, val, c.ttl).Result()
	return err
}

func (c *RedisCache) Remove(key int) error {
	ctx := context.Background()
	intStr := strconv.Itoa(key)
	_, err := c.client.Del(ctx, intStr).Result()
	return err
}
