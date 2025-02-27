package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fupengl/goutils/utils"
)

// TestRetry_Success 测试函数第一次执行成功
func TestRetry_Success(t *testing.T) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		return "success", nil
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		Delay:       100 * time.Millisecond,
	}

	result, err := Retry(fn, opt)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != "success" {
		t.Fatalf("Expected result 'success', got: %v", result)
	}
	if attempts != 1 {
		t.Fatalf("Expected 1 attempt, got: %d", attempts)
	}
}

// TestRetry_RetryThenSuccess 测试函数在重试几次后成功
func TestRetry_RetryThenSuccess(t *testing.T) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("temporary failure")
		}
		return "success", nil
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		Delay:       100 * time.Millisecond,
	}

	result, err := Retry(fn, opt)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != "success" {
		t.Fatalf("Expected result 'success', got: %v", result)
	}
	if attempts != 3 {
		t.Fatalf("Expected 3 attempts, got: %d", attempts)
	}
}

// TestRetry_MaxAttempts 测试达到最大重试次数
func TestRetry_MaxAttempts(t *testing.T) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		return "", errors.New("temporary failure")
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		Delay:       100 * time.Millisecond,
	}

	_, err := Retry(fn, opt)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	expectedErr := "could not complete function within 3 attempts: temporary failure"
	if err.Error() != expectedErr {
		t.Fatalf("Expected error '%s', got: %v", expectedErr, err)
	}
	if attempts != 3 {
		t.Fatalf("Expected 3 attempts, got: %d", attempts)
	}
}

// TestRetry_NonRetryableError 测试不可重试的错误
func TestRetry_NonRetryableError(t *testing.T) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		return "", errors.New("non-retryable failure")
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		Delay:       100 * time.Millisecond,
		IsRetryable: func(err error) bool {
			return err.Error() == "temporary failure"
		},
	}

	_, err := Retry(fn, opt)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	if err.Error() != "non-retryable failure" {
		t.Fatalf("Expected error 'non-retryable failure', got: %v", err)
	}
	if attempts != 1 {
		t.Fatalf("Expected 1 attempt, got: %d", attempts)
	}
}

// TestRetry_Timeout 测试超时场景
func TestRetry_Timeout(t *testing.T) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		time.Sleep(200 * time.Millisecond)
		return "", errors.New("temporary failure")
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(5),
		Delay:       100 * time.Millisecond,
		Timeout:     300 * time.Millisecond,
	}

	_, err := Retry(fn, opt)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	expectedErr := "retry cancelled: context deadline exceeded"
	if err.Error() != expectedErr {
		t.Fatalf("Expected error '%s', got: %v", expectedErr, err)
	}
	if attempts < 1 {
		t.Fatalf("Expected at least 1 attempt, got: %d", attempts)
	}
}

// TestRetry_ContextCancel 测试上下文取消
func TestRetry_ContextCancel(t *testing.T) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		return "", errors.New("temporary failure")
	}

	ctx, cancel := context.WithCancel(context.Background())
	opt := RetryOption{
		MaxAttempts: utils.Ptr(5),
		Delay:       100 * time.Millisecond,
		Context:     ctx,
	}

	// 取消上下文
	go func() {
		time.Sleep(150 * time.Millisecond)
		cancel()
	}()

	_, err := Retry(fn, opt)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	expectedErr := "retry cancelled: context canceled"
	if err.Error() != expectedErr {
		t.Fatalf("Expected error '%s', got: %v", expectedErr, err)
	}
	if attempts < 1 {
		t.Fatalf("Expected at least 1 attempt, got: %d", attempts)
	}
}

// TestRetry_DynamicDelay 测试动态延迟
func TestRetry_DynamicDelay(t *testing.T) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("temporary failure")
		}
		return "success", nil
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		DelayFn: func(attempts int) time.Duration {
			return time.Duration(attempts) * 100 * time.Millisecond
		},
	}

	start := time.Now()
	result, err := Retry(fn, opt)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != "success" {
		t.Fatalf("Expected result 'success', got: %v", result)
	}
	if attempts != 3 {
		t.Fatalf("Expected 3 attempts, got: %d", attempts)
	}

	// 验证动态延迟
	elapsed := time.Since(start)
	expectedMinDelay := 100*time.Millisecond + 200*time.Millisecond // 第1次和第2次重试的延迟
	if elapsed < expectedMinDelay {
		t.Fatalf("Expected elapsed time >= %v, got: %v", expectedMinDelay, elapsed)
	}
}

// BenchmarkRetry_Success 测试无重试的成功场景
func BenchmarkRetry_Success(b *testing.B) {
	fn := func() (string, error) {
		return "success", nil
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		Delay:       100 * time.Millisecond,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Retry(fn, opt)
	}
}

// BenchmarkRetry_RetryThenSuccess 测试多次重试后成功的场景
func BenchmarkRetry_RetryThenSuccess(b *testing.B) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("temporary failure")
		}
		return "success", nil
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		Delay:       100 * time.Millisecond,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attempts = 0
		_, _ = Retry(fn, opt)
	}
}

// BenchmarkRetry_Concurrent 测试高并发场景下的性能
func BenchmarkRetry_Concurrent(b *testing.B) {
	fn := func() (string, error) {
		return "success", nil
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		Delay:       100 * time.Millisecond,
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = Retry(fn, opt)
		}
	})
}

// BenchmarkRetry_DynamicDelay 测试动态延迟的性能
func BenchmarkRetry_DynamicDelay(b *testing.B) {
	attempts := 0
	fn := func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("temporary failure")
		}
		return "success", nil
	}

	opt := RetryOption{
		MaxAttempts: utils.Ptr(3),
		DelayFn: func(attempts int) time.Duration {
			return time.Duration(attempts) * 100 * time.Millisecond
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attempts = 0
		_, _ = Retry(fn, opt)
	}
}
