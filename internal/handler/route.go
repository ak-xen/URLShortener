package handler

import (
	"log/slog"

	"github.com/ak-xen/URLShortener.git/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter собирает все маршруты
func NewRouter(h Handlers, logger *slog.Logger, repository db.RepositoryInterface) *chi.Mux {
	r := chi.NewRouter()
	middlewareLogger := RequestLogger(logger)
	// Подключаем стандартные middleware
	r.Use(middlewareLogger)
	r.Use(middleware.Recoverer)

	// Разделяем маршруты по группам
	r.Route("/api", func(r chi.Router) {
		r.Post("/v1", h.Shorten)
		r.With(IfExistShortUrl(repository)).Get("/{shorten_url}", h.Redirect)
	})

	return r
}
