package retry

import (
	"context"
	"time"
)

// RetryOption 定义重试选项
type RetryOption struct {
	// MaxAttempts 最大重试次数，0表示无限重试
	MaxAttempts int
	// Delay 重试间隔时间
	Delay time.Duration
	// DelayFn 动态返回间隔时间
	DelayFn func(attempts int, err error) time.Duration
	// IsRetryable 判断错误是否可重试
	IsRetryable func(err error) bool
	// Timeout 全局超时时间
	Timeout time.Duration
	// Context 上下文，用于取消重试
	Context context.Context
	// OnRetry 重试回调函数
	OnRetry func(attempts int, err error)
}

// WithMaxAttempts 设置最大重试次数
func WithMaxAttempts(attempts int) func(*RetryOption) {
	return func(o *RetryOption) {
		o.MaxAttempts = attempts
	}
}

// WithDelay 设置固定延迟时间
func WithDelay(delay time.Duration) func(*RetryOption) {
	return func(o *RetryOption) {
		o.Delay = delay
	}
}

// WithDelayFn 设置动态延迟函数
func WithDelayFn(fn func(attempts int, err error) time.Duration) func(*RetryOption) {
	return func(o *RetryOption) {
		o.DelayFn = fn
	}
}

// WithIsRetryable 设置错误重试判断函数
func WithIsRetryable(fn func(err error) bool) func(*RetryOption) {
	return func(o *RetryOption) {
		o.IsRetryable = fn
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) func(*RetryOption) {
	return func(o *RetryOption) {
		o.Timeout = timeout
	}
}

// WithContext 设置上下文
func WithContext(ctx context.Context) func(*RetryOption) {
	return func(o *RetryOption) {
		o.Context = ctx
	}
}

// WithOnRetry 设置重试回调函数
func WithOnRetry(fn func(attempts int, err error)) func(*RetryOption) {
	return func(o *RetryOption) {
		o.OnRetry = fn
	}
}
