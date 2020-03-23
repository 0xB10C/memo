package database

import (
	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/logger"
	"github.com/gomodule/redigo/redis"
)

// RedisPool holds a pool of redis connections
type RedisPool struct {
	*redis.Pool
}

func newPool() *redis.Pool {
	dbUser := config.GetString("redis.user")
	dbPasswd := config.GetString("redis.passwd")
	dbHost := config.GetString("redis.host")
	dbPort := config.GetString("redis.port")
	dbConnection := config.GetString("redis.connection")

	connectionString := dbConnection + "://" + dbUser + ":" + dbPasswd + "@" + dbHost + ":" + dbPort

	return &redis.Pool{
		MaxIdle:   40,
		MaxActive: 1200, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(connectionString)
			if err != nil {
				return c, err
			}
			return c, err
		},
	}
}

// SetupRedis sets up a new Redis Pool
func SetupRedis() (pool *RedisPool, err error) {
	p := newPool()

	c := p.Get() // get a new connection
	defer c.Close()

	_, err = c.Do("PING")
	if err != nil {
		return nil, err
	}

	logger.Info.Println("Setup redis database connection pool")

	return &RedisPool{p}, nil
}
