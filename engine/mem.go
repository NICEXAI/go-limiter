package engine

import (
	"math"
	"sync"
	"time"
)

type MemData struct {
	Counter  int
	LastTime time.Time
}

type Memory struct {
	mutex   sync.Mutex
	storage map[string]*MemData
}

func (m *Memory) Get(key string) (int, error) {
	memData := m.storage[key]
	if memData == nil {
		return 0, nil
	}
	return memData.Counter, nil
}

func (m *Memory) Increment(key string, delta, min, max int) (bool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	memData := m.storage[key]
	if memData == nil {
		memData = &MemData{}
		m.storage[key] = memData
	}
	if memData.Counter+delta > max || memData.Counter+delta < min {
		return false, nil
	}
	memData.Counter += delta
	return true, nil
}

func (m *Memory) IncrementTo(key string, delta, min, max, incr int) (bool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	memData := m.storage[key]
	if memData == nil {
		memData = &MemData{
			Counter:  max,
			LastTime: time.Now(),
		}
		m.storage[key] = memData
	}

	delayTime := time.Since(memData.LastTime)
	offsetTime := delayTime.Milliseconds() % 1000

	incrNum := int(math.Floor(delayTime.Seconds()) * float64(incr))
	if incrNum > 0 {
		if memData.Counter+incrNum > max {
			memData.Counter = max
		} else if memData.Counter+incrNum < min {
			memData.Counter = min
		} else {
			memData.Counter += incrNum
		}
		memData.LastTime = time.UnixMilli(time.Now().UnixMilli() - offsetTime)
	}

	if memData.Counter+delta > max || memData.Counter+delta < min {
		return false, nil
	}
	memData.Counter += delta
	return true, nil
}

func NewEngineByMemory() Engine {
	return &Memory{
		mutex:   sync.Mutex{},
		storage: map[string]*MemData{},
	}
}
