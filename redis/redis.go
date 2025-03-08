package redis

import (
	rdb "github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string, password string, DB int) *rdb.Client {
	client := rdb.NewClient(&rdb.Options{
		Addr:     addr,
		Password: password,
		DB:       DB,
	})

	return client
}
