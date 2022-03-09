package bucket

import "sync/atomic"

type Token struct {
	bucket *Bucket
	key    string
	tokens int
}

func (t *Token) Free() {
	var (
		oCounter interface{}
		counter  *int64
		ok       bool
	)

	oCounter, ok = t.bucket.manager.Get(t.key)
	if !ok {
		return
	}

	counter, ok = oCounter.(*int64)
	if !ok {
		return
	}

	atomic.AddInt64(counter, -1)
	_, _ = t.bucket.engine.Increment(t.key, t.tokens, 0, t.bucket.burst)
}
