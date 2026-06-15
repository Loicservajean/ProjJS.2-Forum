package dto

type ApiError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
