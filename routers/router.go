package routers

import (
	"github.com/gin-gonic/gin"
	"log"
	"zgin/global"
	"zgin/internal/api/test"
	"zgin/internal/middlewares"
)

func InitRouter() *gin.Engine {
	if !global.Config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestLog())
	r.Use(middlewares.ApiRisk())
	r.Use(middlewares.Recovery())

	api := r.Group("api")
	{
		testApi := api.Group("test").Use(middlewares.ApiCocurrent())
		{
			testApi.GET("test", test.Test1)
			testApi.GET("test2", test.Test2)
		}

	}

	log.Println("路由初始化成功")
	return r
}
