# MyRateLimiter
---

### 功能
用于接口限流

</br>

### 特点
包含固定窗口、滑动窗口和令牌桶三种常用限流算法，可以根据业务灵活选择

</br>

### 用法
```go
import (
    "github.com/RGBli/MyRateLimiter"
    "time"
)

rl := MyRateLimiter.NewFixedWindowRateLimiter(1000, time.Second)
if ok := rl.Limit(); ok {
    // pass rate limit
}
```

</br>