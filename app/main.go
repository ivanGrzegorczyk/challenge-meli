package main

import (
	"github.com/fufuok/favicon"
	"github.com/gin-gonic/gin"

	"github.com/ivanGrzegorczyk/challenge_meli/database"
	"github.com/ivanGrzegorczyk/challenge_meli/handler"
)

func main() {
	rdb := database.NewRedisClient()

	r := gin.Default()
	r.Use(favicon.New(favicon.Config{}))

	r.GET("/*path", func(c *gin.Context) {
		handler.ProxyHandler(c, rdb)
	})

	r.Run(":5000")
}
