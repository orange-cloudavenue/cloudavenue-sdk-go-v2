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

	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/http-client"
	itypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ subClientInterface = &ose{}

type ose struct {
	subclient
}

func newOSEClient() subClientInterface {
	return &ose{}
}

func (o *ose) getID() string {
	return string(ClientOSE)
}

func (o *ose) newHTTPClient(ctx context.Context) (*resty.Client, error) {
	hC := httpclient.NewHTTPClient().
		SetBaseURL(o.console.GetAPIOSEEndpoint()).
		SetHeader("Accept", "application/json").
		SetResultError(oseError{})

	if !o.credential.IsInitialized() {
		if err := o.credential.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	hC.
		SetHeaders(o.credential.Headers())

	return hC, nil
}

func (o *ose) parseAPIError(operation string, resp *resty.Response) *errors.APIError {
	if resp == nil || resp.StatusCode() < http.StatusBadRequest {
		return nil
	}

	if err, ok := resp.ResultError().(*oseError); ok {
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
		Message:    "Unknown error occurred",
		Duration:   resp.Duration(),
		Endpoint:   resp.Request.URL,
		Method:     resp.Request.Method,
		Err:        classifyStatusCode(resp.StatusCode()),
	}
}

func (o *ose) idempotentRetryCondition() resty.RetryConditionFunc {
	return func(resp *resty.Response, err error) bool {
		return false
	}
}

func (o *ose) ContextData(_ context.Context) ContextData {
	if o == nil || o.getCredential() == nil {
		return ContextData{}
	}

	extra := o.getCredential().getExtraData()
	return ContextData{
		OrganizationID: extra["organizationID"],
		SiteID:         extra["siteID"],
	}
}

// GetOrganizationID retrieves the organization ID associated with the current user.
func (o *ose) GetOrganizationID(ctx context.Context) (string, error) {
	hC, err := o.newHTTPClient(ctx)
	if err != nil {
		return "", err
	}

	var tenants itypes.APIResponseListAssociatedTenants
	resp, err := hC.R().
		SetContext(ctx).
		SetQueryParam("accessible-only", "true").
		SetResult(&tenants).
		Get("/api/v1/core/associated-tenants")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return "", fmt.Errorf("failed to list associated tenants: %s", resp.Status())
	}

	if len(tenants) == 0 {
		return "", fmt.Errorf("no associated tenants found")
	}

	return tenants[0].ID, nil
}

// GetS3Credentials retrieves S3 access keys for a user in an organization.
func (o *ose) GetS3Credentials(ctx context.Context, orgID, username string) (accessKey, secretKey string, err error) {
	hC, err := o.newHTTPClient(ctx)
	if err != nil {
		return "", "", err
	}

	var creds itypes.APIResponseS3Credentials
	resp, err := hC.R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"organizationID": orgID,
			"userName":       username,
		}).
		SetResult(&creds).
		Get("/api/v1/core/tenants/{organizationID}/users/{userName}/credentials")
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusBadRequest {
		return "", "", fmt.Errorf("failed to get S3 credentials: %s", resp.Status())
	}

	return creds.AccessKey, creds.SecretKey, nil
}

type oseError struct {
	Code    string `json:"code" fake:"{regex:err-[0-9]{4}}"`
	Message string `json:"message" fake:"{sentence:3,10}"`
}
