package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {
	Create(ctx context.Context, url string, shortURL string) error
	Get(ctx context.Context, shortURL string) (string, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, url string, shortURL string) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO urls (shortURL) VALUES ($1, $2)", url, shortURL)
	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) Get(ctx context.Context, shortURL string) (string, error) {
	var url string
	err := r.pool.QueryRow(ctx, `SELECT url FROM urls WHERE url=$1`, shortURL).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}
