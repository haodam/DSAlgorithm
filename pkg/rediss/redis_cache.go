package rediss

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(c *redis.Client) *RedisCache {
	return &RedisCache{
		client: c,
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
	_, err := c.client.Set(ctx, intStr, val, 0).Result()
	return err
}

func (c *RedisCache) Remove(key int) error {
	ctx := context.Background()
	intStr := strconv.Itoa(key)
	_, err := c.client.Del(ctx, intStr).Result()
	return err
}
