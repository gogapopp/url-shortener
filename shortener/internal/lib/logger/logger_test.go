package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}
