package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/shobhit-Creator/IPLocationService/internal/models"
	"github.com/shobhit-Creator/IPLocationService/internal/service/interfaces"
	"github.com/spf13/viper"
)
type ProviderSelector struct {
	providers []models.Provider
	providersMetrics []models.ProviderQualityMetrics
	mu sync.RWMutex
}

var _ interfaces.ProviderSelector = &ProviderSelector{}

var (  
	providerSelectorInstance *ProviderSelector
	once sync.Once
)
func GetProviderSelectorService() *ProviderSelector {
    return providerSelectorInstance
}

func InitProviders() {
	once.Do(func() {
		LoadProvidersFromConfig()
	})
}

func LoadProvidersFromConfig(){
	var newProviders []models.Provider
	err := viper.UnmarshalKey("providers", &newProviders)
	if err != nil {
		log.Fatalf("unable to decode into providers, %v", err)
		return
	}

	err = godotenv.Load("../../.env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
		return;
    }

	for i := range newProviders {
			if newProviders[i].Token != "" {
					fmt.Println(os.Getenv(newProviders[i].Token))
					newProviders[i].Token = os.Getenv(newProviders[i].Token)
			}
	}

	var metrics []models.ProviderQualityMetrics
    for _, provider := range newProviders {
		metrics = append(metrics, models.ProviderQualityMetrics{ Name: provider.Name })
    }

	providerSelectorInstance = &ProviderSelector{
		providers: newProviders,
		providersMetrics: metrics,
	}
}

func (ps *ProviderSelector) GetBestProvider(excludeProvider string) (*models.Provider, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	var tempfilteredProviders []models.ProviderQualityMetrics
	for _, providerMetric := range ps.providersMetrics {
		if providerMetric.IsDailyLimitReached  {
			continue
		}
		tempfilteredProviders = append(tempfilteredProviders, providerMetric)
	}

	var filteredProviders []models.ProviderQualityMetrics
	for _, provider := range tempfilteredProviders {
		if provider.Name == excludeProvider  {
			continue
		}
		filteredProviders = append(filteredProviders, provider)
	}
	
	sort.Slice(filteredProviders, func(i, j int) bool {
        if filteredProviders[i].RequestCount != filteredProviders[j].RequestCount {
            return filteredProviders[i].RequestCount < filteredProviders[j].RequestCount
        }
        if filteredProviders[i].ErrorCount != filteredProviders[j].ErrorCount {
            return filteredProviders[i].ErrorCount < filteredProviders[j].ErrorCount
        }
        return filteredProviders[i].AvgResponseTime < filteredProviders[j].AvgResponseTime
    })

	bestProvider := filteredProviders[0]
	for _, provider := range ps.providers {
		if provider.Name == bestProvider.Name {
			return &provider, nil
		}
	}
	return nil, errors.New("No provider found")
}

func (ps *ProviderSelector) UpdateProviderMetrics(name string, isError bool, ResponseTime time.Duration) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var (  
		metrics *models.ProviderQualityMetrics
		currProvider *models.Provider 
	)

    jsonBytes, err := json.Marshal(ps.providersMetrics)
    if err != nil {
        fmt.Println("Error marshalling data:", err)
        return
    }

    // Print the JSON string
    fmt.Println(string(jsonBytes))

	for i, providerMetric := range ps.providersMetrics {
		if providerMetric.Name == name {
			metrics = &ps.providersMetrics[i]
			break
		}
	}

	for i, provider := range ps.providers {
		if provider.Name == name {
			currProvider = &ps.providers[i]
			break
		}
	}

	if(metrics == nil) {
		log.Printf("Provider %s not found in metrics", name)
		return
	}

	metrics.RequestCount++

	if metrics.RequestCount > currProvider.DailyLimit {
		metrics.IsDailyLimitReached = true
	}

	if isError {
		metrics.ErrorCount++
	}

	metrics.ResponseTimes = append(metrics.ResponseTimes, ResponseTime)
	totalResponseTime := time.Duration(0)
	for _, responseTime := range metrics.ResponseTimes {
		totalResponseTime += responseTime
	}
	metrics.AvgResponseTime = float64(totalResponseTime / time.Duration(len(metrics.ResponseTimes)))

    // Marshal the slice to JSON
    jsonBytes, err = json.Marshal(ps.providersMetrics)
    if err != nil {
        fmt.Println("Error marshalling data:", err)
        return
    }

    // Print the JSON string
    fmt.Println(string(jsonBytes))
}
