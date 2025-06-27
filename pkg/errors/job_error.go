/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package errors

import (
	"fmt"
	"time"
)

type (
	JobError struct {
		Status   string
		Message  string
		Duration time.Duration
		Endpoint string
	}
)

// Error returns the error message for JobError.
func (e *JobError) Error() string {
	if e == nil {
		return "nil JobError"
	}
	return fmt.Sprintf("request Job error: %s (status code: %s, duration: %s, endpoint: %s)",
		e.Message, e.Status, e.Duration, e.Endpoint,
	)
}
