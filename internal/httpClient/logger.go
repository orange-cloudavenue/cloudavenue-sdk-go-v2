/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
