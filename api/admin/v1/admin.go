/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package admin

import (
	"context"
	"fmt"
	"net/http"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListAdminOrgs = "Admin.ListOrgs"
	opGetAdminOrg   = "Admin.GetOrg"
)

// ListAdminOrgs lists all organizations in the admin scope.
func (c *Client) ListAdminOrgs(ctx context.Context) ([]*types.ModelAdminOrg, error) {
	ep := endpoints.ListAdminOrgs()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.AdminOrgs)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListAdminOrgs, err)
	}

	orgs, ok := resp.Result().(*itypes.AdminOrgs)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListAdminOrgs, resp.Result())
	}

	result := make([]*types.ModelAdminOrg, len(orgs.Orgs))
	for i, o := range orgs.Orgs {
		result[i] = o.ToModel()
	}

	return result, nil
}

// GetAdminOrg retrieves an organization by ID or name in the admin scope.
func (c *Client) GetAdminOrg(ctx context.Context, params types.ParamsGetAdminOrg) (*types.ModelAdminOrg, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetAdminOrg, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	// If only the name is provided, resolve it against the admin org list so the
	// path parameter receives a value vCD accepts (name or id).
	if params.ID == "" && params.Name != "" {
		orgs, err := c.ListAdminOrgs(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s: list orgs: %w", opGetAdminOrg, err)
		}

		var matched *types.ModelAdminOrg
		for _, o := range orgs {
			if o.Name == params.Name {
				matched = o
				break
			}
		}
		if matched == nil {
			return nil, fmt.Errorf("%s: organization with name %q not found", opGetAdminOrg, params.Name)
		}

		// Use the resolved ID for the path parameter when available.
		if matched.ID != "" {
			idOrName = matched.ID
		}
	}

	ep := endpoints.GetAdminOrg()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.AdminOrg)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetAdminOrg, err)
	}

	org, ok := resp.Result().(*itypes.AdminOrg)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetAdminOrg, resp.Result())
	}

	return org.ToModel(), nil
}

// EnableOrg enables an organization by ID or name in the admin scope.
func (c *Client) EnableOrg(ctx context.Context, params types.ParamsGetAdminOrg) (*types.ModelAdminOrg, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetAdminOrg, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	ep := endpoints.GetAdminOrg()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.AdminOrg)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetAdminOrg, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	org, ok := resp.Result().(*itypes.AdminOrg)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetAdminOrg, resp.Result())
	}

	return org.ToModel(), nil
}

// DisableOrg disables an organization by ID or name in the admin scope.
func (c *Client) DisableOrg(ctx context.Context, params types.ParamsGetAdminOrg) (*types.ModelAdminOrg, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetAdminOrg, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	ep := endpoints.GetAdminOrg()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.AdminOrg)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetAdminOrg, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	org, ok := resp.Result().(*itypes.AdminOrg)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetAdminOrg, resp.Result())
	}

	return org.ToModel(), nil
}

// setXMLHeaders sets Accept and Content-Type to application/xml.
func setXMLHeaders(req *resty.Request) {
	req.SetHeader("Accept", "application/xml")
	req.SetHeader("Content-Type", "application/xml")
}
