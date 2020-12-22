package util

import (
	"github.com/bingoohuang/golog"
	"github.com/bingoohuang/golog/pkg/ginlogrus"
	"github.com/gin-gonic/gin"
)

// InitGin initialized the gin.
func InitGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	golog.SetupLogrus(nil, "", "")

	r := gin.New()
	r.Use(ginlogrus.Logger(nil, true), gin.Recovery())

	return r
}
