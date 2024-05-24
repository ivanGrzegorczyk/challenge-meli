package database

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})
}

func IncrementCounter(ip string, path string, time_seconds int) error {
	key := fmt.Sprintf("ip:%s:path:%s:%d", ip, path, time.Now().UnixMicro())
	err := rdb.Set(key, "1", time.Duration(time_seconds)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetCounterForIp(ip string) (int, error) {
	keys, err := rdb.Do("KEYS", fmt.Sprintf("ip:%s:*", ip)).Result()
	if err != nil {
		return 0, err
	}
	return len(keys.([]interface{})), nil
}

func GetCounterForPath(path string) (int, error) {
	keys, err := rdb.Do("KEYS", fmt.Sprintf("ip:*:path:%s:*", path)).Result()
	if err != nil {
		return 0, err
	}
	return len(keys.([]interface{})), nil
}

func GetCounterForIpAndPath(ip string, path string) (int, error) {
	keys, err := rdb.Do("KEYS", fmt.Sprintf("ip:%s:path:%s:*", ip, path)).Result()
	if err != nil {
		return 0, err
	}
	return len(keys.([]interface{})), nil
}
