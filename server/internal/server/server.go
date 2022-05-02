package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"quick-url/server/api"
)

func New(ctx context.Context, host, port string, rdb *redis.Client) error {
	config := fiber.Config{
		Prefork:      false,
		ServerHeader: "quick-url-server",
		AppName:      "quick-url",
	}
	handlers := api.NewHandlers(ctx, rdb)
	app := fiber.New(config)
	app.Post("/shorten", handlers.ShortUrl)
	app.Get("/long", func(c *fiber.Ctx) error {
		return c.SendString("long")
	})

	err := app.Listen(host + ":" + port)
	if err != nil {
		return err
	}
	return nil
}
