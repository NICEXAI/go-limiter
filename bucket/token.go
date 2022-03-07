package bucket

type Token struct {
	bucket *Bucket
	key    string
	tokens int
}

func (t *Token) Free() {
	_, _ = t.bucket.engine.Increment(t.key, t.tokens, 0, t.bucket.burst)
}
