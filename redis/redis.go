package redis

import (
	"context"
	"fmt"
	rdb "github.com/redis/go-redis/v9"
	"time"
)

func NewRedisClient(addr string, password string, DB int) (*rdb.Client, error) {
	client := rdb.NewClient(&rdb.Options{
		Addr:     addr,
		Password: password,
		DB:       DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %s", err)
	}

	return client, nil
}
