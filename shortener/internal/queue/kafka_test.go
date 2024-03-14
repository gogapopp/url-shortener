package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKafkaProducer(t *testing.T) {
	producer, err := NewKafkaProducer()
	assert.NotNil(t, producer)
	assert.NoError(t, err)
}
