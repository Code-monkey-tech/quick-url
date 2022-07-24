package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/xlab/closer"
	"os"
	"quick-url/server/configs"
	"quick-url/server/internal/server"
	"quick-url/server/internal/storage/pg"
	"quick-url/server/internal/storage/redis"
)

func main() {
	ctx := context.Background()
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
	config := configs.ReadConfig()

	pgdb := pg.New(ctx, config.PostgresURL)
	closer.Bind(func() {
		err := pgdb.Close(ctx)
		if err != nil {
			log.Error().Err(err).Str("postgres", "close").Msg("postgres close connection failed")
		}
	})

	rdb, err := redis.New(ctx, config.RedisURI, config.RedisPassword)
	if err != nil {
		log.Fatal().Err(err).Str("redis", "connection").Msg("redis connection failed")
	}
	closer.Bind(func() {
		err := rdb.Close()
		if err != nil {
			log.Error().Err(err).Str("redis", "connection").Msg("postgres close connection failed")
		}
	})

	err = server.New(ctx, config.ServerHost, config.ServerPort, pgdb, rdb)
	if err != nil {
		log.Fatal().Err(err).Str("server", "failed").Msg("start new server failed")
	}
}
