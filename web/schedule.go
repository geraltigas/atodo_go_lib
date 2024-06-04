package web

import (
	"atodo_go/schedule"
	"github.com/gin-gonic/gin"
)

func InitScheduleWebInterface(engine *gin.Engine) {
	engine.POST("/schedule", func(c *gin.Context) {
		data, err := schedule.Schedule()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, data)
	})
}
