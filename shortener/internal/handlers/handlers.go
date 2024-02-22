package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	urlgenerator "github.com/gogapopp/url-shortener/shortener/internal/lib/url-generator"
	"github.com/gogapopp/url-shortener/shortener/internal/repository"
	"go.uber.org/zap"
)

type Handlers struct {
	service Service
	logger  *zap.SugaredLogger
}

type Service interface {
	Save(ctx context.Context, longURL, shortURL string) error
	Get(ctx context.Context, shortURL string) (string, error)
}

func NewHandlers(service Service, logger *zap.SugaredLogger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

func (h *Handlers) PostURLSaveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.PostURLSaveHandler"
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Errorf("%s: %w", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		shortURL, err := urlgenerator.GenerateShortURL()
		if err != nil {
			h.logger.Errorf("%s: %w", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		err = h.service.Save(ctx, string(body), shortURL)
		if err != nil {
			if errors.Is(err, repository.ErrURLAlreadyExists) {
				h.logger.Errorf("%s: %w", op, err)
				http.Error(w, "url already exists", http.StatusBadRequest)
				return
			}
			h.logger.Errorf("%s: %w", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, shortURL)
	}
}

func (h *Handlers) GetURLGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GettURLGetHandler"
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Errorf("%s: %w", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		longURL, err := h.service.Get(ctx, string(body))
		if err != nil {
			if errors.Is(err, repository.ErrURLNotExists) {
				h.logger.Errorf("%s: %w", op, err)
				http.Error(w, "url not exists", http.StatusBadRequest)
				return
			}
			h.logger.Errorf("%s: %w", op, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, longURL)
	}
}
