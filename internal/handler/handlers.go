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

type ToURl struct {
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

	var url ToURl
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if !IsValidate(url.URL) {
		http.Error(w, `{"error": "URL is required"}`, http.StatusBadRequest)
		return
	}

	shortUrl := CreateShortURl(url.URL)

	err := h.db.Create(r.Context(), url.URL, shortUrl)
	if err != nil {
		return
	}
	err = json.NewEncoder(w).Encode(ToURl{URL: shortUrl})
	if err != nil {
		return
	}

}

func (h Handlers) Redirect(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var shortUrl ToURl

	if err := json.NewDecoder(r.Body).Decode(&shortUrl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	fullUrl, err := h.db.Get(r.Context(), shortUrl.URL)
	if err != nil {
		http.Error(w, `{"error": "URL is not find"}`, http.StatusBadRequest)
	}

	http.Redirect(w, r, fullUrl, http.StatusFound)

}
