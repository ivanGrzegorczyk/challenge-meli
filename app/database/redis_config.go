package database

import (
	"os"

	"github.com/go-redis/redis"
)

func NewRedisClient() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	return redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})
}
