package mock

import (
	"log/slog"
)

type OptionFunc func(*Options) error

type Options struct {
	logger *slog.Logger
}

func WithLogger(logger *slog.Logger) OptionFunc {
	return func(c *Options) error {
		if logger == nil {
			return nil
		}

		c.logger = logger
		return nil
	}
}
