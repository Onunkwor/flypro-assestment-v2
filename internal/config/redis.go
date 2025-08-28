package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis() {
	Addr, err := Getenv("REDIS_ADDR")
	if err != nil {
		log.Fatalf("environment variable REDIS_ADDR not set: %v", err)
	}
	Redis = redis.NewClient(&redis.Options{
		Addr: Addr,
		DB:   0,
	})

	_, err = Redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("✅ Successfully connected to Redis")
}
