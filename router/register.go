package router

import (
	"go_gin_demo/controller"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.POST("/report", new(controller.BaseController).Report)
	engine.GET("/query", new(controller.BaseController).Query)
}
