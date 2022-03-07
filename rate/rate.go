package rate

import (
	"github.com/NICEXAI/go-limiter/engine"
	"math"
	"time"
)

type Options struct {
	Engine engine.Engine
	Period time.Duration
	Limit  uint
}

type Rate struct {
	engine   engine.Engine
	burst    int
	period   time.Duration
	limit    uint
	speed    int
	lastTime time.Time
}

func (r *Rate) Allow(key string) bool {
	ok, err := r.engine.IncrementTo(key, -1, 0, r.burst, r.speed)
	if err != nil {
		return false
	}
	return ok
}

func (r *Rate) Count(key string) int {
	counter, _ := r.engine.Get(key)
	return counter
}

func NewRate(opt Options) *Rate {
	speed := math.Floor(float64(opt.Limit) / opt.Period.Seconds())
	// At least one access per second is allowed
	if speed == 0 {
		speed = 1
	}

	return &Rate{
		engine: opt.Engine,
		burst:  int(speed),
		period: opt.Period,
		limit:  opt.Limit,
		speed:  int(speed),
	}
}
