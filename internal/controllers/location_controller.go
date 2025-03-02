package controllers

import (
	"github.com/shobhit-Creator/IPLocationService/internal/models"
	"github.com/shobhit-Creator/IPLocationService/internal/service"
)

func GetLocation(ip string, resultChan chan<- models.ApiResult) {
	providerSelector := service.GetProviderSelectorService()
	provider, err := providerSelector.GetBestProvider("")
	if err != nil {
		resultChan<- models.ApiResult {
			Result: "",
			Error: err,
		}
	}
	providerCaller := service.GetProviderCallerService()
	result, err := providerCaller.CallToGetLocation(provider, ip)
	// Send result back via channel
	resultChan <- models.ApiResult{Result: result, Error: err}
}