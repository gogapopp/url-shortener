package queue

import "github.com/confluentinc/confluent-kafka-go/kafka"

func NewKafkaProducer() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
	if err != nil {
		return nil, err
	}
	return p, nil
}
