package engine

import (
	"sync"
	"testing"
)

func TestMemory_Increment(t *testing.T) {
	limiter := NewEngineByMemory()
	key := "key"
	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 50000; j++ {
				if ok, _ := limiter.Increment(key, 1, 0, 1000000); !ok {
					t.Fail()
				}
			}
		}()
	}

	wg.Wait()
	curNum, _ := limiter.Get(key)
	if curNum != 1000000 {
		t.Fail()
	}
}
