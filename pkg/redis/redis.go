package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
}

func GetClient() *redis.Client {
	return client
}
