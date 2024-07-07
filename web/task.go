package web

import (
	"atodo_go/table"
	"github.com/gin-gonic/gin"
)

type IDRequest struct {
	ID int `json:"id"`
}

type TaskDefaultRequest struct {
	Name       string `json:"name"`
	Goal       string `json:"goal"`
	Deadline   int64  `json:"deadline"`
	InWorkTime bool   `json:"in_work_time"`
}

func InitTaskWebInterface(engine *gin.Engine) {
	engine.POST("/task/eliminate_task", func(c *gin.Context) {
		var request IDRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		err = table.EliminateTask(request.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/task/complete_task", func(c *gin.Context) {
		var request IDRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		err = table.CompleteTask(request.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/task/add_task_default", func(c *gin.Context) {
		var request TaskDefaultRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		id, err := table.CreateTask(request.Name, request.Goal, request.Deadline, request.InWorkTime) // replace with your actual function
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"id": id})
	})

	engine.POST("/task/get_detailed_task", func(c *gin.Context) {
		var request IDRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		task, err := table.GetDetailedTask(request.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, task)
	})

	engine.POST("/task/set_detailed_task", func(c *gin.Context) {
		var taskDetail table.TaskDetail
		err := c.BindJSON(&taskDetail)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		err = table.SetDetailedTask(taskDetail)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/task/update_task_ui", func(c *gin.Context) {
		var request table.UpdateTaskUIs
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		err = table.UpdatePositions(request)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/task/copy_task", func(c *gin.Context) {
		var request IDRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		_, err = table.CopyTask(request.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})
}
