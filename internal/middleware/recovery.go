package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"runtime/debug"
	"zgin/common"
	"zgin/pkg/constant"
	"zgin/pkg/errcode"
	"zgin/pkg/sflogger"
)

//异常处理
func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if res, ok := err.(errcode.ApiPanic); ok {
					log.Println("错误信息：", err)
					errMsg := res.Msg
					switch res.Code {
					case errcode.ERROR_CODE_REDIS:
						errMsg = redisErrorHandle(res, ctx)
					case errcode.ERROR_CODE_MYSQL:
						errMsg = mysqlErrorHandle(res, ctx)
					default:
						requestErrorHandle(res, ctx)
					}
					common.Error(ctx, res.Code, errMsg)
				} else {
					serverErrorHandle(err, ctx)
					common.Error(ctx, 500, "系统错误，请联系客服")
				}
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

//请求错误处理
func requestErrorHandle(res errcode.ApiPanic, ctx *gin.Context) {
	sflogger.Error(constant.LOG_REQUEST_ERR, "错误信息", map[string]interface{}{
		"【code】":     res.Code,
		"【msg】":      res.Msg,
		"【url】":      ctx.Request.RequestURI,
		"【params】":   fmt.Sprintf("%+v", res.Data),
		"【clientIp】": ctx.ClientIP(),
	})
}

//服务异常处理
func serverErrorHandle(err interface{}, ctx *gin.Context) {
	m := debug.Stack()
	sflogger.Error(constant.LOG_PANIC, "异常信息", map[string]interface{}{
		"【错误信息】":     err,
		"【异常信息】":     string(m),
		"【url】":      ctx.Request.RequestURI,
		"【params】":   "",
		"【clientIp】": ctx.ClientIP(),
	})
}

//redis服务器异常处理
func redisErrorHandle(res errcode.ApiPanic, ctx *gin.Context) string {
	sflogger.Error(constant.LOG_REDIS_ERR, "错误信息", map[string]interface{}{
		"【code】":     res.Code,
		"【msg】":      res.Msg,
		"【url】":      ctx.Request.RequestURI,
		"【params】":   fmt.Sprintf("%+v", res.Data),
		"【clientIp】": ctx.ClientIP(),
	})
	return "系统错误，请联系客服"
}

//mysql服务器异常处理
func mysqlErrorHandle(res errcode.ApiPanic, ctx *gin.Context) string {
	sflogger.Error(constant.LOG_MYSQL_ERR, "错误信息", map[string]interface{}{
		"【code】":     res.Code,
		"【msg】":      res.Msg,
		"【url】":      ctx.Request.RequestURI,
		"【params】":   fmt.Sprintf("%+v", res.Data),
		"【clientIp】": ctx.ClientIP(),
	})
	return "系统错误，请联系客服"
}
