package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ak-xen/URLShortener.git/internal/db"
)

type Handlers struct {
	db db.RepositoryInterface
}

type CreateShortURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func (h Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler.ServeHTTP")
}

func NewHandler(db *db.Repository) Handlers {
	return Handlers{db: db}
}

func (h Handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var url CreateShortURLRequest
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if !IsValidate(url.URL) {
		http.Error(w, `{"error": "URL is required"}`, http.StatusBadRequest)
		return
	}

	shortUrl := CreateShortURl(url.URL)
	_ = shortUrl

}
