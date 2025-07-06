package httpclient

import (
	"log/slog"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"
)

var _ resty.Logger = &restyLogger{}

type restyLogger struct {
	s *slog.Logger
}

func (l *restyLogger) Debugf(msg string, keysAndValues ...interface{}) {
	l.s.Debug(msg, keysAndValues...)
}

func (l *restyLogger) Warnf(msg string, keysAndValues ...interface{}) {
	l.s.Warn(msg, keysAndValues...)
}

func (l *restyLogger) Errorf(msg string, keysAndValues ...interface{}) {
	l.s.Error(msg, keysAndValues...)
}

var logger = func() resty.Logger {
	gLogger := xlog.GetGlobalLogger()

	x := &restyLogger{
		s: gLogger,
	}
	return x
}
