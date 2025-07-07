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
	"errors"
	"log/slog"

	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

var _ auth = (*cloudavenueCredential)(nil)

// cloudavenueCredential implements the auth interface
// for Cloudavenue authentication using a username and password.
type cloudavenueCredential struct {
	logger       *slog.Logger
	httpC        *resty.Client
	username     string `validate:"required"`
	password     string `validate:"required"`
	bearer       string
	organization string `validate:"required"`
	console      consoles.Console
}

// cloudavenueCredentialXVmwareAccessToken is the header used to retrieve the Bearer token in the authentication process.
const cloudavenueCredentialXVmwareAccessToken = "X-VMWARE-VCLOUD-ACCESS-TOKEN" // #nosec G101

// NewCloudavenueCredential creates a new CloudavenueCredential
// with the given username and password.
var NewCloudavenueCredential = newCloudavenueCredential

func newCloudavenueCredential(c consoles.Console, organization, username, password string) (auth, error) {
	cc := &cloudavenueCredential{
		logger:       xlogger.WithGroup("auth"),
		console:      c,
		organization: organization,
		username:     username,
		password:     password,
	}

	// Validator struct doesn't work because the struct is not exported.
	if err := validators.New().Var(cc.username, "required"); err != nil {
		cc.logger.Error("Failed to validate username", "error", err)
		return nil, err
	}

	if err := validators.New().Var(cc.password, "required"); err != nil {
		cc.logger.Error("Failed to validate password", "error", err)
		return nil, err
	}

	if err := validators.New().Var(cc.organization, "required"); err != nil {
		cc.logger.Error("Failed to validate organization", "error", err)
		return nil, err
	}

	if ok := consoles.CheckOrganizationName(organization); !ok {
		cc.logger.Error("Invalid organization name", "organization", organization)
		return nil, errors.New("invalid organization name")
	}

	// Set the logger with the organization name for better context in logs.
	cc.logger = cc.logger.With("organization", cc.organization)
	cc.httpC = httpclient.NewHTTPClient().SetBaseURL(c.GetAPIVCDEndpoint())

	return cc, nil
}

// Headers returns the HTTP headers required for authentication
// using the CloudavenueCredential.
func (c *cloudavenueCredential) Headers() map[string]string {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.bearer
	return headers
}

// Refresh is a placeholder method for refreshing the authentication token.
func (c *cloudavenueCredential) Refresh(ctx context.Context) error {
	logger := c.logger.WithGroup("refresh")
	ep, err := GetEndpoint("SessionVmware", MethodPOST)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get endpoint for CreateSessionVmware", "error", err)
		return errors.New("failed to get endpoint for CreateSessionVmware: " + err.Error())
	}

	opts := []EndpointRequestOption{}
	if c.bearer != "" {
		logger.DebugContext(ctx, "Using existing bearer token for authentication")
		opts = append(opts, SetCustomRestyOption(func(r *resty.Request) {
			r.SetAuthToken(c.bearer)
		}))
	} else {
		logger.DebugContext(ctx, "Using username and password for authentication")
		opts = append(opts, SetCustomRestyOption(func(r *resty.Request) {
			r.SetBasicAuth(c.username+"@"+c.organization, c.password)
		}))
	}

	opts = append(
		opts,
		SetCustomRestyOption(func(r *resty.Request) { r.SetError(&vmwareError{}) }),
		SetCustomRestyOption(func(r *resty.Request) { r.SetURL(c.console.GetAPIVCDEndpoint()) }),
	)

	resp, err := ep.requestInternalFunc(ctx, c.httpC, ep, opts...)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to refresh session", "error", err)
		return err
	}

	if err := (&vmware{}).ParseAPIError("SessionVmware", resp); err != nil {
		c.bearer = ""
		logger.ErrorContext(ctx, "Failed to refresh session", "error", err)
		return err
	}

	logger.DebugContext(ctx, "Successfully refreshed session",
		"endpoint", ep.PathTemplate,
		"method", ep.Method,
		"status", resp.StatusCode(),
	)
	c.bearer = resp.Header().Get(cloudavenueCredentialXVmwareAccessToken)

	return nil
}

// IsInitialized checks if the CloudavenueCredential is initialized.
func (c *cloudavenueCredential) IsInitialized() bool {
	return c.bearer != ""
}
