package utils

import (
	"sync"

	"github.com/shobhit-Creator/IPLocationService/internal/models"
)

func GetOrCreateBucket(buckets *sync.Map, key string, capacity int, refillrate float64) *models.TokenBucket {
	bucket, ok := buckets.Load(key)
	if !ok {
		buckets.Store(key, models.NewTokenBucket(capacity, refillrate))
	}
	return bucket.(*models.TokenBucket)
}