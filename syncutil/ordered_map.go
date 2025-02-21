package syncutil

import (
	"sync"
)

type OrderedMap[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
	keys  []K
	index map[K]int
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		items: make(map[K]V),
		keys:  make([]K, 0),
		index: make(map[K]int),
	}
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, exists := om.items[key]; !exists {
		om.keys = append(om.keys, key)
		om.index[key] = len(om.keys) - 1
	}
	om.items[key] = value
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	value, exists := om.items[key]
	return value, exists
}

func (om *OrderedMap[K, V]) Delete(key K) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if idx, exists := om.index[key]; exists {
		om.keys = append(om.keys[:idx], om.keys[idx+1:]...)
		delete(om.index, key)
		delete(om.items, key)
		for i := idx; i < len(om.keys); i++ {
			om.index[om.keys[i]] = i
		}
	}
}

func (om *OrderedMap[K, V]) Keys() []K {
	om.mu.RLock()
	defer om.mu.RUnlock()

	keys := make([]K, len(om.keys))
	copy(keys, om.keys)
	return keys
}

func (om *OrderedMap[K, V]) Values() []V {
	om.mu.RLock()
	defer om.mu.RUnlock()

	values := make([]V, 0, len(om.keys))
	for _, key := range om.keys {
		values = append(values, om.items[key])
	}
	return values
}

func (om *OrderedMap[K, V]) Len() int {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return len(om.keys)
}

func (om *OrderedMap[K, V]) Range(f func(key K, value V) bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for _, key := range om.keys {
		if !f(key, om.items[key]) {
			break
		}
	}
}
