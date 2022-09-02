package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"shrty/internal/handlers"
)

func New(ctx context.Context, port string, pgdb *pgxpool.Pool, rdb *redis.Client) error {
	config := fiber.Config{
		Prefork:      false,
		ServerHeader: "shrty-server",
		AppName:      "shrty",
	}
	handle := handlers.NewHandlers(ctx, pgdb, rdb)
	app := fiber.New(config)
	app.Post("/shorten", handle.Shorter)
	app.Get("/long", handle.Longer)
	app.Get("/", handle.HealthCheck)

	err := app.Listen(":" + port)
	if err != nil {
		return err
	}
	return nil
}
