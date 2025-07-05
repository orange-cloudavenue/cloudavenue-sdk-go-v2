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
	"fmt"
)

// Newf creates a new error with a formatted message.
func Newf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

var (
	// Reimplement golang "errors" package function to avoid package name conflict

	// It creates a new error with the given message.
	New = errors.New

	// As reimplements the standard library's errors.As function.
	As = errors.As

	// Is reimplements the standard library's errors.Is function.
	Is = errors.Is

	// Unwrap reimplements the standard library's errors.Unwrap function.
	Unwrap = errors.Unwrap
)
