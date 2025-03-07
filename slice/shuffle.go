package slice

import (
	"math/rand"
)

// Shuffle 随机打乱给定的切片。
func Shuffle[T any, Slice ~[]T](slice Slice) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
