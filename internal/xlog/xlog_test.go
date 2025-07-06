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
