package handler

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivanGrzegorczyk/challenge_meli/database"
)

var testIPs = []string{
	"192.168.0.1",
	"192.168.0.2",
	"192.168.0.3",
}

func ProxyHandler(c *gin.Context) {
	//ip := c.ClientIP()
	ip := testIPs[rand.Intn(len(testIPs))]
	path := c.Param("path")

	rules, err := database.GetRulesByIpOrPath(ip, path)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting rules: %s", err)
		return
	}

	for _, rule := range rules {
		if (rule.Ip == ip && rule.Path == path) || (rule.Ip == ip && rule.Path == "") || (rule.Ip == "" && rule.Path == path) {
			count, err := database.GetCounter(ip)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error getting counter: %s", err)
				return
			}

			if count > rule.MaxRequests {
				c.String(http.StatusTooManyRequests, "Too many requests")
				return
			}
		}
	}

	// TODO: Incrementar el contador para ip y path
	for _, rule := range rules {
		if rule.Ip == ip {
			err = database.IncrementCounter(ip, rule.Time)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error incrementing counter: %s", err)
				return
			}
		}
	}

	c.String(http.StatusOK, "OK for ip %s and path %s", ip, path)
}
