package engine

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"testing"
)

func TestRedis_Increment(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	limiter := NewEngineByRedis(RedisOption{
		Client: client,
		Expire: 10,
	})
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

func TestRedis_IncrementTo(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	limiter := NewEngineByRedis(RedisOption{
		Client: client,
		Expire: 20,
	})
	key := "key"

	client.Del(context.Background(), key)
	_, err := limiter.IncrementTo(key, -1, 0, 5, 2)
	if err != nil {
		t.Fail()
	}
}

func BenchmarkRedis_Increment(b *testing.B) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	limiter := NewEngineByRedis(RedisOption{
		Client: client,
		Expire: 10,
	})
	key := "key"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = limiter.Increment(key, 1, 0, 100000000)
	}
}
