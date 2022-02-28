package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/omygoden/gotools/sfslice"
	"log"
	"net/http"
)

var BLACK_IP = []string{
	//"127.0.0.1",
}

//风控规则
func ApiRisk() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//生产环境才校验
		if !ClientIpRule(ctx) {
			ctx.String(http.StatusNotFound, "404 page not found")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

//ip黑名单限制，硬性条件
func ClientIpRule(ctx *gin.Context) bool {
	if sfslice.SliceContainString(BLACK_IP, ctx.ClientIP()) {
		log.Println(fmt.Sprintf("命中黑名单IP：%s", ctx.ClientIP()))
		return false
	}
	return true
}
