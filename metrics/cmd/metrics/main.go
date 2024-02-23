package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogapopp/url-shortener/metrics/lib/logger"
	"github.com/gogapopp/url-shortener/metrics/prom"
	"github.com/gogapopp/url-shortener/metrics/queue"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal(err)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			logger.Fatal(err)
		}
	}()
	logger.Infof("http: Server running at %s", ":8081")

	consumer, err := queue.NewKafkaConsumer([]string{"URL-save", "URL-get"})
	if err != nil {
		logger.Fatal(err)
	}

	prometheus := prom.NewPrometheus(logger)
	prometheus.ConsumeAndExportMetrics(consumer, []string{"URL-save", "URL-get"})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint
}
