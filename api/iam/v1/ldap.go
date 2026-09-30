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
	opTestLDAP         = "IAM.TestLDAP"
	opSyncLDAP         = "IAM.SyncLDAP"
	opSearchLDAPUsers  = "IAM.SearchLDAPUsers"
	opSearchLDAPGroups = "IAM.SearchLDAPGroups"
)

// TestLDAP tests an LDAP connection using the provided configuration.
// This is a system-level (non org-scoped) admin operation.
func (c *Client) TestLDAP(ctx context.Context, params types.ParamsTestLDAP) (*types.ModelLDAPTestResult, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opTestLDAP, err)
	}

	body := itypes.APIRequestLDAPTest{
		Config: itypes.LDAPConfigToAPIRequest(params),
	}

	ep := endpoints.TestLDAP()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.SetCustomRestyOption(setJSONHeaders),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseLDAPTestResult)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opTestLDAP, err)
	}

	result, ok := resp.Result().(*itypes.APIResponseLDAPTestResult)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opTestLDAP, resp.Result())
	}

	return result.ToModel(), nil
}

// SyncLDAP triggers an LDAP directory synchronization.
// This is a system-level (non org-scoped) admin operation.
func (c *Client) SyncLDAP(ctx context.Context, params types.ParamsSyncLDAP) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opSyncLDAP, err)
	}

	body := itypes.APIRequestLDAPSync{
		Config: itypes.LDAPConfigToAPIRequest(types.ParamsTestLDAP{}),
	}

	ep := endpoints.SyncLDAP()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.SetCustomRestyOption(setJSONHeaders),
		cav.SetBody(body),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opSyncLDAP, err)
	}

	return nil
}

// SearchLDAPUsers searches the LDAP directory for users.
// This is a system-level (non org-scoped) admin operation.
func (c *Client) SearchLDAPUsers(ctx context.Context, params types.ParamsSearchLDAP) ([]*types.ModelLDAPUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opSearchLDAPUsers, err)
	}

	ep := endpoints.SearchLDAPUsers()
	opts := []cav.EndpointRequestOption{
		cav.SetCustomRestyOption(setJSONHeaders),
		cav.OverrideSetResult(new(itypes.APIResponseLDAPUsers)),
	}

	if params.Filter != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[0], params.Filter))
	}
	if params.MaxResults != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[1], params.MaxResults))
	}
	if params.PageSize != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[2], params.PageSize))
	}
	if params.Page != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[3], params.Page))
	}

	resp, err := c.c.Do(ctx, ep, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opSearchLDAPUsers, err)
	}

	users, ok := resp.Result().(*itypes.APIResponseLDAPUsers)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opSearchLDAPUsers, resp.Result())
	}

	return users.ToModel(), nil
}

// SearchLDAPGroups searches the LDAP directory for groups.
// This is a system-level (non org-scoped) admin operation.
func (c *Client) SearchLDAPGroups(ctx context.Context, params types.ParamsSearchLDAP) ([]*types.ModelLDAPGroup, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opSearchLDAPGroups, err)
	}

	ep := endpoints.SearchLDAPGroups()
	opts := []cav.EndpointRequestOption{
		cav.SetCustomRestyOption(setJSONHeaders),
		cav.OverrideSetResult(new(itypes.APIResponseLDAPGroups)),
	}

	if params.Filter != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[0], params.Filter))
	}
	if params.MaxResults != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[1], params.MaxResults))
	}
	if params.PageSize != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[2], params.PageSize))
	}
	if params.Page != "" {
		opts = append(opts, cav.WithQueryParam(ep.QueryParams[3], params.Page))
	}

	resp, err := c.c.Do(ctx, ep, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opSearchLDAPGroups, err)
	}

	groups, ok := resp.Result().(*itypes.APIResponseLDAPGroups)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opSearchLDAPGroups, resp.Result())
	}

	return groups.ToModel(), nil
}
