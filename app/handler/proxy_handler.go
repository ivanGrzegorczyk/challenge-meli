package handler

import (
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/ivanGrzegorczyk/challenge_meli/database"
)

var testIPs = []string{
	"192.168.1.1",
	"192.168.1.2",
	"192.168.1.3",
}

func ProxyHandler(c *gin.Context) {
	//ip := c.ClientIP()
	ip := testIPs[rand.Intn(len(testIPs))]

	err := database.IncrementCounter(ip)
	if err != nil {
		c.String(500, "Error incrementing counter: %s", err)
		return
	}

	count, err := database.GetCounter(ip)
	if err != nil {
		c.String(500, "Error getting counter: %s", err)
		return
	}

	c.String(200, "Requests in last minute from ip: %s: %d", ip, count)
}
