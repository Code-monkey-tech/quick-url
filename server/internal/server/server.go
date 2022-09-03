package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4/pgxpool"
	"shrty/cache"
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
		AllowOrigins: "http://localhost:8080, http://localhost:8080/,",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Post(shortenUrlPath, handle.ShortenUrl)
	app.Get(expandUrlPath, handle.ExpandUrl)
	app.Get(liveUrlPath, handle.HealthCheck)

	err := app.Listen(":" + port)
	if err != nil {
		return err
	}
	return nil
}
