package main

import (
	"go_gin_demo/cache"
	"go_gin_demo/db"
	"go_gin_demo/router"
	"go_gin_demo/service"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.Use(gin.Recovery())

	router.Register(engine)

	db.InitDB()
	cache.InitRedis()
	cache.InitLocalCache()

	service.InitMsgConsumer()

	engine.Run(":8080")
}
