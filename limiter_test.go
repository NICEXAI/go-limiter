package go_limiter

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewLimiterByMemory_NewBucket(t *testing.T) {
	limiter := NewLimiterByMemory()
	bucket := limiter.NewBucket(900)
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

func TestNewLimiterByRedis_NewBucket(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	limiter := NewLimiterByRedis(client)
	bucket := limiter.NewBucket(900)
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

func TestNewLimiterByMemory_NewRate(t *testing.T) {
	limiter := NewLimiterByMemory()
	key := fmt.Sprintf("limiter:rate:%s", "test")
	continueTime := 1000
	tCount := 0

	rate := limiter.NewRate(60, 1*time.Second)

	for {
		var (
			successCount int64
			failedCount  int64
		)
		wg := sync.WaitGroup{}

		for i := 0; i < 400; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				time.Sleep(time.Duration(rand.Intn(continueTime)) * time.Millisecond)
				if rate.Allow(key) {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&failedCount, 1)
				}
			}()
		}
		wg.Wait()
		//log.Printf("success: %v, failed: %v", successCount, failedCount)
		if successCount == 0 || failedCount == 0 {
			t.Fail()
		}
		tCount += 1
		if tCount > 10 {
			return
		}
	}
}

func TestNewLimiterByRedis_NewRate(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	limiter := NewLimiterByRedis(client)
	key := fmt.Sprintf("limiter:rate:%s", "test")
	continueTime := 1000
	tCount := 0

	rate := limiter.NewRate(60, 1*time.Second)

	for {
		var (
			successCount int64
			failedCount  int64
		)
		wg := sync.WaitGroup{}

		for i := 0; i < 400; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				time.Sleep(time.Duration(rand.Intn(continueTime)) * time.Millisecond)
				if rate.Allow(key) {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&failedCount, 1)
				}
			}()
		}
		wg.Wait()
		//log.Printf("success: %v, failed: %v", successCount, failedCount)
		if successCount == 0 || failedCount == 0 {
			t.Fail()
		}
		tCount += 1
		if tCount > 10 {
			return
		}
	}
}
