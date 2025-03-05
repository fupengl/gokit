package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(*int) func() (string, error)
		opts     []func(*RetryOption)
		want     string
		wantErr  bool
		attempts int
	}{
		{
			name: "success on first try",
			fn: func(attempts *int) func() (string, error) {
				return func() (string, error) {
					*attempts++
					return "success", nil
				}
			},
			opts:     []func(*RetryOption){WithMaxAttempts(3)},
			want:     "success",
			wantErr:  false,
			attempts: 1,
		},
		{
			name: "success after retries",
			fn: func(attempts *int) func() (string, error) {
				return func() (string, error) {
					*attempts++
					if *attempts == 1 {
						return "", errors.New("temporary error")
					}
					return "success", nil
				}
			},
			opts:     []func(*RetryOption){WithMaxAttempts(3), WithDelay(time.Millisecond)},
			want:     "success",
			wantErr:  false,
			attempts: 2,
		},
		{
			name: "max attempts reached",
			fn: func(attempts *int) func() (string, error) {
				return func() (string, error) {
					*attempts++
					return "", errors.New("permanent error")
				}
			},
			opts:     []func(*RetryOption){WithMaxAttempts(2)},
			want:     "",
			wantErr:  true,
			attempts: 2,
		},
		{
			name: "context cancelled",
			fn: func(attempts *int) func() (string, error) {
				return func() (string, error) {
					*attempts++
					return "", errors.New("error")
				}
			},
			opts: []func(*RetryOption){
				func() func(*RetryOption) {
					ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
					return WithContext(ctx)
				}(),
				WithDelay(time.Second),
			},
			want:     "",
			wantErr:  true,
			attempts: 1,
		},
		{
			name: "custom retryable error",
			fn: func(attempts *int) func() (string, error) {
				return func() (string, error) {
					*attempts++
					return "", errors.New("temporary error")
				}
			},
			opts: []func(*RetryOption){
				WithMaxAttempts(3),
				WithIsRetryable(func(err error) bool {
					return err.Error() == "temporary error"
				}),
			},
			want:     "",
			wantErr:  true,
			attempts: 3,
		},
		{
			name: "non-retryable error",
			fn: func(attempts *int) func() (string, error) {
				return func() (string, error) {
					*attempts++
					return "", errors.New("permanent error")
				}
			},
			opts: []func(*RetryOption){
				WithMaxAttempts(3),
				WithIsRetryable(func(err error) bool {
					return err.Error() == "temporary error"
				}),
			},
			want:     "",
			wantErr:  true,
			attempts: 1,
		},
		{
			name: "custom delay function",
			fn: func(attempts *int) func() (string, error) {
				return func() (string, error) {
					*attempts++
					if *attempts < 3 {
						return "", errors.New("temporary error")
					}
					return "success", nil
				}
			},
			opts: []func(*RetryOption){
				WithMaxAttempts(3),
				WithDelayFn(func(attempts int, err error) time.Duration {
					return time.Duration(attempts) * time.Millisecond
				}),
			},
			want:     "success",
			wantErr:  false,
			attempts: 3,
		},
		{
			name: "on retry callback",
			fn: func(attempts *int) func() (string, error) {
				return func() (string, error) {
					*attempts++
					if *attempts < 2 {
						return "", errors.New("temporary error")
					}
					return "success", nil
				}
			},
			opts: []func(*RetryOption){
				WithMaxAttempts(3),
				WithDelay(time.Millisecond),
				WithOnRetry(func(attempts int, err error) {
					if attempts != 1 {
						t.Errorf("OnRetry callback got attempts = %v, want 1", attempts)
					}
				}),
			},
			want:     "success",
			wantErr:  false,
			attempts: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var attempts int
			got, err := Retry(tt.fn(&attempts), tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Retry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Retry() = %v, want %v", got, tt.want)
			}
			if attempts != tt.attempts {
				t.Errorf("Retry() attempts = %v, want %v", attempts, tt.attempts)
			}
		})
	}
}
