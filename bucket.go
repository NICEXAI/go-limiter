package go_rate_limiter

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-rate-limiter/store"
)

type Token struct {
	bucket *Bucket
	key    string
	tokens int
}

// Free release a token
func (t *Token) Free() {
	_, _ = t.bucket.store.Increment(t.key, t.tokens, 0, t.bucket.burst)
}

type Bucket struct {
	burst int
	store store.Store
}

// Get get a token
func (b *Bucket) Get(key string) (Token, bool) {
	if ok, err := b.store.Increment(key, 1, 0, b.burst); err != nil || !ok {
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

func (b *Bucket) GetCounter(key string) int {
	counter, _ := b.store.Get(key)
	return counter
}

// NewBucketByMemory init new bucket
func NewBucketByMemory(burst uint) *Bucket {
	return &Bucket{
		burst: int(burst),
		store: store.NewStoreByMemory(),
	}
}

// NewBucketByRedis init new bucket
func NewBucketByRedis(client *redis.Client, burst uint) *Bucket {
	return &Bucket{
		burst: int(burst),
		store: store.NewStoreByRedis(client),
	}
}
