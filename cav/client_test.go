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
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

func Test_NewClient(t *testing.T) {
	_, err := newMockClient()
	assert.Nil(t, err, "Error creating mock client")
}

func Test_NewClient_InvalidOrganization(t *testing.T) {
	_, err := NewClient("invalid_org")
	if err == nil {
		t.Fatal("Expected error for invalid organization, got nil")
	}
}

func TestConflictRetryStopsAfterBusyEntityLimit(t *testing.T) {
	vmwareClient := &vmware{}
	retryCount := 5
	busyEntityRetries := 0

	conflictRetry := func(resp *resty.Response, err error) bool {
		if vmwareClient.idempotentRetryCondition()(resp, err) {
			if busyEntityRetries >= retryCount {
				return false
			}

			busyEntityRetries++
			if resp != nil && resp.Request != nil {
				resp.Request.RetryCount++
			}
			return true
		}

		return false
	}

	resp := &resty.Response{
		RawResponse: &http.Response{StatusCode: http.StatusConflict},
		Request:     &resty.Request{},
	}
	err := errors.New("BUSY_ENTITY")

	for range retryCount {
		assert.True(t, conflictRetry(resp, err))
	}

	assert.False(t, conflictRetry(resp, err))
	assert.Equal(t, retryCount, busyEntityRetries)
	assert.Equal(t, retryCount, resp.Request.RetryCount)
}
