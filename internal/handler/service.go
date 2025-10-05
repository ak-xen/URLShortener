package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
)

func IsValidate(inputUrl string) bool {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return false
	}

	// Проверяем наличие схемы (http/https)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Проверяем наличие хоста
	if u.Host == "" {
		return false
	}

	return true

}

func CreateShortURl(inputUrl string) string {
	b := make([]byte, 4)

	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)

}

func IfUrlExistInDb(h Handlers, w http.ResponseWriter, r *http.Request, url string) bool {
	ch := make(chan string)
	fullUrlCh := make(chan string)
	if ok, _ := h.db.IsLargeUrlInDb(r.Context(), url); ok {
		go func() {
			shortCode := h.db.GetShortCode(r.Context(), url)
			ch <- shortCode
		}()

		go func() {
			fullUrl, err := h.db.Get(r.Context(), <-ch)
			if err != nil {

				http.Error(w, `{"error": "URL is not find"}`, http.StatusBadRequest)
				return
			}
			fullUrlCh <- fullUrl

		}()

		http.Redirect(w, r, <-fullUrlCh, http.StatusFound)
		return true
	}
	return false
}
