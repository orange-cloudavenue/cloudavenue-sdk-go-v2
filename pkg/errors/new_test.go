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
	"errors"
	"testing"
)

func TestNewf(t *testing.T) {
	err := Newf("error: %s %d", "test", 42)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	expected := "error: test 42"
	if !errors.Is(err, err) || err.Error() != expected {
		t.Errorf("expected error message %q, got %q", expected, err.Error())
	}
}
