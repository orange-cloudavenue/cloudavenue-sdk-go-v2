package xlog

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetGlobalLogger(t *testing.T) {
	// Create a new logger with a discard handler
	newLogger := slog.New(slog.DiscardHandler)
	SetGlobalLogger(newLogger)

	got := GetGlobalLogger()
	assert.NotNil(t, got, "GetGlobalLogger() should not return nil")
	assert.Equal(t, newLogger, got, "GetGlobalLogger() should return the logger set by SetGlobalLogger")
}

func TestDefaultGlobalLoggerIsNotNil(t *testing.T) {
	logger := GetGlobalLogger()
	assert.NotNil(t, logger, "Default global logger should not be nil")
}
