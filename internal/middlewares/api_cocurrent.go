package middlewares

import (
	"github.com/gin-gonic/gin"
)

//并发控制
func ApiCocurrent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//可以根据上下文携带的token,针对用户级进行限制
		//redisKey := fmt.Sprintf("%s:%s", ctx.Request.RequestURI, ctx.Value("token"))
		//redisclient.SetNxWait(redisKey, 1, time.Second*10)
		//defer redisclient.Del(redisKey)
		ctx.Next()
	}
}
