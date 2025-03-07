package slice

// Reverse 反转给定的切片。
func Reverse[T any, Slice ~[]T](slice Slice) {
	length := len(slice)
	half := length / 2

	for i := 0; i < half; i = i + 1 {
		j := length - 1 - i
		slice[i], slice[j] = slice[j], slice[i]
	}
}
