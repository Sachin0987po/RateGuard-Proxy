package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"fmt"
	"encoding/json"
	"time"
)


var Client *redis.Client
var Ctx = context.Background()

func InitializeRedisClient() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Ping the Redis server to check if it's accessible
	pong, err := Client.Ping(Ctx).Result()
	fmt.Println(pong)

	if err != nil {
		return err
	}

	return nil
}


func GetDataFromRedis(key string) (RateLimiter, error) {
	rateLimiter := RateLimiter{}
	data, err := Client.Get(Ctx, key).Result()
	json.Unmarshal([]byte(data), &rateLimiter)
	return rateLimiter, err
}

func SetDataInRedis(rateLimiter RateLimiter, key string) {
	data, err := json.Marshal(rateLimiter)
	if err != nil {
		panic(err)
	}
	err = Client.Set(Ctx, key, data, 60*time.Second).Err()
	if err != nil {
		panic(err)
	}
}