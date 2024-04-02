package ks2

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	RequestId string  `json:"requestId,omitempty"`
	Errors    []Error `json:"errors,omitempty"`
}
