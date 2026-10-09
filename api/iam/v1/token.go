/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListTokens  = "IAM.ListTokens" //nolint:gosec
	opGetToken    = "IAM.GetToken"
	opCreateToken = "IAM.CreateToken"
	opUpdateToken = "IAM.UpdateToken"
	opDeleteToken = "IAM.DeleteToken"
)

// ListTokens lists all tokens in the organization.
func (c *Client) ListTokens(ctx context.Context) ([]*types.ModelToken, error) {
	ep := endpoints.ListTokens()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.OverrideSetResult(new(itypes.APIResponseListTokens)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListTokens, err)
	}

	tokens, ok := resp.Result().(*itypes.APIResponseListTokens)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListTokens, resp.Result())
	}

	result := make([]*types.ModelToken, len(tokens.Tokens))
	for i, t := range tokens.Tokens {
		result[i] = t.ToModel()
	}

	return result, nil
}

// GetToken retrieves a token by ID.
func (c *Client) GetToken(ctx context.Context, params types.ParamsGetToken) (*types.ModelToken, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetToken, err)
	}

	ep := endpoints.GetToken()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], params.ID),
		cav.OverrideSetResult(new(itypes.APIResponseToken)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetToken, err)
	}

	token, ok := resp.Result().(*itypes.APIResponseToken)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetToken, resp.Result())
	}

	return token.ToModel(), nil
}

// CreateToken creates a new token in the organization.
func (c *Client) CreateToken(ctx context.Context, params types.ParamsCreateToken) (*types.ModelToken, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateToken, err)
	}

	body, err := modelTokenToAPIRequest(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateToken, err)
	}

	ep := endpoints.CreateToken()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseToken)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateToken, err)
	}

	token, ok := resp.Result().(*itypes.APIResponseToken)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opCreateToken, resp.Result())
	}

	return token.ToModel(), nil
}

// UpdateToken updates an existing token.
func (c *Client) UpdateToken(ctx context.Context, params types.ParamsUpdateToken) (*types.ModelToken, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateToken, err)
	}

	body, err := modelTokenToAPIRequest(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateToken, err)
	}

	ep := endpoints.UpdateToken()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], params.ID),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseToken)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateToken, err)
	}

	token, ok := resp.Result().(*itypes.APIResponseToken)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUpdateToken, resp.Result())
	}

	return token.ToModel(), nil
}

// DeleteToken deletes a token by ID.
func (c *Client) DeleteToken(ctx context.Context, params types.ParamsDeleteToken) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opDeleteToken, err)
	}

	ep := endpoints.DeleteToken()
	_, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], params.ID),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opDeleteToken, err)
	}

	return nil
}

// modelTokenToAPIRequest converts public params to internal API request.
func modelTokenToAPIRequest(params any) (itypes.APIRequestToken, error) {
	switch p := params.(type) {
	case types.ParamsCreateToken:
		return itypes.APIRequestToken{
			Name:        p.Name,
			Description: p.Description,
			Role: itypes.APIObjectReference{
				ID:   p.RoleID,
				Name: p.RoleName,
			},
			Enabled: p.IsEnabled,
		}, nil
	case types.ParamsUpdateToken:
		return itypes.APIRequestToken{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Role: itypes.APIObjectReference{
				ID:   p.RoleID,
				Name: p.RoleName,
			},
			Enabled: p.IsEnabled,
		}, nil
	default:
		return itypes.APIRequestToken{}, fmt.Errorf("unsupported params type %T", params)
	}
}
