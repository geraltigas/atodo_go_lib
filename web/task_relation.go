package web

import (
	"atodo_go/table"
	"fmt"
	"github.com/gin-gonic/gin"
)

type TaskRelationRequest struct {
	TaskRelation struct {
		Source int `json:"source"`
		Target int `json:"target"`
	} `json:"task_relation"`
}

func InitTaskRelationWebInterface(engine *gin.Engine) {
	engine.POST("/task_relation/add_relation_default", func(c *gin.Context) {
		var request TaskRelationRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := table.AddRelationDefault(request.TaskRelation.Source, request.TaskRelation.Target)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/task_relation/delete_relation", func(c *gin.Context) {
		var request TaskRelationRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		fmt.Println(request)
		err = table.DeleteRelation(request.TaskRelation.Source, request.TaskRelation.Target)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})
}
