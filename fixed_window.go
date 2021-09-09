package myratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	mu                        sync.Mutex
	rate                      int
	duration                  time.Duration
	lastRequestTime           time.Time
	lastTimeNode              time.Time
	requestsSinceLastTimeNode int
}

func NewFixedWindowRateLimiter(rate int, duration time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		rate:                      rate,
		duration:                  duration,
		lastRequestTime:           time.Now(),
		lastTimeNode:              time.Now(),
		requestsSinceLastTimeNode: 0,
	}
}

func (rl *FixedWindowRateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	lastRequestTime := rl.lastRequestTime
	rl.lastRequestTime = now

	// 请求较少时可以快速返回
	if now.Sub(lastRequestTime) > rl.duration {
		rl.requestsSinceLastTimeNode = 1
		return true
	}

	if now.Sub(rl.lastTimeNode) > rl.duration {
		i := int64(1)
		for rl.lastTimeNode.Add(time.Duration(i * int64(rl.duration))).Before(now) {
			i++
		}
		rl.lastTimeNode = rl.lastTimeNode.Add(time.Duration((i - 1) * int64(rl.duration)))
		rl.requestsSinceLastTimeNode = 1
		return true
	} else {
		if rl.requestsSinceLastTimeNode >= rl.rate {
			return false
		} else {
			rl.requestsSinceLastTimeNode++
			return true
		}
	}
}

func (rl *FixedWindowRateLimiter) UpdateLimiter(rate int, duration time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rate = rate
	rl.duration = duration
}
