package handler

import (
	"fmt"
	"net/http"

	"github.com/ak-xen/URLShortener.git/internal/db"
)

type Handlers struct {
	db *db.Repository
}

func (h Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler.ServeHTTP")
}

func NewHandler(db *db.Repository) Handlers {
	return Handlers{db: db}
}
