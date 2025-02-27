package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/fupengl/goutils/utils"
)

// Retry 重试函数
func Retry[T any](fn func() (T, error), opt RetryOption) (T, error) {
	var zero T
	var lastErr error
	attempts := 0

	opt.Context = utils.WithDefaultContext(opt.Context)

	if opt.Timeout > 0 {
		var cancel context.CancelFunc
		opt.Context, cancel = context.WithTimeout(opt.Context, opt.Timeout)
		defer cancel()
	}

	for {
		attempts++
		result, err := fn()
		if err == nil {
			return result, nil
		}

		if opt.IsRetryable != nil && !opt.IsRetryable(err) {
			return zero, err
		}

		lastErr = err

		if opt.MaxAttempts != nil && attempts >= *opt.MaxAttempts {
			return zero, fmt.Errorf("could not complete function within %d attempts: %w", *opt.MaxAttempts, lastErr)
		}

		delay := calculateDelay(attempts, opt.Delay, opt.DelayFn)
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-opt.Context.Done():
				return zero, fmt.Errorf("retry cancelled: %w", opt.Context.Err())
			}
		}

		select {
		case <-opt.Context.Done():
			return zero, fmt.Errorf("retry cancelled: %w", opt.Context.Err())
		default:
		}
	}
}

func calculateDelay(attempts int, fixedDelay time.Duration, delayFn func(int) time.Duration) time.Duration {
	if delayFn != nil {
		return delayFn(attempts)
	}
	return fixedDelay
}
