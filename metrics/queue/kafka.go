package queue

import "github.com/confluentinc/confluent-kafka-go/kafka"

func NewKafkaConsumer(consumerTopics []string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "metrics",
	})
	if err != nil {
		return nil, err
	}

	err = consumer.SubscribeTopics(consumerTopics, nil)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
