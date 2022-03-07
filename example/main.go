package main

import (
	"fmt"
	"github.com/NICEXAI/go-limiter/engine"
	rate2 "github.com/NICEXAI/go-limiter/rate"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	key := fmt.Sprintf("limiter:rate:%s", "test1")
	continueTime := 1000

	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	rate := rate2.NewRate(rate2.Options{
		Engine: engine.NewEngineByRedis(client),
		Period: 1 * time.Second,
		Limit:  30,
	})

	for {
		var (
			successCount int64
			failedCount  int64
		)
		wg := sync.WaitGroup{}

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				time.Sleep(time.Duration(rand.Intn(continueTime)) * time.Millisecond)
				lock := rate.Allow(key)
				if lock {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&failedCount, 1)
				}
			}()
		}

		wg.Wait()
		log.Printf("success: %v, failed: %v, counter: %v", successCount, failedCount, rate.Count(key))
		//break
	}
}
