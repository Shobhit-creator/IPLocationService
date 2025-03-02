package interfaces

import "github.com/shobhit-Creator/IPLocationService/internal/models"

type ProviderCaller interface {
  CallToGetLocation(provider *models.Provider, ip string) (string, error)
}