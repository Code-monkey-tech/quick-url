package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/xlab/closer"
	"os"
	"shrty/configs"
	"shrty/internal/server"
	"shrty/internal/storage/pg"
	"shrty/internal/storage/redis"
)

// @title Tomato-timer backend
// @version 1.1
// @description API Server for Tomato-timer
// @contact.name Shrty
// @contact.url https://github.com/code-monkey-tech/shrty
// @BasePath /
func main() {
	ctx := context.Background()
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
	config := configs.ReadConfig()
	pgp := pg.New(ctx, config.PostgresURL)
	closer.Bind(pgp.Close)

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

	err = server.New(ctx, config.ServerPort, pgp, rdb)
	if err != nil {
		log.Fatal().Err(err).Str("server", "failed").Msg("start new server failed")
	}
}
