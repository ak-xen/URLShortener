package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ak-xen/URLShortener.git/config"
	"github.com/ak-xen/URLShortener.git/internal/db"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	db     db.RepositoryInterface
	cfg    config.Config
	logger *slog.Logger
}

type ToURl struct {
	URL string `json:"url" validate:"required,url"`
}

func (h Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler.ServeHTTP")
}

func NewHandler(db *db.Repository, cfg config.Config, logger *slog.Logger) Handlers {
	return Handlers{db: db, cfg: cfg, logger: logger}
}

func (h Handlers) Shorten(w http.ResponseWriter, r *http.Request) {

	var url ToURl
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !IsValidate(url.URL) {
		http.Error(w, `{"error": "URL is required"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	//This func redirect user if url exsit in db
	if IfUrlExistInDb(h, w, r, url.URL) {
		return
	}

	baseUrl := h.cfg.App.BaseURL
	shortCode := CreateShortURl(url.URL)
	err := h.db.Create(r.Context(), url.URL, shortCode)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(ToURl{URL: baseUrl + "/" + shortCode})
	if err != nil {
		return
	}
	h.logger.Info("short URL created",
		"short_id", shortCode,
		"original_url", url.URL,
		"user_ip", r.RemoteAddr,
	)

}

func (h Handlers) Redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var shortUrl ToURl

	shortUrl.URL = chi.URLParam(r, "shorten_url")

	fullUrl, err := h.db.Get(r.Context(), shortUrl.URL)

	if err != nil {

		http.Error(w, `{"error": "URL is not find"}`, http.StatusBadRequest)
	}

	http.Redirect(w, r, fullUrl, http.StatusFound)

}
