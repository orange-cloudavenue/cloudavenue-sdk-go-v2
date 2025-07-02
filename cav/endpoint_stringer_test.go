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
	"testing"
)

func TestEndpoint_String(t *testing.T) {
	e := Endpoint{
		Category:     "cat",
		Version:      "v1",
		Name:         "endpointName",
		Method:       "GET",
		PathTemplate: "/path/{id}",
	}
	expected := "[cat] v1 endpointName GET /path/{id}"
	if got := e.String(); got != expected {
		t.Errorf("Endpoint.String() = %q, want %q", got, expected)
	}
}
