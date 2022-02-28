package test

import (
	"github.com/gin-gonic/gin"
	"zgin/common"
)

func Test1(ctx *gin.Context) {
	common.Success(ctx, "success")
}
