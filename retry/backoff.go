package retry

import "time"

// ExponentialBackoff 指数退避策略
func ExponentialBackoff(baseDelay time.Duration, maxDelay time.Duration) func(int) time.Duration {
	return func(attempts int) time.Duration {
		delay := baseDelay * time.Duration(1<<(attempts-1))
		if delay > maxDelay {
			return maxDelay
		}
		return delay
	}
}
