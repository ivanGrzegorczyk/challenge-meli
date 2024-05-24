package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivanGrzegorczyk/challenge_meli/database"
)

const BASE_URL = "https://api.mercadolibre.com"

func ProxyHandler(c *gin.Context) {
	ip := c.ClientIP()
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
				return
			}

			if count >= rule.MaxRequests {
				status = http.StatusTooManyRequests
				break
			}
		}
	}

	if status == http.StatusTooManyRequests {
		c.String(status, "Too many requests")
		return
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

	url := BASE_URL + path
	log.Println("Haciendo llamada HTTP a:", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error al hacer la llamada HTTP:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al hacer la llamada HTTP"})
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}
