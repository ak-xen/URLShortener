package db

import (
	"context"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {
	Create(ctx context.Context, url string, shortURL string) error
	Get(ctx context.Context, shortURL string) (string, error)
	FindShortUrlInDb(ctx context.Context, shortURL string) (bool, error)
}

type Repository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewRepository(pool *pgxpool.Pool, logger *slog.Logger) *Repository {
	return &Repository{pool: pool, logger: logger}
}

func (r *Repository) Create(ctx context.Context, url string, shortURL string) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO urls (original_url, short_code) VALUES ($1, $2)", url, shortURL)
	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) Get(ctx context.Context, shortURL string) (string, error) {
	var url string
	err := r.pool.QueryRow(ctx, `SELECT original_url FROM urls WHERE short_code=$1`, shortURL).Scan(&url)
	r.logger.Info("err query", '\n', "err", err)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (r *Repository) FindShortUrlInDb(ctx context.Context, shortUrl string) (bool, error) {
	var ans bool
	findUrl := strings.Split(shortUrl, "/")[2]
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code=$1)", findUrl).Scan(&ans)
	if err != nil {
		return false, err
	}
	return ans, nil

}
