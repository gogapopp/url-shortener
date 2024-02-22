package memory

import (
	"context"

	"github.com/gogapopp/url-shortener/shortener/internal/repository"
)

type Repository struct {
	storage map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		storage: make(map[string]string),
	}
}

func (s *Repository) Save(_ context.Context, longURL, shortURL string) error {
	if s.urlExists(longURL) {
		return repository.ErrURLAlreadyExists
	}
	s.storage[shortURL] = longURL
	return nil
}

func (s *Repository) urlExists(value string) bool {
	for _, v := range s.storage {
		if v == value {
			return true
		}
	}
	return false
}

func (s *Repository) Get(_ context.Context, shortURL string) (string, error) {
	longURL, ok := s.storage[shortURL]
	if !ok {
		return "", repository.ErrURLNotExists
	}
	return longURL, nil
}
