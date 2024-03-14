package memory

import (
	"context"
	"testing"

	"github.com/gogapopp/url-shortener/shortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRepositorySave(t *testing.T) {
	tests := []struct {
		name     string
		longURL  string
		shortURL string
		wantErr  error
	}{
		{"#1 no err", "http://yandex.ru/", "http://localhost:8080/sdas", nil},
		{"#2 err", "http://yandex.ru/", "http://localhost:8080/sdas", repository.ErrURLAlreadyExists},
		{"#3 err", "http://yandex.ru/", "http://localhost:8080/sdas", repository.ErrURLAlreadyExists},
		{"#4 no err", "http://example.ru/", "http://localhost:8080/sdas", nil},
		{"#5 no err", "http://examp23le.ru/", "http://localhost:8080/sdas", nil},
		{"#6 no err", "http://example.ru/", "http://localhost:8080/sdas", nil},
		{"#7 err", "http://example.ru/", "http://localhost:8080/sdas", repository.ErrURLAlreadyExists},
	}
	repositroy := NewRepository()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repositroy.Save(context.TODO(), tt.longURL, tt.shortURL)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRepositoryGet(t *testing.T) {
	repo := NewRepository()
	_ = repo.Save(context.Background(), "http://yandex1.ru/", "http://example1.com")
	_ = repo.Save(context.Background(), "http://yandex2.ru/", "http://example2.com")
	_ = repo.Save(context.Background(), "http://yandex3.ru/", "http://example3.com")
	_ = repo.Save(context.Background(), "http://example1.ru/", "http://example4.com")
	_ = repo.Save(context.Background(), "http://example2.ru/", "http://example5.com")
	_ = repo.Save(context.Background(), "http://example3.ru/", "http://example6.com")
	_ = repo.Save(context.Background(), "http://example4.ru/", "http://example7.com")

	tests := []struct {
		name     string
		shortURL string
		wantURL  string
		wantErr  error
	}{
		{"#1 no err", "http://example1.com", "http://yandex1.ru/", nil},
		{"#2 no err", "http://example2.com", "http://yandex2.ru/", nil},
		{"#3 no err", "http://example3.com", "http://yandex3.ru/", nil},
		{"#4 no errr", "http://example4.com", "http://example1.ru/", nil},
		{"#5 no err", "http://example5.com", "http://example2.ru/", nil},
		{"#6 no err", "http://example6.com", "http://example3.ru/", nil},
		{"#7 no err", "http://example7.com", "http://example4.ru/", nil},
		{"#8 err", "http://example8.com", "", repository.ErrURLNotExists},
		{"#9 err", "http://example9.com", "", repository.ErrURLNotExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := repo.Get(context.Background(), tt.shortURL)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantURL, url)
		})
	}
}
