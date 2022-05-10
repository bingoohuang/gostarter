package controllers

import "github.com/bingoohuang/gostarter/pkg/model"

func init() { registerController("POST", "/demo", demo) }

func demo(m *model.DemoReq) *model.DemoRsp {
	return &model.DemoRsp{Name: "Echo: " + m.Name}
}

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
