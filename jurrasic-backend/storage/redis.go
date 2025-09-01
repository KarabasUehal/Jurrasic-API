package storage

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func NewRedis() *redis.Client {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return redisClient
}

func GetRedisClient() *redis.Client {
	return redisClient
}
