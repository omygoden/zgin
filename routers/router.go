package routers

import (
	"github.com/gin-gonic/gin"
	"zgin/internal/api/test"
	"zgin/internal/middlewares"
)

func InitRouter() *gin.Engine {
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

	return r
}
