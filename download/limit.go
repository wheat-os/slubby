package download

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// WithTokenBucketLimit 使用令牌桶限速
func WithTokenBucketLimit(l time.Duration, r int) OptFunc {
	limit := rate.NewLimiter(rate.Every(l), r)

	return func(opt *SlubbyComponent) {
		opt.rateLimit = func(ctx context.Context) error {
			return limit.Wait(ctx)
		}
	}
}