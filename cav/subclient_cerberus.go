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
	"fmt"
	"regexp"

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ SubClient = &cerberus{}

type cerberus struct {
	subclient
}

type cerberusError struct {
	Code    string `json:"code" fake:"{regex:err-[0-9]{4}}"`
	Reason  string `json:"reason" fake:"{regex:mock-[0-9]{4}}"`
	Message string `json:"message" fake:"{sentence:3,10}"`
}

var newCerberusClient = func() SubClient {
	return &cerberus{}
}

const cerberusVCDVersion = vmwareVCDVersion // Reusing the same version as VMware

// NewClient creates a new request for the Cerberus subclient.
func (v *cerberus) NewHTTPClient(ctx context.Context) (*resty.Client, error) {
	v.httpClient = httpclient.NewHTTPClient().
		SetBaseURL(v.console.GetAPICerberusEndpoint()).
		SetHeader("Accept", "application/json;version="+cerberusVCDVersion).
		SetError(cerberusError{})

	if !v.credential.IsInitialized() {
		if err := v.credential.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	v.httpClient.
		SetHeaders(v.credential.Headers())

	return v.httpClient, nil
}

// SetCredential sets the authentication credential for the Cerberus client.
func (v *cerberus) SetCredential(a auth) {
	v.credential = a
}

// SetConsole sets the console for the Cerberus client.
func (v *cerberus) SetConsole(c consoles.Console) {
	v.console = c
}

// ParseAPIError parses the API error response from the Cerberus client.
func (v *cerberus) parseAPIError(operation string, resp *resty.Response) *errors.APIError {
	if resp == nil || !resp.IsError() {
		return nil
	}

	// If resp.Error() is not nil, it means an error occurred.
	// Parse the error response body.
	if err, ok := resp.Error().(*cerberusError); ok {
		return &errors.APIError{
			Operation:  operation,
			StatusCode: resp.StatusCode(),
			Message:    fmt.Sprintf("%s: %s", err.Reason, err.Message),
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

// Regexp to match the error message indicating that a job already exists.
//
//	{
//	   "code": "cf-0002",
//	   "message": "another job present on org xxxx",
//	   "reason": "Job already exists"
//	}
var regexCerberusJobAlreadyExists = regexp.MustCompile(`Job already exists`)

// idempotentRetryCondition returns a retry condition function for idempotent operations.
func (v *cerberus) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(_ *resty.Response, err error) bool {
		// Check if the error message indicates that the job already exists.
		if err != nil && regexCerberusJobAlreadyExists.MatchString(err.Error()) {
			return true // Retry if the error message indicates that the job already exists.
		}

		return false
	}
}
