package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/gogapopp/url-shortener/shortener/internal/handlers"
	"github.com/gogapopp/url-shortener/shortener/internal/lib/logger"
	"github.com/gogapopp/url-shortener/shortener/internal/repository/memory"
	"github.com/gogapopp/url-shortener/shortener/internal/service"
)

const RunAddress = ":8080"

func main() {
	ctx := context.Background()
	logger, err := logger.NewLogger()
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal(err)
		}
	}()
	repository := memory.NewRepository()
	service := service.NewService(repository)
	handler := handlers.NewHandlers(service, logger)
	r := chi.NewRouter()
	r.Post("/save", handler.PostURLSaveHandler())
	r.Get("/get", handler.GetURLGetHandler())
	server := http.Server{
		Addr:         RunAddress,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()
	logger.Infof("http: Server running at %s", RunAddress)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
}
