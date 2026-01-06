package store

import "github.com/shreeram-hegde/go-url-shortener/internal/model"

type Store interface {
	Save(u model.URL) error
	Get(code string) (model.URL, error)
	Delete(code string) error
}
