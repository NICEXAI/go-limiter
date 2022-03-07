package go_limiter

import (
	"github.com/NICEXAI/go-limiter/bucket"
	"github.com/NICEXAI/go-limiter/engine"
	"github.com/NICEXAI/go-limiter/rate"
	"github.com/go-redis/redis/v8"
	"time"
)

type Limiter struct {
	engine engine.Engine
}

func (l *Limiter) NewBucket(burst uint) *bucket.Bucket {
	return bucket.NewBucket(bucket.Options{
		Engine: l.engine,
		Burst:  burst,
	})
}

func (l *Limiter) NewRate(limit uint, period time.Duration) *rate.Rate {
	return rate.NewRate(rate.Options{
		Engine: l.engine,
		Limit:  limit,
		Period: period,
	})
}

func NewLimiterByMemory() *Limiter {
	return &Limiter{engine: engine.NewEngineByMemory()}
}

func NewLimiterByRedis(client *redis.Client) *Limiter {
	return &Limiter{engine: engine.NewEngineByRedis(client)}
}
