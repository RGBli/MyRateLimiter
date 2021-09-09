package myratelimiter

import (
	"math"
	"time"
)

// 滑动窗口限流器
type SlidingWindowRateLimiter struct {
	*FixedWindowRateLimiter
	slidingUnit time.Duration
	// requests 和 lastTimeNode 构成了滑动窗的基本单元
	requests []int
}

func NewSlidingWindowRateLimiter(rate int, duration time.Duration, slidingUnit time.Duration) *SlidingWindowRateLimiter {
	if slidingUnit > duration {
		panic("Duration cannot be less than slidingUnit")
	}
	requests := make([]int, int(math.Ceil(float64(duration.Nanoseconds()/slidingUnit.Nanoseconds()))))
	return &SlidingWindowRateLimiter{
		FixedWindowRateLimiter: NewFixedWindowRateLimiter(rate, duration),
		slidingUnit:            slidingUnit,
		requests:               requests,
	}
}

func (rl *SlidingWindowRateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastTimeNode) < rl.duration {
		i := now.Sub(rl.lastTimeNode) / rl.slidingUnit
		if sliceSum(rl.requests) < rl.rate {
			rl.requests[i]++
			return true
		}
		return false
	} else if int64(now.Sub(rl.lastTimeNode)) < int64(2*len(rl.requests)-1)*int64(rl.slidingUnit) {
		i := int(1 + now.Sub(rl.lastTimeNode)/rl.slidingUnit)
		for j := 0; j < len(rl.requests); j++ {
			if j < 2*len(rl.requests)-i {
				rl.requests[j] = rl.requests[j+i-len(rl.requests)]
			} else {
				rl.requests[j] = 1
			}
		}
		rl.lastTimeNode = rl.lastTimeNode.Add(time.Duration(int64(rl.slidingUnit) * int64(i-len(rl.requests))))
		if sliceSum(rl.requests) < rl.rate {
			rl.requests[i]++
			return true
		}
		return false
	} else {
		units := now.Sub(rl.lastTimeNode) / rl.slidingUnit
		rl.lastTimeNode = rl.lastTimeNode.Add(units * rl.duration)
		setZero(rl.requests)
		return true
	}
}

func (rl *SlidingWindowRateLimiter) UpdateLimiter(rate int, duration time.Duration, slidingUnit time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rate = rate
	rl.duration = duration
	rl.slidingUnit = slidingUnit
}
