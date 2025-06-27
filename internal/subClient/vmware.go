/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package subclient

import (
	"context"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ Client = &vmware{}

var NewVmwareClient = func() Client {
	return &vmware{}
}

// NewClient creates a new request for the VMware subclient.
func (v *vmware) NewHTTPClient(ctx context.Context) (*resty.Client, error) {
	// If the credential is not initialized, refresh it.
	// This is necessary to ensure that the client has the latest authentication token.
	if !v.credential.IsInitialized() {
		if err := v.credential.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	// Create a new HTTP client with the base URL and headers.
	return httpclient.NewHTTPClient().
		SetBaseURL(v.console.GetAPIVCDEndpoint()).
		SetHeaders(v.credential.Headers()).
		SetHeader("Accept", "application/json;version="+auth.VDCVersion), nil
}

// SetCredential sets the authentication credential for the VMware client.
func (v *vmware) SetCredential(a auth.Auth) {
	v.credential = a
}

// SetConsole sets the console for the VMware client.
func (v *vmware) SetConsole(c consoles.Console) {
	v.console = c
}

// ParseAPIError parses the API error response from the VMware client.
func (v *vmware) ParseAPIError(resp *resty.Response) *errors.APIError {
	if resp == nil || !resp.IsError() {
		return nil
	}

	// Documentation for the VCD API error response:
	// https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/types/ErrorType.html

	// If resp.Error() is not nil, it means an error occurred.
	// Parse the error response body.
	if vmwErr, ok := resp.Error().(*VmwareError); ok {
		return &errors.APIError{
			StatusCode:    resp.StatusCode(),
			StatusMessage: vmwErr.StatusMessage,
			Message:       vmwErr.Message,
			Duration:      resp.Duration(),
			Endpoint:      resp.Request.URL,
		}
	}

	return nil
}
