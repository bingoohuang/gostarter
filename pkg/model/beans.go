package model

// Rsp defines the response of http request.
type Rsp struct {
	Status  int
	Message string
	Data    interface{}
}

// DemoReq 演示请求 JSON 参数绑定的结构
type DemoReq struct {
	Name string
}

// DemoRsp 演示响应 JSON 参数绑定的结构
type DemoRsp struct {
	Name string
}

type Login struct {
	UserId    uint32
	Username  string
	Role      string
	TimeStamp string
}

const LoginUser = "_loginUser"
