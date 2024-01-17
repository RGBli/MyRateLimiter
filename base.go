package myratelimiter

import (
	"sync"
	"time"
)

// BaseRateLimiter is the base struct for other rate limiter
type BaseRateLimiter struct {
	mu              sync.Mutex
	limitCount      int64
	duration        int64
	lastRequestTime int64
	windowStartTime int64
}

func NewBaseRateLimiter(limitCount int64, duration int64) *BaseRateLimiter {
	return &BaseRateLimiter{
		limitCount:      limitCount,
		duration:        duration,
		lastRequestTime: time.Now().Unix(),
		windowStartTime: time.Now().Unix(),
	}
}
