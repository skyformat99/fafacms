package kv

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type MyRedisConf struct {
	RedisHost        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
	RedisDB          int
	RedisPass        string
}

func NewRedis(redisConf *MyRedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConf.RedisHost, redis.DialPassword(redisConf.RedisPass), redis.DialDatabase(redisConf.RedisDB))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	return
}
