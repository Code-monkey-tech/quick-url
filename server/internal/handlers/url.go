package handlers

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jxskiss/base62"
	"net/url"
	"shrty/cache"
	"strconv"
	"time"
)

type Handlers struct {
	pgp *pgxpool.Pool
	rdb *redis.Client
	ctx context.Context
}

func NewHandlers(ctx context.Context, pgp *pgxpool.Pool, rdb *redis.Client) *Handlers {
	return &Handlers{pgp: pgp, rdb: rdb, ctx: ctx}
}

type ShortUrlRequest struct {
	LongUrl string `json:"url"`
}

type ShortUrlResponse struct {
	ShortUrl string `json:"url"`
}

type LongUrlRequest struct {
	ShortUrl string `json:"long"`
}

type LongUrlResponse struct {
	LongUrl string `json:"long"`
}

const (
	// cacheTtl default values ttl
	cacheTtl = time.Minute * 5
	// defaultExpireDate expire url, for future expiration
	defaultExpireDate = time.Hour * 24 * 30 // one month
)

// HealthCheck ...
func (h *Handlers) HealthCheck(fc *fiber.Ctx) error {
	return fc.JSON(struct {
		Status string
	}{Status: "live"})
}

// Shorter create short url from long address
func (h *Handlers) Shorter(fc *fiber.Ctx) error {
	sur := new(ShortUrlRequest)
	if err := fc.BodyParser(sur); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	uri, err := url.ParseRequestURI(sur.LongUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid url: %s", err).Error())
	}

	ch := cache.NewCache(h.rdb)
	get, err := ch.Get(fc.Context(), sur.LongUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if get != "" {
		return fc.JSON(ShortUrlResponse{ShortUrl: get})
	}

	newID := new(uint64)
	if err = h.pgp.QueryRow(fc.Context(),
		`select nextval(pg_get_serial_sequence('url','id'))`).Scan(newID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	var bs62 = base62.Encode([]byte(strconv.FormatUint(*newID, 10)))

	_, err = h.pgp.Exec(fc.Context(),
		"insert into public.url (short_url, long_url, insert_date, expire_date) "+
			"values ($1, $2, $3, $4)", string(bs62), uri.String(), time.Now(), time.Now().Add(defaultExpireDate))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := ch.Set(fc.Context(), sur.LongUrl, string(bs62), cacheTtl); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return fc.JSON(ShortUrlResponse{ShortUrl: string(bs62)})
}

// Longer return long url from short address
func (h *Handlers) Longer(fc *fiber.Ctx) error {
	sur := new(LongUrlRequest)
	if err := fc.BodyParser(sur); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	ch := cache.NewCache(h.rdb)
	get, err := ch.Get(fc.Context(), sur.ShortUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if get != "" {
		return fc.JSON(ShortUrlResponse{ShortUrl: get})
	}

	longUrl := new(string)
	err = h.pgp.QueryRow(fc.Context(),
		"select long_url from public.url where short_url ilike $1",
		sur.ShortUrl).Scan(longUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if *longUrl == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unknown short url: %s", sur))
	}

	if err := ch.Set(fc.Context(), sur.ShortUrl, *longUrl, cacheTtl); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := fc.JSON(LongUrlResponse{LongUrl: *longUrl}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
