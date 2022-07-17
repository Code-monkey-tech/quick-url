package main

import (
	"context"
	"github.com/rs/zerolog"
	"os"
	"quick-url/server/configs"
	"quick-url/server/internal/server"
	"quick-url/server/internal/storage/redis"
)

func main() {
	ctx := context.Background()
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
	config := configs.ReadConfig()

	rdb, err := redis.New(ctx, config.RedisURI, config.RedisPassword)
	if err != nil {
		log.Fatal().Err(err).Str("redis", "failed").Msg("redis connection failed")
	}

	err = server.New(ctx, config.ServerHost, config.ServerPort, rdb)
	if err != nil {
		log.Fatal().Err(err).Str("server", "failed").Msg("start new server failed")
	}
}
