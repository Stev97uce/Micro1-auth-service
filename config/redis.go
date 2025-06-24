package config

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic("No se pudo conectar a Redis: " + err.Error())
	}

	err = rdb.Set(Ctx, "auth:ping", "pong", 10*time.Second).Err()
	if err != nil {
		panic("Error escribiendo en Redis: " + err.Error())
	}

	return rdb
}
