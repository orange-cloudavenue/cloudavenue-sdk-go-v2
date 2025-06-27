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
	"testing"
	"time"
)

func TestJobError_Error(t *testing.T) {
	apiErr := &JobError{
		Status:   "OK",
		Message:  "internal error",
		Duration: 2 * time.Second,
		Endpoint: "/test",
	}
	want := "request Job error: internal error (status: OK, duration: 2s, endpoint: /test)"
	if got := apiErr.Error(); got != want {
		t.Errorf("JobError.Error() = %q, want %q", got, want)
	}
	var nilErr *JobError
	if got := nilErr.Error(); got != "nil JobError" {
		t.Errorf("JobError.Error() nil = %q, want %q", got, "nil JobError")
	}
}
