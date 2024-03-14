package service

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gogapopp/url-shortener/shortener/internal/models"
	"go.uber.org/zap"
)

var (
	urlSaveTopic = "URL-save:"
	urlGetTopic  = "URL-get:"
)

type Service struct {
	repository Repository
	producer   Producer
	logger     *zap.SugaredLogger
}

type Producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}

type Repository interface {
	Save(ctx context.Context, longURL, shortURL string) error
	Get(ctx context.Context, shortURL string) (string, error)
}

func NewService(repository Repository, producer Producer, logger *zap.SugaredLogger) *Service {
	return &Service{
		repository: repository,
		producer:   producer,
		logger:     logger,
	}
}

func (s *Service) Save(ctx context.Context, longURL, shortURL string) error {
	err := s.repository.Save(ctx, longURL, shortURL)
	if err != nil {
		return err
	}
	ks := models.KafkaSaveTopic{
		Value:  1,
		Length: len(longURL),
	}
	ksBytes, err := json.Marshal(ks)
	if err != nil {
		return err
	}
	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &urlSaveTopic, Partition: kafka.PartitionAny},
		Value:          []byte(ksBytes),
	}, nil)
	// it dosen't matter if the message was sent to kafka accurately
	if err != nil {
		s.logger.Errorf("kafka produce error: %s", urlSaveTopic)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, shortURL string) (string, error) {
	longURL, err := s.repository.Get(ctx, shortURL)
	if err != nil {
		return "", err
	}
	kg := models.KafkaGetTopic{
		Value: 1,
	}
	kgBytes, err := json.Marshal(kg)
	if err != nil {
		return "", err
	}
	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &urlGetTopic, Partition: kafka.PartitionAny},
		Value:          []byte(kgBytes),
	}, nil)
	// it dosen't matter if the message was sent to kafka accurately
	if err != nil {
		s.logger.Errorf("kafka produce error: %s", urlGetTopic)
	}
	return longURL, nil
}
