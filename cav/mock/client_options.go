/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
