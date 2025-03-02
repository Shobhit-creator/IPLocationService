package models

import "time"

type ProviderQualityMetrics struct {
	Name            string   `json:"name"`
	RequestCount    int      `json:"request_count"`
	ErrorCount      int      `json:"error_count"`
	AvgResponseTime float64     `json:"avg_response_time"`
	ResponseTimes   []time.Duration `json:"response_times"`
	IsDailyLimitReached bool `json:"is_daily_limit_reached"`
}