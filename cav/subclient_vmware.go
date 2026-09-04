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

var _ subClientInterface = &vmware{}

const vmwareVCDVersion = "39.1"

func newVmwareClient() subClientInterface {
	return &vmware{}
}

func (v *vmware) getID() string {
	return string(ClientVmware)
}

func (v *vmware) newHTTPClient(ctx context.Context) (*resty.Client, error) {
	hC := httpclient.NewHTTPClient().
		SetBaseURL(v.console.GetAPIVCDEndpoint()).
		SetHeader("Accept", "application/json;version="+vmwareVCDVersion).
		SetResultError(vmwareError{})

	if !v.credential.IsInitialized() {
		if err := v.credential.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	hC.
		SetHeaders(v.credential.Headers())

	return hC, nil
}

func (v *vmware) parseAPIError(operation string, resp *resty.Response) *errors.APIError {
	if resp == nil || resp.StatusCode() < http.StatusBadRequest {
		return nil
	}

	if err, ok := resp.ResultError().(*vmwareError); ok {
		return &errors.APIError{
			Operation:  operation,
			StatusCode: resp.StatusCode(),
			Message:    err.Message,
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
			Method:     resp.Request.Method,
			Err:        classifyStatusCode(resp.StatusCode()),
		}
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

var regexVmwareBusyEntity = regexp.MustCompile(`BUSY_ENTITY`)

// idempotentRetryCondition retries VMware busy-entity conflicts.
func (v *vmware) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(resp *resty.Response, err error) bool {
		if resp == nil || resp.StatusCode() != 409 {
			return false
		}

		if err != nil && regexVmwareBusyEntity.MatchString(err.Error()) {
			return true
		}

		return false
	}
}

func (v *vmware) ContextData(_ context.Context) ContextData {
	if v == nil || v.getCredential() == nil {
		return ContextData{}
	}

	extra := v.getCredential().getExtraData()
	return ContextData{
		OrganizationID: extra["organizationID"],
		SiteID:         extra["siteID"],
	}
}
