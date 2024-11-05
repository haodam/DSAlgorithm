package redis

import (
	"time"

	redisV9 "github.com/redis/go-redis/v9"
)

type RedisEngine interface {
	Configure(...Option) RedisEngine
	Client() *redisV9.Client
	Set(key string, value any, timeToLive ...time.Duration) error
	Get(key string) ([]byte, bool, error)
	Invalidate(key string) error
	InvalidatePrefix(prefix string) error
	Close()
}
