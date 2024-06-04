package web

import (
	"atodo_go/task_show"
	"github.com/gin-gonic/gin"
)

func InitTaskShowWebInterface(engine *gin.Engine) {
	engine.POST("/task_show/get_show_stack", func(c *gin.Context) {
		stack, err := task_show.GetShowStack()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"stack": stack})
	})

	engine.POST("/task_show/get_show_data", func(c *gin.Context) {
		data, err := task_show.GetShowData()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, data)
	})
}
