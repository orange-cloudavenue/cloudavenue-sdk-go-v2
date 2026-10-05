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
		SetResultError(errors.CustomerAPIErrorBody{})

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
	if resp == nil {
		return nil
	}

	if resp.StatusCode() < http.StatusBadRequest {
		return errors.CustomerAPIWAFError(
			operation,
			resp.StatusCode(),
			resp.Header().Get("Content-Type"),
			resp.String(),
			resp.Duration(),
			resp.Request.URL,
			resp.Request.Method,
		)
	}

	if err, ok := resp.ResultError().(*errors.CustomerAPIErrorBody); ok && err != nil {
		return errors.CustomerAPIStatusError(
			operation,
			resp.StatusCode(),
			err,
			unknownErrorMessage,
			resp.Duration(),
			resp.Request.URL,
			resp.Request.Method,
			classifyStatusCode(resp.StatusCode()),
		)
	}

	return &errors.APIError{
		Operation:  operation,
		StatusCode: resp.StatusCode(),
		Message:    unknownErrorMessage,
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
		Method:     resp.Request.Method,
		Err:        classifyStatusCode(resp.StatusCode()),
	}
}

// regexCerberusJobAlreadyExists matches customer job idempotency conflicts
// returned through Cerberus.
var regexCerberusJobAlreadyExists = regexp.MustCompile(`Job already exists`)

// idempotentRetryCondition retries idempotent customer job conflicts returned
// through Cerberus.
func (v *cerberus) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(resp *resty.Response, err error) bool {
		if err, ok := resp.ResultError().(*errors.CustomerAPIErrorBody); ok {
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
