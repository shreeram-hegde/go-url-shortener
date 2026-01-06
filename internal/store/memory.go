package store

import (
	"errors"
	"sync"

	"github.com/shreeram-hegde/go-url-shortener/internal/model"
)

var ErrNotFound = errors.New("url not found")

type MemoryStore struct {
	mu   sync.RWMutex         //This is needed for concurency as we are gonna use Goroutines
	data map[string]model.URL //Here URL is a struct that has both longURL and short code version of the url
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]model.URL),
	}
}

func (m *MemoryStore) Save(u model.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[u.Code] = u // we get the code itself in the passed object
	return nil
}

func (m *MemoryStore) Get(Code string) (model.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	URL, ok := m.data[Code]

	if !ok {
		return model.URL{}, ErrNotFound
	}
	return URL, nil
}

func (m *MemoryStore) Delete(code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, code)
	return nil
}
