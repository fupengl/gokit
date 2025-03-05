package syncutil

import (
	"testing"
)

func BenchmarkSyncMap_Store(b *testing.B) {
	m := NewSyncMap[int, string]()
	for i := 0; i < b.N; i++ {
		m.Store(i, "value")
	}
}

func BenchmarkSyncMap_Load(b *testing.B) {
	m := NewSyncMap[int, string]()
	for i := 0; i < 1000; i++ {
		m.Store(i, "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(i % 1000)
	}
}
