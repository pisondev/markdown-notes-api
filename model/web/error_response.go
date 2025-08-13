package web

type ErrorResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data"`
}
