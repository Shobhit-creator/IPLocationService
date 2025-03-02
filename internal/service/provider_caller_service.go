package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shobhit-Creator/IPLocationService/internal/models"
	"github.com/shobhit-Creator/IPLocationService/internal/service/interfaces"
)

type ProviderCaller struct {
	ProviderCallerInstance *interfaces.ProviderCaller
}

var _ interfaces.ProviderCaller = &ProviderCaller{}
var (  
	providerCallerInstance *ProviderCaller

)

func GetProviderCallerService() *ProviderCaller {
    return providerCallerInstance;
}

func (pc *ProviderCaller) CallToGetLocation(provider *models.Provider, ip string) (string, error) {
	var url string
 	if provider.Token == "" {
		url = fmt.Sprintf(provider.Url, ip)
	}else {
		url = fmt.Sprintf(provider.Url, ip, provider.Token)
	}
	startTime := time.Now()
	providerSelector := GetProviderSelectorService()

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	endTime := time.Now()

	responseTime := endTime.Sub(startTime)

	data, err := io.ReadAll(response.Body)
	if err != nil {
		providerSelector.UpdateProviderMetrics(provider.Name, true, responseTime);
		return handleFallbackOfCall(pc, provider.Name, ip)
	}
	// this method will be updated when will change request refresh logic for providers
	providerSelector.UpdateProviderMetrics(provider.Name, false, responseTime);
	return string(data), nil
}

func handleFallbackOfCall(pc *ProviderCaller, provider string, ip string) (string, error) {
	providerSelector := GetProviderSelectorService()
	bestProvider, err := providerSelector.GetBestProvider(provider)
	if err != nil {
		return "", err
	}
	return pc.CallToGetLocation(bestProvider, ip)
}