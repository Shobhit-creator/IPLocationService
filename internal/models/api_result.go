package models

type ApiResult struct {
	Result string `json:"response"`
	Error  error  `json:"error"`
}