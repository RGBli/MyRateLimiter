package myratelimiter

import (
	"math"
	"time"
)

type TokenBucketRateLimiter struct {
	*BaseRateLimiter
	maxToken int
	tokens   float64
	speed    float64
}

func NewTokenBucketRateLimiter(limitCount int64, duration int64, maxToken int) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		BaseRateLimiter: NewBaseRateLimiter(limitCount, duration),
		maxToken:        maxToken,
		tokens:          float64(maxToken),
		speed:           float64(limitCount / duration),
	}
}

func (rl *TokenBucketRateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().Unix()
	interval := now - rl.lastRequestTime
	rl.lastRequestTime = now
	rl.tokens += float64(interval) * rl.speed
	rl.tokens = math.Min(rl.tokens, float64(rl.maxToken))
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func (rl *TokenBucketRateLimiter) UpdateLimiter(limitCount int64, duration int64, maxToken int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.limitCount = limitCount
	rl.duration = duration
	rl.maxToken = maxToken
}
