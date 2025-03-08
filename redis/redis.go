package redis

import (
	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string, password string, DB int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       DB,
	})

	return client
}
