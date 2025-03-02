package models

import (
	"sync"
	"time"
)

type TokenBucket struct {
	Capacity int
	Tokens   int
	RefillRate float64
	LastRefillTime time.Time
	Mu sync.Mutex
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		Capacity: capacity,
		Tokens: capacity,
		RefillRate: refillRate,
		LastRefillTime: time.Now(),
		Mu: sync.Mutex{},
	}
}

func (b *TokenBucket) Allow() bool {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	now := time.Now()

	elapsedTime :=  time.Since(b.LastRefillTime).Seconds()
	b.Tokens += int(elapsedTime * b.RefillRate)
	if b.Tokens >= b.Capacity {
		b.Tokens = b.Capacity;
	}
	b.LastRefillTime = now;

	if b.Tokens > 0 {
		b.Tokens--
		return true
	}
	return false
}