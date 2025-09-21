package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ak-xen/URLShortener.git/config"
	"github.com/go-chi/chi"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	_ = cfg
	r := chi.NewRouter()
	slog.Info("Starting server")
	err = http.ListenAndServe(cfg.App.Port, r)
	if err != nil {
		return
	}

}
