package store

type Store interface {
	Get(key string) (int, error)
	Increment(key string, delta, min, max int) (bool, error)
}
