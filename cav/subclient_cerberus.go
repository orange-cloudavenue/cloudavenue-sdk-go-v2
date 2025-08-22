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

var _ subClientInterface = &cerberus{}

type cerberus struct {
	subclient
}

type cerberusError struct {
	Code    string `json:"code" fake:"{regex:err-[0-9]{4}}"`
	Reason  string `json:"reason" fake:"{regex:mock-[0-9]{4}}"`
	Message string `json:"message" fake:"{sentence:3,10}"`
}

var newCerberusClient = func() subClientInterface {
	return &cerberus{}
}

const cerberusVCDVersion = vmwareVCDVersion // Reusing the same version as VMware

// getID returns the unique identifier for the subclient
func (v *cerberus) getID() string {
	return string(ClientCerberus)
}

// NewClient creates a new request for the Cerberus subclient.
func (v *cerberus) newHTTPClient(ctx context.Context) (*resty.Client, error) {
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

// setCredential sets the authentication credential for the Cerberus client.
func (v *cerberus) setCredential(a auth) {
	v.credential = a
}

// setConsole sets the console for the Cerberus client.
func (v *cerberus) setConsole(c consoles.ConsoleName) {
	v.console = c
}

// Close closes the Cerberus client and releases any resources.
func (v *cerberus) close() error {
	// Close the HTTP client if it was created.
	if v.httpClient != nil {
		return v.httpClient.Close()
	}
	return nil
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
// Retries are triggered if the error message indicates that the job already exists.
func (v *cerberus) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(resp *resty.Response, err error) bool {
		if err, ok := resp.Error().(*cerberusError); ok {
			return regexCerberusJobAlreadyExists.MatchString(err.Reason) || regexCerberusJobAlreadyExists.MatchString(err.Message)
		}

		if err != nil {
			return regexCerberusJobAlreadyExists.MatchString(err.Error())
		}

		return false
	}
}
