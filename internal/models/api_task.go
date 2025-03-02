package models

type ApiTask struct {
	Provider *Provider
	Ip       string
	Result   chan<- ApiResult
}