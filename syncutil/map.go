package syncutil

import "sync"

type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

func (sm *SyncMap[K, V]) Store(key K, value V) {
	sm.m.Store(key, value)
}

func (sm *SyncMap[K, V]) Load(key K) (V, bool) {
	value, ok := sm.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return value.(V), true
}

func (sm *SyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := sm.m.LoadOrStore(key, value)
	return actual.(V), loaded
}

func (sm *SyncMap[K, V]) Delete(key K) {
	sm.m.Delete(key)
}

func (sm *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	sm.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func (sm *SyncMap[K, V]) Len() int {
	count := 0
	sm.m.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

func (sm *SyncMap[K, V]) Clear() {
	sm.m.Range(func(key, _ any) bool {
		sm.m.Delete(key)
		return true
	})
}

func (sm *SyncMap[K, V]) Keys() []K {
	keys := make([]K, 0)
	sm.m.Range(func(key, _ any) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

func (sm *SyncMap[K, V]) Values() []V {
	values := make([]V, 0)
	sm.m.Range(func(_, value any) bool {
		values = append(values, value.(V))
		return true
	})
	return values
}

func (sm *SyncMap[K, V]) Copy() *SyncMap[K, V] {
	newMap := NewSyncMap[K, V]()
	sm.Range(func(key K, value V) bool {
		newMap.Store(key, value)
		return true
	})
	return newMap
}

func (sm *SyncMap[K, V]) Merge(other *SyncMap[K, V]) {
	other.Range(func(key K, value V) bool {
		sm.Store(key, value)
		return true
	})
}

func (sm *SyncMap[K, V]) IsEmpty() bool {
	return sm.Len() == 0
}

func (sm *SyncMap[K, V]) Contains(key K) bool {
	_, ok := sm.Load(key)
	return ok
}
