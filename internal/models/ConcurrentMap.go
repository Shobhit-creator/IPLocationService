package models

import (
	"sync"
)

type ConcurrentMap struct {
	mu sync.RWMutex
	store map[string]int
}