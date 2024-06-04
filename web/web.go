package web

import (
	"github.com/gin-gonic/gin"
	"log"
)

func InitWebInterface() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return nil
	}
	InitAppStateWebInterface(router)
	InitTaskWebInterface(router)
	InitTaskRelationWebInterface(router)
	InitTaskShowWebInterface(router)
	InitScheduleWebInterface(router)
	return router
}

func RunWebServer(router *gin.Engine) {
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
