package syncutil

import (
	"fmt"
)

func ExampleOrderedMap() {
	om := NewOrderedMap[int, string]()
	om.Set(1, "one")
	value, _ := om.Get(1)
	fmt.Println(value)
	// Output: one
}
