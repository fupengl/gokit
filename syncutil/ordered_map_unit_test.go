package syncutil

import (
	"testing"
)

func TestOrderedMap(t *testing.T) {
	om := NewOrderedMap[int, string]()

	om.Set(1, "one")
	value, ok := om.Get(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got %v", value)
	}

	om.Delete(1)
	if _, ok := om.Get(1); ok {
		t.Error("expected key 1 to be deleted")
	}
}
