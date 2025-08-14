/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcgroup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewClient_ClientNil(t *testing.T) {
	c, err := New(nil)
	assert.Nil(t, c, "Expected nil client when input is nil")
	assert.Error(t, err, "Expected error when input is nil")
}
