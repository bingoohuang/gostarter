package demo

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ReqBean wraps the request bean.
type ReqBean struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// WrapBindJSONRouter Wrap函数，和处理函数，放在一起，对照出现，
func WrapBindJSONRouter() gin.HandlerFunc { return WrapBindJSONImpl(PostBindJSON, PostBindJSONRouter) }

// PostBindJSONRouter 演示POST函数，预解析请求体到第二个参数
func PostBindJSONRouter(_ *gin.Context) interface{} {
	return &ReqBean{}
}

// WrapBindJSON Wrap函数，和处理函数，放在一起，对照出现，
func WrapBindJSON() gin.HandlerFunc { return WrapBindJSONImpl(PostBindJSON, &ReqBean{}) }

// PostBindJSON 演示POST函数，预解析请求体到第二个参数
func PostBindJSON(ctx *gin.Context, req *ReqBean) {
	ctx.JSON(http.StatusOK, Result{Status: 200, Message: "v2", Data: req})
}

// Result wraps results of REST.
type Result struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// WrapBindJSONImpl 包装请求体JSON解析
func WrapBindJSONImpl(handler interface{}, req interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqv := reflect.ValueOf(req)
		if reqv.Kind() == reflect.Func {
			router := req.(func(*gin.Context) interface{})
			req = router(ctx)
			reqv = reflect.ValueOf(req)
		}

		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.JSON(http.StatusOK, Result{Status: 400, Message: err.Error()})
			logrus.Errorf("handler %v", err)
			return
		}

		reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(ctx), reqv})
	}
}
