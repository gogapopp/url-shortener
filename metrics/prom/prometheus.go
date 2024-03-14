package prom

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gogapopp/url-shortener/metrics/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type Prometheus struct {
	logger *zap.SugaredLogger
}

func NewPrometheus(logger *zap.SugaredLogger) *Prometheus {
	return &Prometheus{
		logger: logger,
	}
}

var (
	receivedLinks = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_links",
		Help: "The total number of received links",
	})
	receivedLinksLength = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "received_links_length",
		Help: "The total length of received links",
	})
	savedLinks = promauto.NewCounter(prometheus.CounterOpts{
		Name: "saved_links",
		Help: "The total number of saved links",
	})
	savedLinksLength = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "saved_links_length",
		Help: "The total length of saved links",
	})
)

func (p *Prometheus) ProcessMessage(msg *kafka.Message) {
	if msg.TopicPartition.Topic == nil {
		p.logger.Errorf("message has no topic")
		return
	}
	topic := *msg.TopicPartition.Topic

	var m models.Message
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		p.logger.Errorf("error parsing message value: %w", err)
		return
	}

	switch topic {
	case "URL-get":
		receivedLinks.Add(float64(m.Value))
		receivedLinksLength.Observe(float64(m.Length))
	case "URL-save":
		savedLinks.Add(float64(m.Value))
		savedLinksLength.Observe(float64(m.Length))
	default:
		p.logger.Infof("unknown topic: %s", topic)
	}
}

func (p *Prometheus) ConsumeAndExportMetrics(consumer *kafka.Consumer, topics []string) {
	err := consumer.SubscribeTopics(topics, nil)
	if err != nil {
		p.logger.Errorf("error subscribing to topics: %w", err)
		return
	}
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			p.logger.Errorf("error reading message: %w", err)
			continue
		}
		p.ProcessMessage(msg)
	}
}
