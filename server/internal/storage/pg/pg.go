package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, url string) *pgxpool.Pool {
	pc, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("Unable to parse config url: %v\n", err)).Str("postgres", "config")
	}

	pc.ConnConfig.Logger = zerologadapter.NewLogger(log.Logger)

	cp, err := pgxpool.ConnectConfig(ctx, pc)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("Unable to connect to database: %v\n", err)).Str("postgres", "connection")
	}
	return cp
}
