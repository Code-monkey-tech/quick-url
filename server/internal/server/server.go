package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"shrty/api"
)

func New(ctx context.Context, port string, pgdb *pgxpool.Pool, rdb *redis.Client) error {
	config := fiber.Config{
		Prefork:      false,
		ServerHeader: "shrty-server",
		AppName:      "shrty",
	}
	handlers := api.NewHandlers(ctx, pgdb, rdb)
	app := fiber.New(config)
	app.Post("/shorten", handlers.Shorter)
	app.Get("/long", handlers.Longer)

	err := app.Listen(":" + port)
	if err != nil {
		return err
	}
	return nil
}
