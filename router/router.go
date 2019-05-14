package router

import "github.com/gomodule/redigo/redis"

type Router struct {
	Conn redis.Conn
}



