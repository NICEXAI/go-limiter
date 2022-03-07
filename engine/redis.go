package engine

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	luaIncrementScript = redis.NewScript(`
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

	luaIncrementToScript = redis.NewScript(`
local key = KEYS[1]
local delta = tonumber(ARGV[1])
local min = tonumber(ARGV[2])
local max = tonumber(ARGV[3])
local incr = tonumber(ARGV[4])

local counter = redis.call("HGET", key, "counter")
local last_time = redis.call("HGET", key, "last_time")

local cur_time = redis.call("TIME")[1]

if not counter or not last_time then
	counter = max
else
	local incrNum = (cur_time - last_time) * incr

	if incrNum > 0 then
		if counter + incrNum > max then
			counter = max
		elseif counter + incrNum < min then
			counter = min
		else
			counter = counter + incrNum	
		end
	end
end

if (counter + delta > max) or (counter + delta < min) then
	return 0
end

counter = counter + delta
redis.call("HSET", key, "counter", counter)
redis.call("HSET", key, "last_time", cur_time)

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
	ok, err := luaIncrementScript.Run(context.Background(), r.client, []string{key}, delta, min, max).Bool()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *Redis) IncrementTo(key string, delta, min, max, incr int) (bool, error) {
	ok, err := luaIncrementToScript.Run(context.Background(), r.client, []string{key}, delta, min, max, incr).Bool()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func NewEngineByRedis(client *redis.Client) Engine {
	return &Redis{client: client}
}
