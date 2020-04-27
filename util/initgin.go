package util

import (
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// InitGin initialized the gin.
func InitGin(wr io.Writer) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithConfig(
		gin.LoggerConfig{
			Output: ioutil.Discard,
			Formatter: func(p gin.LogFormatterParams) string {
				logrus.Debugf("[GIN] %3d| %13v | %15s |%-7s %s %s",
					p.StatusCode,
					p.Latency,
					p.ClientIP,
					p.Method,
					p.Path,
					p.ErrorMessage,
				)

				return ""
			},
		}),

		MakeGinRecovery().RecoveryWithWriter(wr),
	)

	return r
}
