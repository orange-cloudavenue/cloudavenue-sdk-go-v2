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

	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	_, err := newMockClient()
	assert.Nil(t, err, "Error creating mock client")
}

func Test_NewClient_InvalidOrganization(t *testing.T) {
	// Example test case for NewClient with an invalid organization
	_, err := NewClient("invalid_org")
	if err == nil {
		t.Fatal("Expected error for invalid organization, got nil")
	}
}
