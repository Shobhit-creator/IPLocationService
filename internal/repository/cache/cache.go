package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/shobhit-Creator/IPLocationService/internal/models"
	"github.com/shobhit-Creator/IPLocationService/internal/repository/cache/interfaces"
)

type MemoryCache struct {
    concurrentMap *models.ConcurrentMap
}


var _ interfaces.Cache = &MemoryCache{}

var (
    instance *MemoryCache
    once     sync.Once
)

func GetInstance() *MemoryCache {
    once.Do(func() {
        instance = &MemoryCache{
            concurrentMap: models.NewConcurrentMap(),
        }
    })
    return instance
}

func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) error {
    c.concurrentMap.Set(key, value, duration)
    return nil
}

func (c *MemoryCache) Get(key string) (interface{}, error) {
    value, exists := c.concurrentMap.Get(key)
    if !exists {
        return "", errors.New("key not found or expired")
	}
    return value, nil
}

func (c *MemoryCache) Delete(key string) error {
    c.concurrentMap.Delete(key)
    return nil
}