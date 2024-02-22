package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogapopp/url-shortener/metrics/lib/logger"
)

func main() {
	ctx := context.Background()
	_ = ctx
	logger, err := logger.NewLogger()
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal(err)
		}
	}()

	// logger.Infof("http: Server running at %s", RunAddress)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
}
