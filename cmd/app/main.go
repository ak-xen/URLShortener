package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ak-xen/URLShortener.git/config"
	"github.com/ak-xen/URLShortener.git/internal/db"
	"github.com/ak-xen/URLShortener.git/internal/handler"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	configDsn, _ := pgxpool.ParseConfig(config.GetDsn(cfg))
	pool, err := pgxpool.NewWithConfig(ctx, configDsn)
	rep := db.NewRepository(pool, logger)
	h := handler.NewHandler(rep, *cfg, logger)
	r := handler.NewRouter(h, logger, rep)
	slog.Info("Starting server")
	err = http.ListenAndServe(cfg.App.Port, r)
	if err != nil {
		return
	}

}
