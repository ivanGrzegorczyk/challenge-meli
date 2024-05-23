package main

import (
	"github.com/fufuok/favicon"
	"github.com/gin-gonic/gin"

	"github.com/ivanGrzegorczyk/challenge_meli/handler"
)

func main() {
	r := gin.Default()
	r.Use(favicon.New(favicon.Config{}))

	r.GET("/*path", func(c *gin.Context) {
		handler.ProxyHandler(c)
	})

	r.POST("/rules", func(c *gin.Context) {
		handler.AddRuleHandler(c)
	})

	r.Run(":5000")
}
