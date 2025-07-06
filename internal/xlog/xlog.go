/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package xlog

import (
	"log/slog"
)

// New creates a new slog.Logger with a discard handler.
// This logger will not output any logs, effectively silencing all log messages.
// It can be used in tests or when logging is not required.
var logger = slog.New(slog.DiscardHandler)

// SetGlobalLogger sets the global logger to the provided slog.Logger.
// This allows the SDK to use a custom logger for logging messages.
func SetGlobalLogger(l *slog.Logger) {
	logger = l
}

// GetGlobalLogger returns the current global logger.
// This logger can be used by the SDK for logging messages.
func GetGlobalLogger() *slog.Logger {
	return logger
}
