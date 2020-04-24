package util

import (
	"io"
	"io/ioutil"
	"time"

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
				if p.Latency > time.Minute {
					// Truncate in a golang < 1.8 safe way
					p.Latency -= p.Latency % time.Second
				}

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
