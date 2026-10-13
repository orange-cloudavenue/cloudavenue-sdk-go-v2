/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"context"
	"fmt"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const opCreateOAuthClient = "IAM.CreateOAuthClient"

const opCreateOAuthToken = "IAM.CreateOAuthToken" //nolint: gosec

// CreateOAuthClient registers an OAuth client for the specified tenant.
func (c *Client) CreateOAuthClient(ctx context.Context, params types.ParamsCreateOAuthClient) (*types.ModelOAuthClient, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateOAuthClient, err)
	}
	body := itypes.APIRequestOAuthClient{
		ClientName: params.ClientName, SoftwareID: params.SoftwareID,
		SoftwareVersion: params.SoftwareVersion, ClientURI: params.ClientURI, Scope: params.Scope,
	}
	ep := endpoints.CreateOAuthClient()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.TenantID), cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseOAuthClient)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateOAuthClient, err)
	}
	result, ok := resp.Result().(*itypes.APIResponseOAuthClient)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opCreateOAuthClient, resp.Result())
	}
	return &types.ModelOAuthClient{
		ClientName:              result.ClientName,
		ClientID:                result.ClientID,
		GrantTypes:              result.GrantTypes,
		TokenEndpointAuthMethod: result.TokenEndpointAuthMethod,
		ClientURI:               result.ClientURI,
		SoftwareID:              result.SoftwareID,
		SoftwareVersion:         result.SoftwareVersion,
		Scope:                   result.Scope,
	}, nil
}

// CreateOAuthToken exchanges a tenant OAuth assertion or device code for an access token.
func (c *Client) CreateOAuthToken(ctx context.Context, params types.ParamsCreateOAuthToken) (*types.ModelOAuthToken, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateOAuthToken, err)
	}
	ep := endpoints.CreateOAuthToken()
	resp, err := c.c.Do(
		ctx, ep,
		cav.WithPathParam(ep.PathParams[0], params.TenantID),
		cav.SetCustomRestyOption(func(req *resty.Request) {
			form := map[string]string{
				"client_id":  params.ClientID,
				"grant_type": params.GrantType,
			}
			if params.Assertion != "" {
				form["assertion"] = params.Assertion
			}
			if params.DeviceCode != "" {
				form["device_code"] = params.DeviceCode
			}
			req.SetFormData(form)
		}),
		cav.OverrideSetResult(new(itypes.APIResponseOAuthToken)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateOAuthToken, err)
	}
	result, ok := resp.Result().(*itypes.APIResponseOAuthToken)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opCreateOAuthToken, resp.Result())
	}
	return &types.ModelOAuthToken{
		AccessToken: result.AccessToken, TokenType: result.TokenType,
		ExpiresIn: result.ExpiresIn, RefreshToken: result.RefreshToken,
	}, nil
}
