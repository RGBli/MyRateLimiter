package myratelimiter

import (
	"testing"
)

// 对三种限流器的基准测试
func BenchmarkFixedWindowRateLimiter_Limit(b *testing.B) {
	rl := NewFixedWindowRateLimiter(1000, 1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rl.Limit()
		}
	})
}

func BenchmarkSlidingWindowRateLimiter_Limit(b *testing.B) {
	rl := NewSlidingWindowRateLimiter(1000, 1, 5)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rl.Limit()
		}
	})
}

func BenchmarkTokenBucketRateLimiter_Limit(b *testing.B) {
	rl := NewTokenBucketRateLimiter(1000, 1, 2000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rl.Limit()
		}
	})
}
