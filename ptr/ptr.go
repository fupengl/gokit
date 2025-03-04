package ptr

// Ptr 返回值的指针
func Ptr[T any](v T) *T {
	return &v
}
