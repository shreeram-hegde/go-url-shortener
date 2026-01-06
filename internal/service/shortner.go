package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/shreeram-hegde/go-url-shortener/internal/model"
	"github.com/shreeram-hegde/go-url-shortener/internal/store"
)

var ErrInvalidURL = errors.New("invalid url")

type ShortenerService struct {
	store store.Store
}

func NewShortenerService(s store.Store) *ShortenerService {
	return &ShortenerService{store: s}
}

func generateCode(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:n], nil
}

func (s *ShortenerService) Create(longURL string, expiry time.Duration) (model.URL, error) {
	if longURL == "" {
		return model.URL{}, ErrInvalidURL
	}

	code, err := generateCode(6)
	if err != nil {
		return model.URL{}, err
	}

	now := time.Now()

	u := model.URL{
		Code:      code,
		LongURL:   longURL,
		CreatedAt: now,
		ExpiresAt: now.Add(expiry),
	}

	if err := s.store.Save(u); err != nil {
		return model.URL{}, err
	}

	return u, nil
}

func (s *ShortenerService) Resolve(code string) (model.URL, error) {
	u, err := s.store.Get(code)
	if err != nil {
		return model.URL{}, err
	}

	if time.Now().After(u.ExpiresAt) {
		_ = s.store.Delete(code)
		return model.URL{}, store.ErrNotFound
	}

	return u, nil
}
