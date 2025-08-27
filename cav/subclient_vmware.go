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
	"context"
	"regexp"

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ subClientInterface = &vmware{}

const vmwareVCDVersion = "38.1"

var newVmwareClient = func() subClientInterface {
	return &vmware{}
}

// getID returns the unique identifier for the subclient
func (v *vmware) getID() string {
	return string(ClientVmware)
}

// NewClient creates a new request for the VMware subclient.
func (v *vmware) newHTTPClient(ctx context.Context) (*resty.Client, error) {
	// Create a new HTTP client with the base URL and headers.
	hC := httpclient.NewHTTPClient().
		SetBaseURL(v.console.GetAPIVCDEndpoint()).
		SetHeader("Accept", "application/json;version="+vmwareVCDVersion).
		SetError(vmwareError{})

	// If the credential is not initialized, refresh it.
	// This is necessary to ensure that the client has the latest authentication token.
	if !v.credential.IsInitialized() {
		if err := v.credential.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	hC.
		SetHeaders(v.credential.Headers())

	return hC, nil
}

// Close closes the VMware client and releases any resources.
func (v *vmware) close() error {
	return nil
}

// ParseAPIError parses the API error response from the VMware client.
func (v *vmware) parseAPIError(operation string, resp *resty.Response) *errors.APIError {
	if resp == nil || !resp.IsError() {
		return nil
	}

	// If resp.Error() is not nil, it means an error occurred.
	// Parse the error response body.
	if err, ok := resp.Error().(*vmwareError); ok {
		return &errors.APIError{
			Operation:  operation,
			StatusCode: resp.StatusCode(),
			Message:    err.Message,
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
			Method:     resp.Request.Method,
		}
	}

	// This is used to prevent nil pointer dereference if SetError() was not called or overrided by other object.
	return &errors.APIError{
		Operation:  operation,
		StatusCode: resp.StatusCode(),
		Message:    "Unknown error occurred",
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
		Method:     resp.Request.Method,
	}
}

var regexVmwareBusyEntity = regexp.MustCompile(`BUSY_ENTITY`)

// idempotentRetryCondition returns a retry condition function for the VMware client.
func (v *vmware) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(resp *resty.Response, err error) bool {
		// If the response is nil or the status code is not 409, do not retry.
		if resp == nil || resp.StatusCode() != 409 {
			return false
		}

		// Check if the error message indicates that the entity is busy.
		if err != nil && regexVmwareBusyEntity.MatchString(err.Error()) {
			return true // Retry if the error message indicates that the entity is busy.
		}

		return false
	}
}
