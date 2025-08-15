package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func New(addr string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr, DB: db})
}

func Ping(ctx context.Context, rdb *redis.Client) error {
	return rdb.Ping(ctx).Err()
}
