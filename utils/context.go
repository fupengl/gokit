package utils

import "context"

// WithDefaultContext 如果 ctx 为 nil，返回 context.Background()
func WithDefaultContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
