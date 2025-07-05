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
)

func TestClientError_Error(t *testing.T) {
	err := &ClientError{Message: "client error"}
	if got := err.Error(); got != "client error" {
		t.Errorf("ClientError.Error() = %q, want %q", got, "client error")
	}
	var nilErr *ClientError
	if got := nilErr.Error(); got != "nil ClientError" {
		t.Errorf("ClientError.Error() nil = %q, want %q", got, "nil ClientError")
	}
}
