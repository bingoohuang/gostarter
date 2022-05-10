package controllers

import (
	"github.com/bingoohuang/gg/pkg/ginx/adapt"
	"github.com/bingoohuang/gg/pkg/ginx/anyfn"
)

type registeredController struct {
	Method string
	Path   string
	F      interface{}
}

var registeredControllers []registeredController

func registerController(method, path string, f interface{}) {
	registeredControllers = append(registeredControllers, registeredController{
		Method: method,
		Path:   path,
		F:      f,
	})
}

func Register(ar *adapt.Adaptee, af *anyfn.Adapter) {
	for _, c := range registeredControllers {
		ar.Handle(c.Method, c.Path, af.F(c.F))
	}
}
