package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"testing"
)

func TestRedis_Increment(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	limiter := NewStoreByRedis(client)
	key := "key"
	wg := sync.WaitGroup{}

	client.Del(context.Background(), key)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 5000; j++ {
				if ok, err := limiter.Increment(key, 1, 0, 100000); err != nil || !ok {
					t.Fail()
				}
			}
		}()
	}

	wg.Wait()
	curNum, _ := limiter.Get(key)
	if curNum != 100000 {
		t.Fail()
	}
}
