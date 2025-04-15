package dto

type Response struct {
	StatusCode uint   `json:"status_code"`
	Message    string `json:"message"`
}

var OKStatus = Response {
	StatusCode: 200,
	Message: "OK",
}
