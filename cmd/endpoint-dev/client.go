package main

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
)

var logger *slog.Logger

func newClient() (cav.Client, error) {
	handler := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})
	logLevel := log.InfoLevel

	switch loggerLevel {
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "warn":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	default:
		return nil, errors.New("invalid logger level")
	}

	handler.SetLevel(logLevel)
	logger = slog.New(handler)

	// Check if the client is a mock client
	if mockFlag {
		logger.Info("Using mock client")
		return mock.NewClient(
			mock.WithLogger(logger),
		)
	}

	path := configFile
	if path == "" {
		path = defaultConfigPath()
	}

	config, err := loadConfig(path)
	if err != nil {
		return nil, err
	}

	logger.Info("Using real client")
	return cav.NewClient(
		config.Organization,
		cav.WithLogger(logger),
		cav.WithCloudAvenueCredential(config.Username, config.Password),
	)
}
