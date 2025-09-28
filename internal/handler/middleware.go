package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ak-xen/URLShortener.git/internal/db"
	"github.com/go-chi/chi/v5/middleware"
)

func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// выполняем следующий handler
			next.ServeHTTP(ww, r)

			// после выполнения логируем
			logger.Info("request handled",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
				"remote", r.RemoteAddr,
			)
		})
	}
}

func IfExistShortUrl(rep db.RepositoryInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestUrl := r.URL.Path
			answer, _ := rep.FindShortUrlInDb(r.Context(), requestUrl)
			if answer != true {
				http.Error(w, `{"error": "Not find short Url"}`, http.StatusBadRequest)
			}
			next.ServeHTTP(w, r)
		})

	}
}
