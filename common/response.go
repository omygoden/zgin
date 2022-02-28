package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


/**
 * 成功时返回
 */
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

/**
 * 失败时返回
 */
func Error(ctx *gin.Context, code int, message interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  message,
		"data": nil,
	})
}