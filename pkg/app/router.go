package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bingoohuang/gg/pkg/ginx/adapt"
	"github.com/bingoohuang/gg/pkg/ginx/anyfn"
	"github.com/bingoohuang/golog/pkg/ginlogrus"
	"github.com/bingoohuang/gostarter/pkg/conf"

	"github.com/bingoohuang/gg/pkg/logx"
	"github.com/bingoohuang/gostarter/pkg/ging"
	"github.com/bingoohuang/gostarter/pkg/model"
	"github.com/bingoohuang/gostarter/pkg/ui"

	"github.com/gin-gonic/gin"
)

var startupTime = time.Now().Format(`2006-01-02 15:03:04`)

// StartWeb start the web
func StartWeb() {
	r := gin.New()
	r.Use(ginlogrus.Logger(nil, true), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		ging.JSON(c, model.Rsp{Status: http.StatusOK, Message: "started time:" + startupTime})
	})
	r.GET("/ping", func(c *gin.Context) {
		ging.JSON(c, model.Rsp{Status: http.StatusOK, Message: "pong"})
	})
	r.GET("/stats", func(c *gin.Context) {
		ging.JSON(c, struct {
			StartupTime string
			Count       int
		}{
			StartupTime: startupTime,
			Count:       100,
		})
	})

	ar := adapt.Adapt(r)
	af := anyfn.NewAdapter()
	ar.RegisterAdapter(af)

	ar.POST("/demo", af.F(func(m *model.DemoReq) *model.DemoRsp {
		return &model.DemoRsp{Name: "Echo: " + m.Name}
	}))
	/*
		$ gurl POST :1235/demo name=@name
		POST /demo HTTP/1.1
		Host: 127.0.0.1:1235
		Accept: application/json
		Accept-Encoding: gzip, deflate
		Content-Type: application/json
		Gurl-Date: Tue, 10 May 2022 06:20:41 GMT
		User-Agent: gurl/1.0.0

		{
		  "name": "Fairyink"
		}

		HTTP/1.1 200 OK
		Content-Type: application/json; charset=utf-8
		Date: Tue, 10 May 2022 06:20:41 GMT
		Content-Length: 25

		{
		  "name": "Echo: Fairyink"
		}
	*/

	g := r.Group("/", ui.Auth)
	g.GET("/", ui.Handler(""))

	addr := fmt.Sprintf(":%d", conf.Conf.Port)
	log.Printf("start to listen on %s", addr)
	err := r.Run(addr)
	logx.Fatalf(err, "run failed: %v", err)
}
