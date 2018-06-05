package utils

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

//RC is the global redis connection handle
var RC redis.Conn

//KS is the global redis key separator
var KS = ":"

//GetRedisConnection returns a connection to Redis
func GetRedisConnection() redis.Conn {
	rc, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v\n", err)
	}
	return rc
}
