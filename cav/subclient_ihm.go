/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
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

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/http-client"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ subClientInterface = &ihm{}

type ihm struct{ subclient }

func newIHMClient() subClientInterface { return &ihm{} }
func (i *ihm) getID() string           { return string(ClientIHM) }

func (i *ihm) newHTTPClient(ctx context.Context) (*resty.Client, error) {
	hC := httpclient.NewHTTPClient().SetBaseURL(i.console.GetIHMEndpoint()).
		SetHeader("Accept", "application/json").SetHeader("Content-Type", "application/json").
		SetResultError(errors.CustomerAPIErrorBody{})
	if !i.credential.IsInitialized() {
		if err := i.credential.Refresh(ctx); err != nil {
			return nil, err
		}
	}
	hC.SetHeaders(i.credential.Headers())
	return hC, nil
}

func (i *ihm) parseAPIError(operation string, resp *resty.Response) *errors.APIError {
	if resp == nil || resp.StatusCode() < http.StatusBadRequest {
		return nil
	}
	message := unknownErrorMessage
	if err, ok := resp.ResultError().(*errors.CustomerAPIErrorBody); ok && err != nil {
		message = err.Message
	}
	return &errors.APIError{
		Operation: operation, StatusCode: resp.StatusCode(), Message: message,
		Duration: resp.Duration(), Endpoint: resp.Request.URL, Method: resp.Request.Method,
		Err: classifyStatusCode(resp.StatusCode()),
	}
}

func (i *ihm) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(_ *resty.Response, _ error) bool { return false }
}

func (i *ihm) ContextData(_ context.Context) ContextData { return ContextData{} }
