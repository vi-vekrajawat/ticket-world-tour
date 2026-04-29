package config

import (
	"context"
	"log"
	"github.com/redis/go-redis/v9"
)

var ClientRadis *redis.Client
var Ctx = context.Background()

func ConnectRadis(){
	ClientRadis = redis.NewClient(&redis.Options{
		Addr : "localhost:6379",
		Password : "",
		DB : 0,
	})

	_ , err := ClientRadis.Ping(Ctx).Result()

	if err!=nil{
		log.Println("Failed to integrate the redis database")
	}

	log.Println("Redis connected successfully")
}