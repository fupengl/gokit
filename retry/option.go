package retry

import (
	"context"
	"time"
)

type RetryOption struct {
	// MaxAttempts 最大重试次数，如果不设置则一直重试
	MaxAttempts *int
	// Delay 重试间隔时间
	Delay time.Duration
	// DelayFn 动态返回间隔时间
	DelayFn func(attempts int) time.Duration
	// IsRetryable 判断错误是否可重试
	IsRetryable func(err error) bool
	// Timeout 全局超时时间
	Timeout time.Duration
	// Context 上下文，用于取消重试
	Context context.Context
}
