package api

type Response struct {
	IsError bool   `json:"is_error"`
	Message string `json:"message"`
	Resp    any    `json:"resp,omitempty"`
}
