package service

import (
	"context"
)

type Service struct {
	repository Repository
}

type Repository interface {
	Save(ctx context.Context, longURL, shortURL string) error
	Get(ctx context.Context, shortURL string) (string, error)
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Save(ctx context.Context, longURL, shortURL string) error {
	const op = "service.Save"
	_ = op
	return s.repository.Save(ctx, longURL, shortURL)
}

func (s *Service) Get(ctx context.Context, shortURL string) (string, error) {
	const op = "service.Get"
	_ = op
	return s.repository.Get(ctx, shortURL)
}
