package retry

import (
	"errors"
	"testing"
	"time"
)

func BenchmarkRetry(b *testing.B) {
	// 测试成功场景
	b.Run("success", func(b *testing.B) {
		var attempts int
		fn := func() (string, error) {
			attempts++
			return "success", nil
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Retry(fn, WithMaxAttempts(3))
		}
	})

	// 测试失败场景
	b.Run("failure", func(b *testing.B) {
		var attempts int
		fn := func() (string, error) {
			attempts++
			return "", errors.New("error")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Retry(fn, WithMaxAttempts(3))
		}
	})

	// 测试重试成功场景
	b.Run("retry_success", func(b *testing.B) {
		var attempts int
		fn := func() (string, error) {
			attempts++
			if attempts == 1 {
				return "", errors.New("temporary error")
			}
			return "success", nil
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			attempts = 0
			Retry(fn, WithMaxAttempts(3), WithDelay(time.Millisecond))
		}
	})

	// 测试自定义延迟函数
	b.Run("custom_delay", func(b *testing.B) {
		var attempts int
		fn := func() (string, error) {
			attempts++
			return "success", nil
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Retry(fn, WithDelayFn(func(attempts int, err error) time.Duration {
				return time.Duration(attempts) * time.Millisecond
			}))
		}
	})

	// 测试错误判断函数
	b.Run("error_check", func(b *testing.B) {
		var attempts int
		fn := func() (string, error) {
			attempts++
			return "", errors.New("error")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Retry(fn, WithIsRetryable(func(err error) bool {
				return err.Error() == "error"
			}))
		}
	})

	// 测试重试回调
	b.Run("retry_callback", func(b *testing.B) {
		var attempts int
		fn := func() (string, error) {
			attempts++
			return "", errors.New("error")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Retry(fn, WithOnRetry(func(attempts int, err error) {}))
		}
	})
}
