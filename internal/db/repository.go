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
	IsShortUrlInDb(ctx context.Context, shortURL string) (bool, error)
	IsLargeUrlInDb(ctx context.Context, originalUrl string) (bool, error)
	GetShortCode(ctx context.Context, originalUrl string) string
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
	if err != nil {
		r.logger.Info("err query", '\n', "err", err)

		return "", err
	}
	return url, nil
}

func (r *Repository) IsShortUrlInDb(ctx context.Context, shortUrl string) (bool, error) {
	var ans bool
	findUrl := strings.Split(shortUrl, "/")[2]
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code=$1)", findUrl).Scan(&ans)
	if err != nil {
		return false, err
	}
	return ans, nil

}

func (r *Repository) IsLargeUrlInDb(ctx context.Context, originalUrl string) (bool, error) {
	var ans bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM urls WHERE original_url=$1)", originalUrl).Scan(&ans)
	if err != nil {
		return false, err
	}
	return ans, nil

}

func (r *Repository) GetShortCode(ctx context.Context, originalUrl string) string {
	var shortCode string
	err := r.pool.QueryRow(ctx, "SELECT short_code FROM urls WHERE original_url=$1 ", originalUrl).Scan(&shortCode)
	if err != nil {
		r.logger.Info("Error query DB GetShortCode")
	}
	return shortCode
}
