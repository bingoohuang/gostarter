package model

// Rsp defines the response of http request.
type Rsp struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
