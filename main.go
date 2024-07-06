package main

import (
	"go_gin_demo/router"
	"go_gin_demo/service"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use(gin.Recovery())

	router.Register(engine)
	service.InitMsgConsumer()

	engine.Run(":8080")
}
