package util

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Auth(c *gin.Context) {
	basicAuth := viper.GetString("BasicAuth")
	if basicAuth == "" {
		return
	}

	if PassBasicAuth(c, basicAuth) {
		return
	}

	c.Header("WWW-Authenticate", `Basic realm="gostarter Server"`)
	c.AbortWithStatus(http.StatusUnauthorized)
}

func PassBasicAuth(c *gin.Context, basicAuth string) bool {
	authHeader := c.GetHeader("Authorization")
	if strings.Index(authHeader, "Basic ") != 0 {
		return false
	}

	base := authHeader[len("Basic "):]
	userPass, err := base64.StdEncoding.DecodeString(base)
	if err != nil {
		return false
	}

	return string(userPass) == basicAuth
}
