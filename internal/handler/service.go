package handler

import (
	"crypto/rand"
	"encoding/base64"
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
