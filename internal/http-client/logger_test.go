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
	"testing"
)

func TestRestyLogger_xf(_ *testing.T) {
	logger := &restyLogger{s: slog.New(slog.DiscardHandler)}

	logger.Debugf("debug message", "key", "value")
	logger.Warnf("warn message", "key1", "value1", "key2", "value2")
	logger.Errorf("error message", "key3", "value3")
}
