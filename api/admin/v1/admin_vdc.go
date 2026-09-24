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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListAdminVDCs = "Admin.ListVDCs"
	opGetAdminVDC   = "Admin.GetVDC"
)

// ListAdminVDCs lists all VDCs in the admin scope.
func (c *Client) ListAdminVDCs(ctx context.Context) ([]*types.ModelAdminVDC, error) {
	ep := endpoints.ListAdminVDCs()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.AdminVDCs)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListAdminVDCs, err)
	}

	vdcs, ok := resp.Result().(*itypes.AdminVDCs)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListAdminVDCs, resp.Result())
	}

	result := make([]*types.ModelAdminVDC, len(vdcs.VDCs))
	for i, v := range vdcs.VDCs {
		result[i] = v.ToModel()
	}

	return result, nil
}

// GetAdminVDC retrieves a VDC by ID or name in the admin scope.
func (c *Client) GetAdminVDC(ctx context.Context, params types.ParamsGetAdminVDC) (*types.ModelAdminVDC, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetAdminVDC, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	// If only the name is provided, resolve it against the admin VDC list so the
	// path parameter receives a value vCD accepts (name or id).
	if params.ID == "" && params.Name != "" {
		vdcs, err := c.ListAdminVDCs(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s: list VDCs: %w", opGetAdminVDC, err)
		}

		var matched *types.ModelAdminVDC
		for _, v := range vdcs {
			if v.Name == params.Name {
				matched = v
				break
			}
		}
		if matched == nil {
			return nil, fmt.Errorf("%s: VDC with name %q not found", opGetAdminVDC, params.Name)
		}

		// Use the resolved ID for the path parameter when available.
		if matched.ID != "" {
			idOrName = matched.ID
		}
	}

	ep := endpoints.GetAdminVDC()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.AdminVDC)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetAdminVDC, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	vdc, ok := resp.Result().(*itypes.AdminVDC)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetAdminVDC, resp.Result())
	}

	return vdc.ToModel(), nil
}
