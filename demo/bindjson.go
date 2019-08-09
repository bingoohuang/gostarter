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

func PostBindJSON(ctx *gin.Context, req *ReqBean) {
	ctx.JSON(http.StatusOK, Result{Status: 200, Message: "v2", Data: req})
}

type Result struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	//Metadata Metadata    `json:"metadata"`
	Data interface{} `json:"data"`
}

func WrapBindJSON(reqBody interface{}, handler interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBindJSON(reqBody); err != nil {
			ctx.JSON(http.StatusOK, Result{Status: 400, Message: err.Error()})
			logrus.Errorf("handler %v", err)
			return
		}

		reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(reqBody)})
	}
}
