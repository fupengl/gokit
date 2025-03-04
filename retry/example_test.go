package retry

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func ExampleRetry() {
	// 基本用法
	result, err := Retry(func() (string, error) {
		return "success", nil
	}, WithMaxAttempts(3), WithDelay(time.Second))
	fmt.Printf("Result: %s, Error: %v\n", result, err)

	// 使用指数退避
	exponentialBackoff := func(attempts int, err error) time.Duration {
		// 基础延迟时间
		baseDelay := time.Second
		// 最大延迟时间
		maxDelay := time.Minute
		// 退避因子
		factor := 2.0
		// 随机抖动因子
		jitter := 0.1

		// 计算基础延迟时间
		delay := float64(baseDelay)
		delay *= math.Pow(factor, float64(attempts-1))

		// 应用最大延迟限制
		if delay > float64(maxDelay) {
			delay = float64(maxDelay)
		}

		// 添加随机抖动
		if jitter > 0 {
			jitterAmount := delay * jitter
			delay += rand.Float64()*jitterAmount - jitterAmount/2
		}

		return time.Duration(delay)
	}

	result, err = Retry(func() (string, error) {
		return "", errors.New("temporary error")
	}, WithDelayFn(exponentialBackoff))

	// 使用自定义错误判断
	result, err = Retry(func() (string, error) {
		return "", errors.New("temporary error")
	}, WithIsRetryable(func(err error) bool {
		return err.Error() == "temporary error"
	}))

	// 使用重试回调
	result, err = Retry(func() (string, error) {
		return "", errors.New("error")
	}, WithOnRetry(func(attempts int, err error) {
		fmt.Printf("Attempt %d failed: %v\n", attempts, err)
	}))
}
