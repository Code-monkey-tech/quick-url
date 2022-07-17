package api

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"net/url"
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
	ShortUrl string `json:"url"`
}

type ShortUrlResponse struct {
	ShortUrl string `json:"url"`
}

type LongUrlRequest struct {
	LongUrl string `json:"long"`
}

type LongUrlResponse struct {
	LongUrl string `json:"long"`
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

	uri, err := url.ParseRequestURI(sur.ShortUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid url: %s", err).Error())
	}

	bs62 := hex.EncodeToString([]byte(uri.String()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	val, err := h.DupCheck(bs62)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if val == "" {
		err = h.rdb.Set(h.ctx, bs62, uri, defaultTtl).Err()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if err := fc.JSON(bs62); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		return fc.JSON(val)
	}

	return nil
}

func (h *Handlers) LongUrl(fc *fiber.Ctx) error {
	sur := new(LongUrlRequest)
	if err := fc.BodyParser(sur); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	cmd := h.rdb.Get(h.ctx, sur.LongUrl)
	err := cmd.Err()
	if (err != nil) && err == redis.Nil {
		return fiber.NewError(fiber.StatusGone, "")
	} else if cmd.Err() != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	str, err := cmd.Result()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := fc.JSON(str); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// GenFreeKey checking new id in storage by pseudo-random generator
func (h *Handlers) GenFreeKey() (uint64, error) {
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

// DupCheck check duplicates
func (h *Handlers) DupCheck(hex string) (string, error) {
	cmd := h.rdb.Get(h.ctx, hex)
	err := cmd.Err()
	if (err != nil) && err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return cmd.String(), nil
}
