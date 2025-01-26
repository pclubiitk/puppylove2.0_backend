package redisclient

import (
	"context"
	"github.com/redis/go-redis/v9"
	"fmt"
)

var (
	Ctx = context.Background()
	RedisClient *redis.Client
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func ViewRedis() {
	keys, err := RedisClient.Keys(Ctx, "*").Result()
	if err != nil {
		fmt.Println("Error fetching keys:", err)
		return
	}
	for _, key := range keys {
		val, err := RedisClient.HGetAll(Ctx, key).Result()
		if err != nil {
			fmt.Printf("Error fetching value for key %s: %v\n", key, err)
			continue
		}
		fmt.Printf("Key: %s, Value: %v\n", key, val)
	}
}