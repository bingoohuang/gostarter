package model

// Rsp defines the response of http request.
type Rsp struct {
	Status  int
	Message string
	Data    interface{}
}
