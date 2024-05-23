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

func IncrementCounter(ip string) error {
	key := fmt.Sprintf("ip:%s:%d", ip, time.Now().UnixMicro())
	err := rdb.Set(key, "1", 60*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetCounter(ip string) (int, error) {
	keys, err := rdb.Do("KEYS", fmt.Sprintf("ip:%s:*", ip)).Result()
	if err != nil {
		return 0, err
	}
	return len(keys.([]interface{})), nil
}
