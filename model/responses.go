package model

type HeartbeatResponse struct {
	Code int `json:"code"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

