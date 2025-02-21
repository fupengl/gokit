package syncutil

import (
	"testing"
)

func BenchmarkOrderedMap_Set(b *testing.B) {
	om := NewOrderedMap[int, int]()
	for i := 0; i < b.N; i++ {
		om.Set(i, i)
	}
}

func BenchmarkOrderedMap_Get(b *testing.B) {
	om := NewOrderedMap[int, int]()
	for i := 0; i < 1000; i++ {
		om.Set(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		om.Get(i % 1000)
	}
}

func BenchmarkOrderedMap_Delete(b *testing.B) {
	om := NewOrderedMap[int, int]()
	for i := 0; i < 1000; i++ {
		om.Set(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		om.Delete(i % 1000)
	}
}
