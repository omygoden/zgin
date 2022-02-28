package routers

import (
	"github.com/gin-gonic/gin"
	"zgin/internal/api/test"
	"zgin/internal/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLog())
	r.Use(middleware.ApiRisk())
	r.Use(middleware.Recovery())

	api := r.Group("api")
	{
		testApi := api.Group("test")
		{
			testApi.POST("test", test.Test1)
		}
	}

	return r
}
