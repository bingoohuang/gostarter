package demo

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ReqBean struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Wrap函数，和处理函数，放在一起，对照出现，
func WrapBindJSON() gin.HandlerFunc { return WrapBindJSONImpl(PostBindJSON, &ReqBean{}) }

// PostBindJSON 演示POST函数，预解析请求体到第二个参数
func PostBindJSON(ctx *gin.Context, req *ReqBean) {
	ctx.JSON(http.StatusOK, Result{Status: 200, Message: "v2", Data: req})
}

type Result struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	//Metadata Metadata    `json:"metadata"`
	Data interface{} `json:"data"`
}

// WrapBindJSONImpl 包装请求体JSON解析
func WrapBindJSONImpl(handler interface{}, req interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.JSON(http.StatusOK, Result{Status: 400, Message: err.Error()})
			logrus.Errorf("handler %v", err)
			return
		}

		reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
	}
}
