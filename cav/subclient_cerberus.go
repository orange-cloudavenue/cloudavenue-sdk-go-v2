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
	"net/http"
	"regexp"

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/http-client"
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

const cerberusVCDVersion = vmwareVCDVersion

func newCerberusClient() subClientInterface {
	return &cerberus{}
}

func (v *cerberus) getID() string {
	return string(ClientCerberus)
}

func (v *cerberus) newHTTPClient(ctx context.Context) (*resty.Client, error) {
	hC := httpclient.NewHTTPClient().
		SetBaseURL(v.console.GetAPICerberusEndpoint()).
		SetHeader("Accept", "application/json;version="+cerberusVCDVersion).
		SetResultError(cerberusError{})

	if !v.credential.IsInitialized() {
		if err := v.credential.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	hC.
		SetHeaders(v.credential.Headers())

	return hC, nil
}

func (v *cerberus) parseAPIError(operation string, resp *resty.Response) *errors.APIError {
	if resp == nil || resp.StatusCode() < http.StatusBadRequest {
		return nil
	}

	if err, ok := resp.ResultError().(*cerberusError); ok {
		return &errors.APIError{
			Operation:  operation,
			StatusCode: resp.StatusCode(),
			Message:    fmt.Sprintf("%s: %s", err.Reason, err.Message),
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
			Method:     resp.Request.Method,
			Err:        classifyStatusCode(resp.StatusCode()),
		}
	}

	return &errors.APIError{
		Operation:  operation,
		StatusCode: resp.StatusCode(),
		Message:    "Unknown error occurred",
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
		Method:     resp.Request.Method,
		Err:        classifyStatusCode(resp.StatusCode()),
	}
}

// regexCerberusJobAlreadyExists matches Cerberus idempotency conflicts.
var regexCerberusJobAlreadyExists = regexp.MustCompile(`Job already exists`)

// idempotentRetryCondition retries Cerberus idempotent conflicts.
func (v *cerberus) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(resp *resty.Response, err error) bool {
		if err, ok := resp.ResultError().(*cerberusError); ok {
			return regexCerberusJobAlreadyExists.MatchString(err.Reason) || regexCerberusJobAlreadyExists.MatchString(err.Message)
		}

		if err != nil {
			return regexCerberusJobAlreadyExists.MatchString(err.Error())
		}

		return false
	}
}

func (v *cerberus) ContextData(_ context.Context) ContextData {
	return ContextData{}
}
