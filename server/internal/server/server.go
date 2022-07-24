package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"quick-url/server/api"
)

func New(ctx context.Context, host, port string, pgdb *pgx.Conn, rdb *redis.Client) error {
	config := fiber.Config{
		Prefork:      false,
		ServerHeader: "quick-url-server",
		AppName:      "quick-url",
	}
	handlers := api.NewHandlers(ctx, pgdb, rdb)
	app := fiber.New(config)
	app.Post("/shorten", handlers.ShortUrl)
	app.Get("/long", handlers.LongUrl)

	err := app.Listen(host + ":" + port)
	if err != nil {
		return err
	}
	return nil
}
