package web

import (
	"atodo_go/table"
	"github.com/gin-gonic/gin"
)

type WorkTimeRequest struct {
	WorkTime int64 `json:"work_time"`
}

type NowIsWorkTimeRequest struct {
	NowIsWorkTime bool `json:"now_is_work_time"`
}

func InitAppStateWebInterface(engine *gin.Engine) {
	engine.POST("/app_state/get_now_viewing_task", func(c *gin.Context) {
		task, err := table.GetNowViewingTask()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"task": task})
	})

	engine.POST("/app_state/set_now_viewing_task", func(c *gin.Context) {
		var request IDRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		err = table.SetNowViewingTask(request.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/app_state/back_to_parent_task", func(c *gin.Context) {
		err := table.BackToParentTask()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/app_state/set_work_time", func(c *gin.Context) {
		var request WorkTimeRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		err = table.SetWorkTime(request.WorkTime)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/app_state/get_work_time", func(c *gin.Context) {
		workTime, err := table.GetWorkTime()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"work_time": workTime})
	})

	engine.POST("/app_state/set_now_doing_task", func(c *gin.Context) {
		var request IDRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		err = table.SetNowDoingTask(request.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/app_state/get_now_doing_task", func(c *gin.Context) {
		task, err := table.GetNowDoingTask()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"task": task})
	})

	engine.POST("/app_state/set_now_is_work_time", func(c *gin.Context) {
		var request NowIsWorkTimeRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}
		err = table.SetNowIsWorkTime(request.NowIsWorkTime)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.POST("/app_state/get_now_is_work_time", func(c *gin.Context) {
		nowIsWorkTime, err := table.GetNowIsWorkTime()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal error: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"now_is_work_time": nowIsWorkTime})
	})
}
