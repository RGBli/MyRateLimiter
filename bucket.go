package myratelimiter

import (
	"math"
	"time"
)

type BucketRateLimiter struct {
	*FixedWindowRateLimiter
	maxToken int
	tokens   float64
	speed    float64
}

func NewBucketRateLimiter(rate int, duration time.Duration, maxToken int) *BucketRateLimiter {
	return &BucketRateLimiter{
		FixedWindowRateLimiter: NewFixedWindowRateLimiter(rate, duration),
		maxToken:               maxToken,
		tokens:                 float64(maxToken),
		speed:                  float64(int64(rate) / int64(duration)),
	}
}

func (rl *BucketRateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	passed := now.Sub(rl.lastRequestTime)
	rl.lastRequestTime = now
	rl.tokens += float64(passed) * float64(rl.rate) / float64(rl.duration)
	rl.tokens = math.Min(rl.tokens, float64(rl.maxToken))
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func (rl *BucketRateLimiter) UpdateLimiter(rate int, duration time.Duration, maxToken int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rate = rate
	rl.duration = duration
	rl.maxToken = maxToken
}
