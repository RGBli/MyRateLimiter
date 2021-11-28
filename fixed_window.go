package myratelimiter

import (
	"sync"
	"time"
)

// FixedWindowRateLimiter is the base struct for other rate limiter
type FixedWindowRateLimiter struct {
	mu              sync.Mutex
	rate            int
	duration        time.Duration
	lastRequestTime time.Time
	// left edge of window
	window            time.Time
	requestsInWindows int
}

func NewFixedWindowRateLimiter(rate int, duration time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		rate:              rate,
		duration:          duration,
		lastRequestTime:   time.Now(),
		window:            time.Now(),
		requestsInWindows: 0,
	}
}

func (rl *FixedWindowRateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	lastRequestTime := rl.lastRequestTime
	rl.lastRequestTime = now

	// fast return
	if now.Sub(lastRequestTime) > rl.duration {
		rl.requestsInWindows = 1
		return true
	}

	if now.Sub(rl.window) <= rl.duration {
		if rl.requestsInWindows >= rl.rate {
			return false
		} else {
			rl.requestsInWindows++
			return true
		}
	} else {
		// update window
		rl.window = now
		rl.requestsInWindows = 1
		return true
	}
}

func (rl *FixedWindowRateLimiter) UpdateLimiter(rate int, duration time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rate = rate
	rl.duration = duration
}
