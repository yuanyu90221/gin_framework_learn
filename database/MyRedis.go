package database

import (
	"fmt"
	"time"
	"web/config"

	"github.com/gomodule/redigo/redis"
)

var RedisDefaultPool *redis.Pool

func newPool() *redis.Pool {
	Config := config.GetConfig()

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", Config.RedisHost, Config.RedisPort), redis.DialPassword(Config.RedisPassword))
		},
	}
}

func init() {
	RedisDefaultPool = newPool()
}
