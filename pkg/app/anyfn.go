package app

import (
	"reflect"

	"github.com/bingoohuang/gg/pkg/ginx/anyfn"
	"github.com/bingoohuang/gostarter/pkg/ging"
	"github.com/bingoohuang/gostarter/pkg/model"
	"github.com/gin-gonic/gin"
)

func registerWrappers(af *anyfn.Adapter) {
	// 注册如何处理成功返回一个值
	af.PrependOutSupport(anyfn.OutSupportFn(func(v interface{}, vs []interface{}, c *gin.Context) (bool, error) {
		if err, ok := v.(error); ok {
			if err != nil {
				ging.JSON(c, model.Rsp{Status: 500, Message: "error", Data: err.Error()})
			} else if len(vs) == 1 { // 只有一个 error 返回
				ging.JSON(c, model.Rsp{Status: 200, Message: "OK"})
			}
			return true, nil
		}

		if _, ok := vs[0].(anyfn.DirectDealer); ok {
			return false, nil
		}

		if anyfn.IndirectTypeOf(vs[0]) != RspType {
			vs[0] = model.Rsp{Status: 200, Message: "OK", Data: v}
		}

		ging.JSON(c, vs[0])
		return true, nil
	}))

	af.PrependInSupport(anyfn.InSupportFn(func(argIn anyfn.ArgIn, argsIn []anyfn.ArgIn, c *gin.Context) (reflect.Value, error) {
		if argIn.Type == LoginType {
			return anyfn.ConvertPtr(argIn.Ptr, reflect.ValueOf(GetLogin(c))), nil
		}

		return reflect.Value{}, nil
	}))
}

var (
	RspType   = reflect.TypeOf((*model.Rsp)(nil)).Elem()
	LoginType = reflect.TypeOf((*model.Login)(nil)).Elem()
)
