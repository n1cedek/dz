package rate_limiter

import (
	"golang.org/x/net/context"
	"time"
)

type TokenBucketLimiter struct {
	tokenBucketCh chan struct{}
}

func NewTokenBucketLimiter(ctx context.Context, limit int, period time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{tokenBucketCh: make(chan struct{}, limit)}

	for i := 0; i < limit; i++ {
		limiter.tokenBucketCh <- struct{}{}
	}
	replInterval := period.Nanoseconds() / int64(limit)
	go limiter.stPer(ctx, time.Duration(replInterval))

	return limiter
}

func (t *TokenBucketLimiter) stPer(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.tokenBucketCh <- struct{}{}
		}
	}
}

func (t *TokenBucketLimiter) Allow() bool {
	select {
	case <-t.tokenBucketCh:
		return true
	default:
		return false
	}
}
