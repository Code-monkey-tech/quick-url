package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) Сaches {
	return &Cache{rdb: rdb}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	cmd := c.rdb.Get(ctx, key)
	err := cmd.Err()
	if (err != nil) && err == redis.Nil {
		return "", nil
	} else if cmd.Err() != nil {
		return "", cmd.Err()
	}

	str, err := cmd.Result()
	if err != nil {
		return "", err
	}
	return str, nil
}

func (c *Cache) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	cmd := c.rdb.Set(ctx, key, value, expiration)
	err := cmd.Err()
	if err != nil {
		return err
	}

	return nil
}
