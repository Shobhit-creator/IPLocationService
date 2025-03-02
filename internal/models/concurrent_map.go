package models

import (
	"sync"
	"time"
)

type CacheEntry struct {
    Value      interface{}
    Expiration int64
}

type ConcurrentMap struct {
    Mu            sync.RWMutex
    Store         map[string]CacheEntry
    CleanupTicker *time.Ticker
}

func NewConcurrentMap() *ConcurrentMap {
    cm := &ConcurrentMap{
        Store: make(map[string]CacheEntry),
        CleanupTicker: time.NewTicker(1 * time.Minute), // Cleanup interval
    }
    go cm.cleanupExpiredEntries()
    return cm
}

func (cm *ConcurrentMap) Set(key string, value interface{}, duration time.Duration) {
    cm.Mu.Lock()
    defer cm.Mu.Unlock()
    expiration := time.Now().Add(duration).UnixNano()
    cm.Store[key] = CacheEntry{
        Value:      value,
        Expiration: expiration,
    }
}

func (cm *ConcurrentMap) Get(key string) (interface{}, bool) {
    cm.Mu.RLock()
    defer cm.Mu.RUnlock()
    entry, exists := cm.Store[key]
    if !exists || time.Now().UnixNano() > entry.Expiration {
        return nil, false
    }
    return entry.Value, true
}

func (cm *ConcurrentMap) Delete(key string) {
    cm.Mu.Lock()
    defer cm.Mu.Unlock()
    delete(cm.Store, key)
}

func (cm *ConcurrentMap) cleanupExpiredEntries() {
    for range cm.CleanupTicker.C {
        cm.Mu.Lock()
        now := time.Now().UnixNano()
        for key, entry := range cm.Store {
            if now > entry.Expiration {
                delete(cm.Store, key)
            }
        }
        cm.Mu.Unlock()
    }
}