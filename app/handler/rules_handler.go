package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ivanGrzegorczyk/challenge_meli/database"
	"github.com/ivanGrzegorczyk/challenge_meli/model"
)

func AddRuleHandler(c *gin.Context) {
	var rule model.Rule
	err := c.BindJSON(&rule)
	if err != nil {
		c.String(400, "Error binding rule: %s", err)
		return
	}

	err = database.InsertRule(rule)
	if err != nil {
		c.String(500, "Error inserting rule: %s", err)
		return
	}

	c.JSON(200, rule)
}
