package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, url string) *pgx.Conn {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("Unable to connect to database: %v\n", err)).Str("postgres", "connection")
	}
	return conn
}
