package util

import (
	"io"

	"github.com/gin-gonic/gin"
)

func InitGin(wr io.Writer) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithWriter(wr), MakeGinRecovery().RecoveryWithWriter(wr))
	return r
}
