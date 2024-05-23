package handler

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var testIPs = []string{
	"192.168.1.1",
	"192.168.1.2",
	"192.168.1.3",
}

func ProxyHandler(c *gin.Context, rdb *redis.Client) {
	//ip := c.ClientIP()
	ip := testIPs[rand.Intn(len(testIPs))]

	err := incrementCounter(rdb, ip)
	if err != nil {
		c.String(500, "Error incrementing counter: %s", err)
		return
	}

	count, err := getCounter(rdb, ip)
	if err != nil {
		c.String(500, "Error getting counter: %s", err)
		return
	}

	c.String(200, "Requests in last minute from ip: %s: %d", ip, count)
}

func incrementCounter(rdb *redis.Client, ip string) error {
	key := fmt.Sprintf("ip:%s:%d", ip, time.Now().UnixMicro())
	err := rdb.Set(key, "1", 60*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func getCounter(rdb *redis.Client, ip string) (int, error) {
	keys, err := rdb.Do("KEYS", fmt.Sprintf("ip:%s:*", ip)).Result()
	if err != nil {
		return 0, err
	}
	return len(keys.([]interface{})), nil
}
