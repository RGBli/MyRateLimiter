package myratelimiter

import (
	"time"
)

// FixedWindowRateLimiter is the base struct for other rate limiter
type FixedWindowRateLimiter struct {
	*BaseRateLimiter
	reqsInWindows int64
}

func NewFixedWindowRateLimiter(limitCount int64, duration int64) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		BaseRateLimiter: NewBaseRateLimiter(limitCount, duration),
		reqsInWindows:   0,
	}
}

func (rl *FixedWindowRateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().Unix()
	rl.lastRequestTime = now

	if now-rl.windowStartTime > rl.duration {
		// update window
		rl.windowStartTime = now
		rl.reqsInWindows = 1
		return true
	}

	if rl.reqsInWindows < rl.limitCount {
		rl.reqsInWindows++
		return true
	}

	return false
}

func (rl *FixedWindowRateLimiter) UpdateLimiter(limitCount int64, duration int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limitCount = limitCount
	rl.duration = duration
}
