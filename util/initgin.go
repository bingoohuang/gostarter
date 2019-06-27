package util

import (
	"github.com/gin-gonic/gin"
	"io"
)

func InitGin(wr io.Writer) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithWriter(wr), RecoveryWithWriter(wr))
	return r
}
