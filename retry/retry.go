package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/fupengl/gokit/contextutil"
)

// Retry 重试函数
func Retry[T any](fn func() (T, error), opts ...func(*RetryOption)) (T, error) {
	var zero T
	opt := &RetryOption{
		MaxAttempts: 0, // 默认无限重试
	}

	// 应用选项
	for _, o := range opts {
		o(opt)
	}

	// 设置默认上下文
	opt.Context = contextutil.WithDefaultContext(opt.Context)

	// 设置超时
	if opt.Timeout > 0 {
		var cancel context.CancelFunc
		opt.Context, cancel = context.WithTimeout(opt.Context, opt.Timeout)
		defer cancel()
	}

	var lastErr error
	attempts := 0

	for {
		attempts++
		result, err := fn()
		if err == nil {
			return result, nil
		}

		// 检查错误是否可重试
		if opt.IsRetryable != nil && !opt.IsRetryable(err) {
			return zero, err
		}

		lastErr = err

		// 检查是否达到最大重试次数
		if opt.MaxAttempts > 0 && attempts >= opt.MaxAttempts {
			return zero, fmt.Errorf("could not complete function within %d attempts: %w", opt.MaxAttempts, lastErr)
		}

		// 计算延迟时间
		delay := calculateDelay(attempts, err, opt)
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-opt.Context.Done():
				return zero, fmt.Errorf("retry cancelled: %w", opt.Context.Err())
			}
		}

		// 调用重试回调
		if opt.OnRetry != nil {
			opt.OnRetry(attempts, err)
		}

		select {
		case <-opt.Context.Done():
			return zero, fmt.Errorf("retry cancelled: %w", opt.Context.Err())
		default:
		}
	}
}

// calculateDelay 计算重试延迟时间
func calculateDelay(attempts int, err error, opt *RetryOption) time.Duration {
	// 优先使用自定义延迟函数
	if opt.DelayFn != nil {
		return opt.DelayFn(attempts, err)
	}

	// 使用固定延迟
	return opt.Delay
}
