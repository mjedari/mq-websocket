package utils

import (
	"fmt"
	"repo.abanicon.com/public-library/glogger"
	"sync"
)

type SafeMap struct {
	mu    sync.RWMutex
	items map[string]any
}

func NewSafeMap() *SafeMap {
	return &SafeMap{mu: sync.RWMutex{}, items: make(map[string]any)}
}

func (m *SafeMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
}

func (m *SafeMap) Load(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.items[key]
	fmt.Println("v, ok", v, ok)
	return v, ok
}

func (m *SafeMap) Store(key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = value
}

func (m *SafeMap) LoadOrStore(key string, value any) (any, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.items[key]
	if ok {
		return v, ok
	}
	m.items[key] = value
	return m.items[key], false
}

func (m *SafeMap) Range(f func(key any, value any) bool) {
	// todo: check this out whether we should lock or not
	//m.mu.RLock()
	//defer m.mu.RUnlock()
	for key, value := range m.items {
		if !f(key, value) {
			glogger.Errorf("error in iterating over map; key: %v value: %v\n", key, value)
			continue
		}
	}
}

func (m *SafeMap) Len() int {
	return len(m.items)
}

func (m *SafeMap) LoadAndDelete(key string) (any, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.items[key]
	delete(m.items, key)
	return v, ok
}
