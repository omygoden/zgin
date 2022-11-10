package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"zgin/pkg/constant"
	"zgin/pkg/sflogger"
)

type MyResponseWrite struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *MyResponseWrite) Write(s []byte) (int, error) {
	w.body.Write(s)
	return w.ResponseWriter.Write(s)
}

//请求日志
func RequestLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("token", ctx.Request.RemoteAddr)

		w := &MyResponseWrite{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = w

		if ctx.Request.Method != "OPTIONS" {
			params, _ := ioutil.ReadAll(ctx.Request.Body)

			sflogger.Info(constant.LOG_REQUEST, "", map[string]interface{}{
				"【url】":      ctx.Request.RequestURI,
				"【params】":   string(params),
				"【clientIp】": ctx.ClientIP(),
				"【header】":   ctx.Request.Header,
			})
			ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(params))
		}
		ctx.Next()

		log.Println("response:", w.body.String())

	}
}
