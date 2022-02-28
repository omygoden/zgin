package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"zgin/pkg/constant"
	"zgin/pkg/sflogger"
)

//请求日志
func RequestLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != "OPTIONS" {
			params, _ := ioutil.ReadAll(ctx.Request.Body)

			sflogger.Info(constant.LOG_REQUEST, "", map[string]interface{}{
				"【url】":      ctx.Request.RequestURI,
				"【params】":   string(params),
				"【clientIp】": ctx.ClientIP(),
				"【header】":   ctx.Request.Header,
				"【macAddr】":  ctx.GetHeader("Mac-Addr"),
			})
			ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(params))
		}
		ctx.Next()
	}
}
