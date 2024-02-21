package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogapopp/url-shortener/shortener/internal/lib/logger"
	"github.com/gogapopp/url-shortener/shortener/internal/repository"
	"github.com/gogapopp/url-shortener/shortener/internal/repository/memory"
	"github.com/gogapopp/url-shortener/shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestPostHandlers(t *testing.T) {
	tests := []struct {
		name     string
		longURL  string
		wantCode int
		wantErr  error
	}{
		{"#1 no err", "http://yandex.ru/", http.StatusCreated, nil},
		{"#2 err", "http://yandex.ru/", http.StatusBadRequest, repository.ErrURLAlreadyExists},
		{"#3 err", "http://yandex.ru/", http.StatusBadRequest, repository.ErrURLAlreadyExists},
		{"#4 no err", "http://yandex111.ru/", http.StatusCreated, nil},
		{"#5 no err", "http://yandex222.ru/", http.StatusCreated, nil},
	}

	shortURLs := make(map[int]string)
	logger, err := logger.NewLogger()
	assert.NoError(t, err)
	defer logger.Sync()
	repo := memory.NewRepository()
	service := service.NewService(repo)
	handlers := NewHandlers(service, logger)

	for i, tt := range tests {
		t.Run("Save: "+tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/save", bytes.NewBufferString(tt.longURL))
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.PostURLSaveHandler())
			handler.ServeHTTP(rr, req)
			assert.NotNil(t, rr.Body.String())
			assert.Equal(t, tt.wantCode, rr.Code)
			if tt.wantErr == nil {
				shortURLs[i] = rr.Body.String()
			}
		})
	}

	for i, tt := range tests {
		t.Run("Get: "+tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/get", bytes.NewBufferString(shortURLs[i]))
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.GetURLGetHandler())
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantCode, rr.Code)
			if tt.wantErr == nil {
				assert.Equal(t, tt.longURL, rr.Body.String())
			}
		})
	}
}
