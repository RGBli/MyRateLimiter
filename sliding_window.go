package myratelimiter

import (
	"time"
)

type SlidingWindowRateLimiter struct {
	*FixedWindowRateLimiter
	slidingUnits int64
	unitDuration time.Duration
	requests     []int
}

func NewSlidingWindowRateLimiter(rate int, duration time.Duration, slidingUnits int64) *SlidingWindowRateLimiter {
	unitDuration := time.Duration(duration.Nanoseconds() / slidingUnits)
	requests := make([]int, slidingUnits)
	return &SlidingWindowRateLimiter{
		FixedWindowRateLimiter: NewFixedWindowRateLimiter(rate, duration),
		slidingUnits:           slidingUnits,
		unitDuration:           unitDuration,
		requests:               requests,
	}
}

func (rl *SlidingWindowRateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.window) < rl.duration {
		i := now.Sub(rl.window) / rl.unitDuration
		if sliceSum(rl.requests) < rl.rate {
			rl.requests[i]++
			return true
		}
		return false
	} else if now.Before(rl.window.Add(rl.unitDuration * time.Duration(2*rl.slidingUnits-1))) {
		step := int(now.Sub(rl.window)/rl.unitDuration) + 1
		for i := 0; i < int(rl.slidingUnits)-step; i++ {
			rl.requests[i] = rl.requests[i+step]
		}
		for i := int(rl.slidingUnits) - step; i < len(rl.requests); i++ {
			rl.requests[i] = 0
		}
		rl.window = rl.window.Add(rl.unitDuration * time.Duration(step))
		if sliceSum(rl.requests) < rl.rate {
			rl.requests[len(rl.requests)-1] = 1
			return true
		}
		return false
	} else {
		rl.window = now
		setSliceZero(rl.requests)
		return true
	}
}

func (rl *SlidingWindowRateLimiter) UpdateLimiter(rate int, duration time.Duration, slidingUnits int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rate = rate
	rl.duration = duration
	rl.slidingUnits = slidingUnits
}
