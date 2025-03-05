package syncutil

import (
	"testing"
)

func TestSyncMap(t *testing.T) {
	m := NewSyncMap[int, string]()

	m.Store(1, "one")
	value, ok := m.Load(1)
	if !ok || value != "one" {
		t.Errorf("expected 'one', got %v", value)
	}

	m.Delete(1)
	if _, ok := m.Load(1); ok {
		t.Error("expected key 1 to be deleted")
	}
}
