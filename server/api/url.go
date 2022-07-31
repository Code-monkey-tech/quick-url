package api

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
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
	pgdb *pgx.Conn
	rdb  *redis.Client
	ctx  context.Context
}

func NewHandlers(ctx context.Context, pgdb *pgx.Conn, rdb *redis.Client) *Handlers {
	return &Handlers{pgdb: pgdb, rdb: rdb, ctx: ctx}
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
	// defaultExpireDate expire url
	defaultExpireDate = time.Hour * 24 * 30
)

// ShortUrl create short url
func (h *Handlers) ShortUrl(fc *fiber.Ctx) error {
	sur := new(ShortUrlRequest)
	if err := fc.BodyParser(sur); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	uri, err := url.ParseRequestURI(sur.ShortUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid url: %s", err).Error())
	}

	newID := new(uint64)
	err = h.pgdb.QueryRow(fc.Context(), `select nextval(pg_get_serial_sequence('url','id'))`).Scan(newID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bs62 := hex.EncodeToString([]byte(strconv.FormatUint(*newID, 10)))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = h.pgdb.Exec(fc.Context(), "insert into public.url (short_url, long_url, insert_date, expire_date) "+
		"values ($1, $2, $3, $4)", bs62, uri.String(), time.Now(), time.Now().Add(defaultExpireDate))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// TODO: hashing
	//val, err := h.DupCheck(bs62)
	//if err != nil {
	//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//}
	//
	//if val == "" {
	//	err = h.rdb.Set(h.ctx, bs62, uri, defaultTtl).Err()
	//	if err != nil {
	//		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//	}
	//	if err := fc.JSON(bs62); err != nil {
	//		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//	}
	//} else {
	//	return fc.JSON(val)
	//}

	return fc.JSON(ShortUrlResponse{ShortUrl: bs62})
}

// LongUrl get long url
func (h *Handlers) LongUrl(fc *fiber.Ctx) error {
	sur := new(LongUrlRequest)
	if err := fc.BodyParser(sur); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// TODO: hashing
	//cmd := h.rdb.Get(h.ctx, sur.LongUrl)
	//err := cmd.Err()
	//if (err != nil) && err == redis.Nil {
	//	return fiber.NewError(fiber.StatusGone, "")
	//} else if cmd.Err() != nil {
	//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//}
	//
	//str, err := cmd.Result()
	//if err != nil {
	//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//}

	shortUrl := new(string)
	err := h.pgdb.QueryRow(fc.Context(),
		"select long_url from public.url where short_url ilike $1", sur.LongUrl).Scan(shortUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if *shortUrl == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unknown short url: %s", sur))
	}

	if err := fc.JSON(LongUrlResponse{LongUrl: *shortUrl}); err != nil {
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
