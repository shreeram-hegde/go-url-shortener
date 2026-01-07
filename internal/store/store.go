package store

import (
	"time"

	"github.com/shreeram-hegde/go-url-shortener/internal/model"
)

type Store interface {
	Save(u model.URL) error
	Get(code string) (model.URL, error)
	Delete(code string) error
	DeleteExpired(now time.Time) error
}
