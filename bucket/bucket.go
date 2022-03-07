package bucket

import (
	"fmt"
	"github.com/NICEXAI/go-limiter/engine"
)

type Options struct {
	Engine engine.Engine
	Burst  uint
}

type Bucket struct {
	burst  int
	engine engine.Engine
}

func (b *Bucket) Get(key string) (Token, bool) {
	if ok, err := b.engine.Increment(key, 1, 0, b.burst); err != nil || !ok {
		if err != nil {
			fmt.Println(err)
		}
		return Token{}, false
	}

	return Token{
		bucket: b,
		key:    key,
		tokens: -1,
	}, true
}

func (b *Bucket) Count(key string) int {
	counter, _ := b.engine.Get(key)
	return counter
}

func NewBucket(opt Options) *Bucket {
	return &Bucket{burst: int(opt.Burst), engine: opt.Engine}
}
