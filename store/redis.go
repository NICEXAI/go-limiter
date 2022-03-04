package store

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	luaScript = redis.NewScript(`
local key = KEYS[1]
local delta = tonumber(ARGV[1])
local min = tonumber(ARGV[2])
local max = tonumber(ARGV[3])

local counter = redis.call("GET", key)
if not counter then
	counter = 0
else 
	counter = tonumber(counter)
end

if (counter + delta > max) or (counter + delta < min) then
	return 0
end

counter = counter + delta
redis.call("SET", key, counter)

if counter == 0 then
	redis.call("EXPIRE", key, 5)
else 
	local expire = redis.call("TTL", key)
	if expire ~= -1 then
		redis.call("EXPIRE", key, -1)
	end
end

return 1
`)
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Get(key string) (int, error) {
	return r.client.Get(context.Background(), key).Int()
}

func (r *Redis) Increment(key string, delta, min, max int) (bool, error) {
	ok, err := luaScript.Run(context.Background(), r.client, []string{key}, delta, min, max).Bool()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func NewStoreByRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}
