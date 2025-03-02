package models

import "sync"

type RateLimitConfig struct {
	HourlyLimit int
	MinuteLimit int
	Mu          sync.Mutex
}