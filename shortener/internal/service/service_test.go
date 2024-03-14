package service

import (
	"context"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gogapopp/url-shortener/shortener/internal/lib/logger"
	"github.com/gogapopp/url-shortener/shortener/internal/repository"
	"github.com/gogapopp/url-shortener/shortener/internal/repository/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	args := m.Called(msg, deliveryChan)
	return args.Error(0)
}

func TestServiceSave(t *testing.T) {
	tests := []struct {
		name     string
		longURL  string
		shortURL string
		wantErr  error
	}{
		{"#1 no err", "http://yandex.ru/", "http://localhost:8080/sdas", nil},
		{"#2 err", "http://yandex.ru/", "http://localhost:8080/sdas", repository.ErrURLAlreadyExists},
		{"#3 err", "http://yandex.ru/", "http://localhost:8080/sdas", repository.ErrURLAlreadyExists},
		{"#4 no err", "http://example.ru/", "http://localhost:8080/sdas", nil},
		{"#5 no err", "http://examp23le.ru/", "http://localhost:8080/sdas", nil},
		{"#6 no err", "http://example.ru/", "http://localhost:8080/sdas", nil},
		{"#7 err", "http://example.ru/", "http://localhost:8080/sdas", repository.ErrURLAlreadyExists},
	}

	repo := memory.NewRepository()
	logger, err := logger.NewLogger()
	assert.NoError(t, err)
	defer logger.Sync()
	mockProducer := new(MockProducer)
	mockProducer.On("Produce", mock.Anything, mock.Anything).Return(nil)
	service := NewService(repo, mockProducer, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Save(context.TODO(), tt.longURL, tt.shortURL)
			assert.Equal(t, tt.wantErr, err)
			mockProducer.AssertExpectations(t)
		})
	}
}

func TestServiceGet(t *testing.T) {
	tests := []struct {
		name     string
		shortURL string
		wantURL  string
		wantErr  error
	}{
		{"#1 no err", "http://example1.com", "http://yandex1.ru/", nil},
		{"#2 no err", "http://example2.com", "http://yandex2.ru/", nil},
		{"#3 no err", "http://example3.com", "http://yandex3.ru/", nil},
		{"#4 no errr", "http://example4.com", "http://example1.ru/", nil},
		{"#5 no err", "http://example5.com", "http://example2.ru/", nil},
		{"#6 no err", "http://example6.com", "http://example3.ru/", nil},
		{"#7 no err", "http://example7.com", "http://example4.ru/", nil},
		{"#8 err", "http://example8.com", "", repository.ErrURLNotExists},
		{"#9 err", "http://example9.com", "", repository.ErrURLNotExists},
	}

	repo := memory.NewRepository()
	_ = repo.Save(context.Background(), "http://yandex1.ru/", "http://example1.com")
	_ = repo.Save(context.Background(), "http://yandex2.ru/", "http://example2.com")
	_ = repo.Save(context.Background(), "http://yandex3.ru/", "http://example3.com")
	_ = repo.Save(context.Background(), "http://example1.ru/", "http://example4.com")
	_ = repo.Save(context.Background(), "http://example2.ru/", "http://example5.com")
	_ = repo.Save(context.Background(), "http://example3.ru/", "http://example6.com")
	_ = repo.Save(context.Background(), "http://example4.ru/", "http://example7.com")
	logger, err := logger.NewLogger()
	assert.NoError(t, err)
	defer logger.Sync()
	mockProducer := new(MockProducer)
	mockProducer.On("Produce", mock.Anything, mock.Anything).Return(nil)
	service := NewService(repo, mockProducer, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := service.Get(context.TODO(), tt.shortURL)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantURL, url)
			mockProducer.AssertExpectations(t)
		})
	}
}
