package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

func InitAppWebInterface(router *gin.Engine) {
	router.POST("/close", func(c *gin.Context) {
		// set a setTimeOut callback, and close this app in 3s
		time.AfterFunc(3*time.Second, func() {
			os.Exit(0)
		})
		log.Println("App will be closed in 3s")
		c.JSON(200, gin.H{"message": "App will be closed in 3s"})
	})
}
