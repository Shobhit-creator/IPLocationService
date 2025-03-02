package interfaces

import (
	"time"

	"github.com/shobhit-Creator/IPLocationService/internal/models"
)

type ProviderSelector interface {
    GetBestProvider(excludeProvider string) (*models.Provider, error)
    UpdateProviderMetrics(name string, isError bool, ResponseTime time.Duration) 
	// UpdateProviderMetrics(provider *models.Provider, isError bool, ResponseTime time.Duration) 
}