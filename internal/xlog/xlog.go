package xlog

import (
	"log/slog"
	"os"
)

var New = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelError,
}))
