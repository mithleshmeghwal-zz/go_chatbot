package main

import (
	"time"
	"log"
	"github.com/gomodule/redigo/redis"
)


func ConnectDatabase(ip , port string) (redis.Conn, error) {
	var conn redis.Conn
	var err error
	for {
		time.Sleep(2 *time.Second)
		pool := newPool(ip, port)
		conn = pool.Get()
	
		err = ping(conn)	
		if err == nil {
			break
		}
	}
	if err := conn.Flush(); err != nil {
		log.Println(err)
	}
	return conn, err
}

func ping(c redis.Conn) error {
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}
	log.Println("PING Response = \n", s)
	return err
}

func newPool(ip, port string) *redis.Pool {
	return &redis.Pool{
	  MaxIdle: 3,
	  IdleTimeout: 240 * time.Second,
	  Dial: func () (redis.Conn, error) { 
		  	c, err := redis.Dial("tcp", ip+":"+port)
			if err != nil {
				log.Println(err)
			}
			return c, err
	  },
	}
}