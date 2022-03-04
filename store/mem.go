package store

import (
	"sync"
)

type Memory struct {
	mutex   sync.Mutex
	storage map[string]int
}

func (m *Memory) Get(key string) (int, error) {
	return m.storage[key], nil
}

func (m *Memory) Increment(key string, delta, min, max int) (bool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	counter := m.storage[key]
	if counter+delta > max || counter+delta < min {
		return false, nil
	}
	counter += delta
	m.storage[key] = counter
	return true, nil
}

func NewStoreByMemory() Store {
	return &Memory{
		mutex:   sync.Mutex{},
		storage: map[string]int{},
	}
}
