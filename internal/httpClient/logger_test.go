package httpclient

import (
	"log/slog"
	"testing"
)

func TestRestyLogger_xf(_ *testing.T) {
	logger := &restyLogger{s: slog.New(slog.DiscardHandler)}

	logger.Debugf("debug message", "key", "value")
	logger.Warnf("warn message", "key1", "value1", "key2", "value2")
	logger.Errorf("error message", "key3", "value3")
}
