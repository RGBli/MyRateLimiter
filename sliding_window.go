package myratelimiter

import (
	"time"
)

type SlidingWindowRateLimiter struct {
	*BaseRateLimiter
	units            int64
	reqInUnits       map[int64]int64
	thresholdPerUnit float64
	startTime        int64
}

func NewSlidingWindowRateLimiter(limitCount int64, duration int64, units int64) *SlidingWindowRateLimiter {
	reqInUnits := make(map[int64]int64)
	thresholdPerUnit := float64(limitCount / units)
	return &SlidingWindowRateLimiter{
		BaseRateLimiter:  NewBaseRateLimiter(limitCount, duration),
		units:            units,
		reqInUnits:       reqInUnits,
		thresholdPerUnit: thresholdPerUnit,
		startTime:        time.Now().Unix(),
	}
}

func (rl *SlidingWindowRateLimiter) Limit() bool {
	rl.mu.Lock()
	now := time.Now().Unix()
	index := (now - rl.startTime) / (rl.duration / rl.units)
	defer func() {
		// drop expired records with 10% sample rate
		if now%10 == 0 {
			rl.dropExpiredRecords(index)
		}
		rl.mu.Unlock()
	}()

	if float64(rl.reqInUnits[index]+1) <= rl.thresholdPerUnit {
		rl.reqInUnits[index]++
		return true
	}

	return false
}

func (rl *SlidingWindowRateLimiter) UpdateLimiter(limitCount int64, duration int64, units int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limitCount = limitCount
	rl.duration = duration
	rl.units = units
}

func (rl *SlidingWindowRateLimiter) dropExpiredRecords(i int64) {
	for k := range rl.reqInUnits {
		if k < i-rl.units {
			delete(rl.reqInUnits, k)
		}
	}
}
