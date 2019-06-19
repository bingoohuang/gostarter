package model

type Rsp struct {
	Code    int         `json:"status"`
	State   string      `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
