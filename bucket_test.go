package go_rate_limiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewBucketByMemory(t *testing.T) {
	bucket := NewBucketByMemory(900)
	key := fmt.Sprintf("limiter:bucket:%s", "test")
	tCount := 0

	for {
		var (
			successCount int64
			failedCount  int64
		)
		wg := sync.WaitGroup{}

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				token, ok := bucket.Get(key)
				if ok {
					atomic.AddInt64(&successCount, 1)
					time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)
					token.Free()
				} else {
					atomic.AddInt64(&failedCount, 1)
				}
			}()
		}

		wg.Wait()
		if successCount != 900 || failedCount != 100 {
			t.Fail()
			return
		}
		tCount += 1
		if tCount > 10 {
			return
		}
	}
}

func TestNewBucketByRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	bucket := NewBucketByRedis(client, 900)
	key := fmt.Sprintf("limiter:bucket:%s", "test")
	client.Del(context.Background(), key)
	tCount := 0

	for {
		var (
			successCount int64
			failedCount  int64
		)
		wg := sync.WaitGroup{}

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				token, ok := bucket.Get(key)
				if ok {
					atomic.AddInt64(&successCount, 1)
					time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)
					token.Free()
				} else {
					atomic.AddInt64(&failedCount, 1)
				}
			}()
		}

		wg.Wait()
		if successCount != 900 || failedCount != 100 {
			t.Fail()
			return
		}
		tCount += 1
		if tCount > 10 {
			return
		}
	}
}
