/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package httpclient

import "testing"

func Test_NewHTTPClient(t *testing.T) {
	client := NewHTTPClient()
	if client == nil {
		t.Error("NewHTTPClient() returned nil")
	}
	if client.IsDebug() {
		t.Error("NewHTTPClient() should not be in debug mode by default")
	}
}
