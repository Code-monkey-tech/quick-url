package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"net/url"
	"quick-url/server/pkg/base62"
	"strconv"
	"time"
)

type Encoder interface {
	Encode(number uint64) string
	Decode(encoded string) (uint64, error)
}

type Handlers struct {
	rdb *redis.Client
	ctx context.Context
}

func NewHandlers(ctx context.Context, rdb *redis.Client) *Handlers {
	return &Handlers{rdb: rdb, ctx: ctx}
}

type ShortUrlRequest struct {
	Url string `json:"url"`
}

type ShortUrlResponse struct {
	Url string `json:"url"`
}

const (
	// defaultTtl default values ttl
	defaultTtl = time.Hour
	// maxGenCount max attempt gen new id
	maxGenCount = 10
)

func (h *Handlers) ShortUrl(fc *fiber.Ctx) error {
	sur := new(ShortUrlRequest)
	if err := fc.BodyParser(sur); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	uri, err := url.ParseRequestURI(sur.Url)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid url: %s", err).Error())
	}

	id, err := h.GenFreeId()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = h.rdb.Set(h.ctx, strconv.FormatUint(id, 10), uri, defaultTtl).Err()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := fc.JSON(base62.NewBase62().Encode(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// GenFreeId checking new id in storage by pseudo-random generator
func (h *Handlers) GenFreeId() (uint64, error) {
	for i := 0; i <= maxGenCount; i++ {
		id := rand.Uint64()
		err := h.rdb.Get(h.ctx, strconv.FormatUint(id, 10)).Err()
		if err == redis.Nil {
			return id, nil
		} else if err != nil {
			return 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	return 0, errors.New("max count gen exceeded")
}
