package query

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Queryer interface {
	GetLongUrl(ctx context.Context, hash string) (error, *string)
	NextVal(ctx context.Context) (*uint64, error)
	ShortenUrl(ctx context.Context, bs62 string, uri string, time, expDate time.Time) error
}

type Query struct {
	pool *pgxpool.Pool
}

func NewQuery(pool *pgxpool.Pool) Query {
	return Query{pool: pool}
}

// GetLongUrl return long url by hash
func (q *Query) GetLongUrl(ctx context.Context, hash string) (*string, error) {
	longUrl := new(string)
	return longUrl, q.pool.QueryRow(ctx,
		"select long_url from public.url where short_url ilike $1",
		hash).Scan(longUrl)
}

// NextVal generate new value sequence
func (q *Query) NextVal(ctx context.Context) (*uint64, error) {
	newID := new(uint64)
	return newID, q.pool.QueryRow(ctx,
		`select nextval(pg_get_serial_sequence('url','id'))`).Scan(newID)
}

// ShortenUrl create new short url
func (q *Query) ShortenUrl(
	ctx context.Context, bs62 string, uri string, time, expDate time.Time) error {
	_, err := q.pool.Exec(ctx,
		"insert into public.url (short_url, long_url, insert_date, expire_date) "+
			"values ($1, $2, $3, $4)", bs62, uri, time, expDate)
	return err
}
