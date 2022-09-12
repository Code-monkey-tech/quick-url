package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/gofiber/swagger"
	"github.com/jackc/pgx/v4/pgxpool"
	"shrty/cache"
	_ "shrty/docs"
	"shrty/internal/handlers"
	"shrty/internal/storage/pg/query"
)

const (
	expandUrlPath  = "/expand"
	shortenUrlPath = "/shorten"
	liveUrlPath    = "/live"
)

func New(ctx context.Context, port string, pgdb *pgxpool.Pool, rdb *redis.Client) error {
	config := fiber.Config{
		Prefork:      false,
		ServerHeader: "shrty-server",
		AppName:      "shrty",
	}
	handle := handlers.NewHandlers(ctx, query.NewQuery(pgdb), cache.NewCache(rdb))
	app := fiber.New(config)
	app.Use(cors.New(cors.Config{
		AllowMethods: "POST, GET, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post(shortenUrlPath, handle.ShortenUrl)
	app.Get(expandUrlPath, handle.ExpandUrl)
	app.Get(liveUrlPath, handle.HealthCheck)
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://localhost:" + port + "/swagger/doc.json",
		DeepLinking: false,
	}))

	err := app.Listen(":" + port)
	if err != nil {
		return err
	}
	return nil
}
