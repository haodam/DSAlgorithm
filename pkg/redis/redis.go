package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
	
	redisV9 "github.com/redis/go-redis/v9"
)

const (
	_maxRetries      = 5
	_minRetryBackoff = 300 * time.Millisecond
	_maxRetryBackoff = 500 * time.Millisecond
	_dialTimeout     = 5 * time.Second
	_readTimeout     = 5 * time.Second
	_writeTimeout    = 5 * time.Second
	_minIdleConns    = 20
	_poolTimeout     = 6 * time.Second
	_poolSize        = 300
	_database        = 0
)

var ctx = context.Background()

type RedisConnString string

type redis struct {
	password string
	database int
	poolSize int
	rwMutex  sync.Mutex
	client   *redisV9.Client
}

// Get implements RedisEngine.
func (r *redis) Get(key string) ([]byte, bool, error) {
	byteValue, err := r.client.Get(ctx, key).Bytes()
	if err == redisV9.Nil {
		return nil, false, err
	}
	if err != nil {
		return nil, false, err
	}
	return byteValue, true, nil
}

// Invalidate implements RedisEngine.
func (r *redis) Invalidate(key string) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	//Delete Key From Cache
	return r.client.Del(ctx, key).Err()
}

func (r *redis) InvalidatePrefix(prefix string) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	pattern := prefix + "*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Set implements RedisEngine.
func (r *redis) Set(key string, value any, timeToLive ...time.Duration) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()
	byteValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	//set default value 0
	var ttl time.Duration
	if len(timeToLive) > 0 {
		ttl = timeToLive[0]
	}
	//Set value to redis cache
	return r.client.Set(ctx, key, byteValue, ttl).Err()
}

// Close implements RedisEngine.
func (r *redis) Close() {
	r.client.Close()
}

// Client implements RedisEngine.
func (r *redis) Client() *redisV9.Client {
	return r.client
}

// Configure implements RedisEngine.
func (r *redis) Configure(opts ...Option) RedisEngine {
	for _, opt := range opts {
		opt(r)
	}
	return r
}

var _ RedisEngine = (*redis)(nil)

func NewRedisClient(config *configs.Redis) (RedisEngine, error) {
	redis := &redis{
		poolSize: _poolSize,
		database: _database,
		password: config.Password,
	}
	urlRedis := fmt.Sprintf("%s:%s", config.Host, config.Port)
	redis.client = redisV9.NewClient(
		&redisV9.Options{
			Addr:            urlRedis,
			Password:        redis.password,
			DB:              redis.database,
			MaxRetries:      _maxRetries,
			MinRetryBackoff: _minRetryBackoff,
			MaxRetryBackoff: _maxRetryBackoff,
			DialTimeout:     _dialTimeout,
			ReadTimeout:     _readTimeout,
			WriteTimeout:    _writeTimeout,
			MinIdleConns:    _minIdleConns,
			PoolTimeout:     _poolTimeout,
			PoolSize:        redis.poolSize,
		},
	)
	_, err := redis.client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect redis")
	}
	slog.Info("📫 connected to redis 🎉")
	return redis, nil
}
