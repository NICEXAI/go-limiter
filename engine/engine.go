package engine

type Engine interface {
	Get(key string) (int, error)
	Increment(key string, delta, min, max int) (bool, error)
	IncrementTo(key string, delta, min, max, incr int) (bool, error)
}
