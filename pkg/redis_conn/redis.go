package redis_conn

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func NewRedisConnPool(addr, password string) (*redis.Pool, error) {
	//connectionString := fmt.Sprintf("%s:6379", host)
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 10,
		// Max number of connections.
		MaxActive: 1000,
		// Time to close a idle connection in the pool.
		IdleTimeout: time.Second * 60,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr, redis.DialPassword(password))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}, nil
}

var RedisConnPool *redis.Pool

func GetRedisConn() redis.Conn {
	redisPool := RedisConnPool.Get()
	defer redisPool.Close()
	return redisPool
}
