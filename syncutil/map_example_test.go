package syncutil

import (
	"fmt"
)

func ExampleSyncMap() {
	m := NewSyncMap[int, string]()
	m.Store(1, "one")
	value, _ := m.Load(1)
	fmt.Println(value)
	// Output: one
}
