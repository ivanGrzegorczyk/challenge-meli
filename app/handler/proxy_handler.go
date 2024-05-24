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

	status := http.StatusOK
	for _, rule := range rules {
		// Si hay alguna regla definida reviso los contadores
		if (rule.Ip == "" && rule.Path == path) || (rule.Ip == ip && rule.Path == "") || (rule.Ip == ip && rule.Path == path) {
			var count int
			var err error

			if rule.Ip == ip && rule.Path == path {
				count, err = database.GetCounterForIpAndPath(ip, path)
			} else if rule.Path == path {
				count, err = database.GetCounterForPath(path)
			} else if rule.Ip == ip {
				count, err = database.GetCounterForIp(ip)
			}

			if err != nil {
				c.String(http.StatusInternalServerError, "Error getting counter: %s", err)
				break
			}

			if count >= rule.MaxRequests {
				status = http.StatusTooManyRequests
				break
			}
		}
	}

	for _, rule := range rules {
		// Si alguna regla aplica aumento el contador, la idea es evitar escribir en la base de datos si no es necesario
		if rule.Ip == ip || rule.Path == path {
			err = database.IncrementCounter(ip, path, rule.Time)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error incrementing counter: %s", err)
				return
			}
			break
		}
	}

	switch status {
	case http.StatusOK:
		c.JSON(http.StatusOK, gin.H{
			"ip":   ip,
			"path": path,
		})
	case http.StatusTooManyRequests:
		c.String(http.StatusTooManyRequests, "Too many requests")
	}
}
