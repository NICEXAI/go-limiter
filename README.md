# go-limiter
 Easy-to-use distributed multi-function current limiter.

### Installation

Run the following command under your project:

> go get -u github.com/NICEXAI/go-limiter

### Usage

1. Use the `NewLimiterByMemory` or `NewLimiterByRedis` methods to initialize a Limiter
```
// Better performance, but no support for distributed
limiter := NewLimiterByMemory()

// or

// Distributed calls are supported, but depend on redis
client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
limiter := NewLimiterByRedis(client)
```

2. Rate limit, for example: allow up to 20 requests per second
```
key := "test"
rate := limiter.NewRate(20, 1*time.Second)
if rate.Allow(key) {
    // normal business process
} else {
    // abnormal business process
}
```
3. Limit concurrency, for example: a single user can initiate a maximum of 10 requests at a time
```
key := "test"
bucket := limiter.NewBucket(10)
token, ok := bucket.Get(key)
if ok {
    // normal business process
    ...
    // note: be sure to release the token after use
    token.Free()
} else {
    // abnormal business process
}
```