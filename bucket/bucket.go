package bucket

import (
	"fmt"
	"github.com/NICEXAI/go-limiter/engine"
	cmap "github.com/orcaman/concurrent-map"
	"sync/atomic"
)

type Options struct {
	Engine engine.Engine
	Burst  uint
}

type Bucket struct {
	burst   int
	engine  engine.Engine
	manager cmap.ConcurrentMap
}

func (b *Bucket) Get(key string) (*Token, bool) {
	key = fmt.Sprintf("limiter:bucket:%v", key)

	if ok, err := b.engine.Increment(key, 1, 0, b.burst); err != nil || !ok {
		return &Token{}, false
	}

	token := &Token{
		bucket: b,
		key:    key,
		tokens: -1,
	}

	var (
		oCounter interface{}
		counter  *int64
		ok       bool
	)

	oCounter, ok = b.manager.Get(key)
	if !ok {
		nCounter := int64(0)
		counter = &nCounter
		b.manager.Set(key, counter)
	} else {
		counter, ok = oCounter.(*int64)
		if !ok {
			return token, false
		}
	}

	atomic.AddInt64(counter, 1)
	return token, true
}

func (b *Bucket) FreeAll() {
	b.manager.IterCb(func(key string, v interface{}) {
		counter, ok := v.(*int64)
		if ok {
			_, _ = b.engine.Increment(key, -int(*counter), 0, b.burst)
		}
	})
}

func NewBucket(opt Options) *Bucket {
	return &Bucket{burst: int(opt.Burst), engine: opt.Engine, manager: cmap.New()}
}
