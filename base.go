package myratelimiter

import (
	"sync"
	"time"
)

// BaseRateLimiter is the base struct for other rate limiter
type BaseRateLimiter struct {
	mu              sync.Mutex
	limitCount      int64
	duration        time.Duration
	lastRequestTime time.Time
	windowStartTime time.Time
}

func NewBaseRateLimiter(limitCount int64, duration time.Duration) *BaseRateLimiter {
	return &BaseRateLimiter{
		limitCount:      limitCount,
		duration:        duration,
		lastRequestTime: time.Now(),
		windowStartTime: time.Now(),
	}
}
