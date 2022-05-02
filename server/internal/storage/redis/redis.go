package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, uri, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	log.Info().Str("redis", pong)

	return rdb, nil
}
