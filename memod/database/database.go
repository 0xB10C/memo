package database

import (
	"github.com/gomodule/redigo/redis"

	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/logger"
)

var (
	Pool *redis.Pool
)

func newPool() *redis.Pool {

	dbUser := config.GetString("redis.user")
	dbPasswd := config.GetString("redis.passwd")
	dbHost := config.GetString("redis.host")
	dbPort := config.GetString("redis.port")
	dbConnection := config.GetString("redis.connection")

	connectionString := dbConnection + "://" + dbUser + ":" + dbPasswd + "@" + dbHost + ":" + dbPort

	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(connectionString)
			if err != nil {
				return c, err
			}
			return c, err
		},
	}

}

func SetupRedis() (err error) {
	Pool = newPool()

	c := Pool.Get() // get a new connection
	defer c.Close()

	_, err = c.Do("PING")
	if err != nil {
		return err
	}

	logger.Info.Println("Setup redis database connection pool")

	return nil
}
