package model

import "github.com/gomodule/redigo/redis"

var RedisPool redis.Pool

//连接池
func InitRedis(){
	RedisPool = redis.Pool{
		MaxIdle:20,
		MaxActive:50,
		IdleTimeout:60 * 5,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp","192.168.137.81:6379")
		},
	}
}

//存储验证码
func SaveImgRnd(uuid ,rnd string)error{
	//链接redis
	conn := RedisPool.Get()
	//存储验证码
	_,err := conn.Do("setex",uuid,60 * 5,rnd)
	return err
}
