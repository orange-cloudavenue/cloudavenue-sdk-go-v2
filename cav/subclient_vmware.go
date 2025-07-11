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

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ SubClient = &vmware{}

const vmwareVCDVersion = "38.1"

var newVmwareClient = func() SubClient {
	return &vmware{}
}

// NewClient creates a new request for the VMware subclient.
func (v *vmware) NewHTTPClient(ctx context.Context) (*resty.Client, error) {
	// Create a new HTTP client with the base URL and headers.
	v.httpClient = httpclient.NewHTTPClient().
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

	v.httpClient.
		SetHeaders(v.credential.Headers())

	return v.httpClient, nil
}

// SetCredential sets the authentication credential for the VMware client.
func (v *vmware) SetCredential(a auth) {
	v.credential = a
}

// SetConsole sets the console for the VMware client.
func (v *vmware) SetConsole(c consoles.Console) {
	v.console = c
}

// ParseAPIError parses the API error response from the VMware client.
func (v *vmware) ParseAPIError(operation string, resp *resty.Response) *errors.APIError {
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
		}
	}

	// This is used to prevent nil pointer dereference if SetError() was not called or overrided by other object.
	return &errors.APIError{
		Operation:  operation,
		StatusCode: resp.StatusCode(),
		Message:    "Unknown error occurred",
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
	}
}
