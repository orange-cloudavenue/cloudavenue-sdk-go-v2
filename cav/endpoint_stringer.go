/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"fmt"
	"strings"
)

func (e Endpoint) String() string {
	return fmt.Sprintf("[%s] %s %s",
		e.Name,
		e.Method,
		e.PathTemplate)
}

// String returns a string representation of the Endpoint.
func (e API) String() string {
	return string(e)
}

// String returns a string representation of the version.
func (e Version) String() string {
	return string(e)
}

// String returns a string representation of the Method.
func (e Method) String() string {
	return string(e)
}

// MockPath returns the mock path for the endpoint.
func (e Endpoint) MockPath() string {
	return fmt.Sprintf("/mock/%s%s", strings.ToLower(e.Name), e.PathTemplate)
}
