package redis

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var key, val = "redis_test", []string{"test", "redis", "v9"}

func TestRedisClient(t *testing.T) {
	redisEngine, err := connectRedis()
	require.NoError(t, err)
	require.NotEmpty(t, redisEngine)
}

func TestSetGetRedis(t *testing.T) {
	redisEngine, err := connectRedis()
	require.NoError(t, err)
	require.NotEmpty(t, redisEngine)
	err = redisEngine.Set(key, val, 0)
	require.NoError(t, err)
	valByte, check, err := redisEngine.Get(key)
	require.NoError(t, err)
	require.Equal(t, check, true)
	var val2 []string
	err = json.Unmarshal(valByte, &val2)
	require.NoError(t, err)
	require.Equal(t, val2, val)

}
func TestInvalidateRedis(t *testing.T) {
	redisEngine, err := connectRedis()
	require.NoError(t, err)
	require.NotEmpty(t, redisEngine)
	err = redisEngine.Invalidate(key)
	require.NoError(t, err)
}
func connectRedis() (RedisEngine, error) {

	//cfg := configs
	//if err != nil {
	//	return nil, err
	//}
	redisEngine, err := NewRedisClient("1323", "231", "6464")
	if err != nil {
		return nil, err
	}
	return redisEngine, nil
}
