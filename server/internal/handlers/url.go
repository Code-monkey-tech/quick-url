package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jxskiss/base62"
	"net/url"
	"shrty/cache"
	"shrty/internal/storage/pg/query"
	"strconv"
	"time"
)

type Handlers struct {
	query query.Query
	cache cache.Сacher
	ctx   context.Context
}

func NewHandlers(ctx context.Context, query query.Query, cache cache.Сacher) *Handlers {
	return &Handlers{query: query, cache: cache, ctx: ctx}
}

type ShortUrlRequest struct {
	LongUrl string `json:"url"`
}

type ShortUrlResponse struct {
	ShortUrl string `json:"url"`
}

type ExpandUrlRequest struct {
	Hash string `json:"hash"`
}

type ExpandUrlResponse struct {
	LongUrl string `json:"url"`
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

// ShortenUrl create short url from long address
func (h *Handlers) ShortenUrl(fc *fiber.Ctx) error {
	sur := new(ShortUrlRequest)
	if err := fc.BodyParser(sur); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	uri, err := url.ParseRequestURI(sur.LongUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid url: %s", err).Error())
	}

	get, err := h.cache.Get(fc.Context(), sur.LongUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if get != "" {
		return fc.JSON(ShortUrlResponse{ShortUrl: get})
	}

	newID, err := h.query.NextVal(fc.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var bs62 = base62.Encode([]byte(strconv.FormatUint(*newID, 10)))

	err = h.query.ShortenUrl(
		fc.Context(), string(bs62), uri.String(), time.Now(), time.Now().Add(defaultExpireDate))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.cache.Set(fc.Context(), sur.LongUrl, string(bs62), cacheTtl); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return fc.JSON(ShortUrlResponse{ShortUrl: string(bs62)})
}

// ExpandUrl return long url from short address
func (h *Handlers) ExpandUrl(fc *fiber.Ctx) error {
	sur := new(ExpandUrlRequest)
	if err := fc.BodyParser(sur); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	get, err := h.cache.Get(fc.Context(), sur.Hash)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if get != "" {
		return fc.JSON(ShortUrlResponse{ShortUrl: get})
	}

	longUrl, err := h.query.GetLongUrl(fc.Context(), sur.Hash)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if *longUrl == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unknown short url: %s", sur))
	}

	if err := h.cache.Set(fc.Context(), sur.Hash, *longUrl, cacheTtl); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := fc.JSON(ExpandUrlResponse{LongUrl: *longUrl}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
