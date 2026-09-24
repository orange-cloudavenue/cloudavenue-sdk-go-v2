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
	"strings"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListAdminCatalogs = "Admin.ListCatalogs"
	opGetAdminCatalog   = "Admin.GetCatalog"
)

// ListAdminCatalogs lists all catalogs in the admin scope.
func (c *Client) ListAdminCatalogs(ctx context.Context) ([]*types.ModelAdminCatalog, error) {
	ep := endpoints.ListAdminCatalogs()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.AdminCatalogs)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListAdminCatalogs, err)
	}

	catalogs, ok := resp.Result().(*itypes.AdminCatalogs)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListAdminCatalogs, resp.Result())
	}

	result := make([]*types.ModelAdminCatalog, len(catalogs.Catalogs))
	for i, cat := range catalogs.Catalogs {
		result[i] = cat.ToModel()
	}

	return result, nil
}

// GetAdminCatalog retrieves a catalog by ID or name in the admin scope.
func (c *Client) GetAdminCatalog(ctx context.Context, params types.ParamsGetAdminCatalog) (*types.ModelAdminCatalog, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetAdminCatalog, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	// If only the name is provided, resolve it against the admin catalog list so the
	// path parameter receives a value vCD accepts (name or id).
	if params.ID == "" && params.Name != "" {
		catalogs, err := c.ListAdminCatalogs(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s: list catalogs: %w", opGetAdminCatalog, err)
		}

		var matched *types.ModelAdminCatalog
		for _, cat := range catalogs {
			if cat.Name == params.Name {
				matched = cat
				break
			}
		}
		if matched == nil {
			return nil, fmt.Errorf("%s: catalog with name %q not found", opGetAdminCatalog, params.Name)
		}

		// Use the resolved ID for the path parameter when available.
		if matched.ID != "" {
			idOrName = matched.ID
		}
	}

	ep := endpoints.GetAdminCatalog()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.AdminCatalog)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetAdminCatalog, err)
	}

	cat, ok := resp.Result().(*itypes.AdminCatalog)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetAdminCatalog, resp.Result())
	}

	return cat.ToModel(), nil
}

// CreateAdminCatalog creates a catalog in the admin scope.
func (c *Client) CreateAdminCatalog(ctx context.Context, p types.ParamsCreateAdminCatalog) (*types.ModelAdminCatalog, error) {
	const opCreateAdminCatalog = "Admin.CreateCatalog"

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateAdminCatalog, err)
	}

	storageProfiles := make([]itypes.Reference, 0, len(p.StorageProfileIDs))
	for _, id := range p.StorageProfileIDs {
		storageProfiles = append(storageProfiles, itypes.Reference{Href: id})
	}

	req := &itypes.AdminCatalogRequest{
		Name:            p.Name,
		Description:     p.Description,
		StorageProfiles: storageProfiles,
	}

	ep := endpoints.CreateAdminCatalog()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], p.OrgID),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.SetBody(req),
		cav.OverrideSetResult(new(itypes.AdminCatalog)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateAdminCatalog, err)
	}

	cat, ok := resp.Result().(*itypes.AdminCatalog)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opCreateAdminCatalog, resp.Result())
	}

	return cat.ToModel(), nil
}

// UpdateAdminCatalog updates a catalog in the admin scope.
func (c *Client) UpdateAdminCatalog(ctx context.Context, p types.ParamsUpdateAdminCatalog) (*types.ModelAdminCatalog, error) {
	const opUpdateAdminCatalog = "Admin.UpdateCatalog"

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateAdminCatalog, err)
	}

	storageProfiles := make([]itypes.Reference, 0, len(p.StorageProfileIDs))
	for _, id := range p.StorageProfileIDs {
		storageProfiles = append(storageProfiles, itypes.Reference{Href: id})
	}

	req := &itypes.AdminCatalogRequest{
		Name:            p.Name,
		Description:     p.Description,
		StorageProfiles: storageProfiles,
	}

	ep := endpoints.UpdateAdminCatalog()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], p.ID),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.SetBody(req),
		cav.OverrideSetResult(new(itypes.AdminCatalog)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateAdminCatalog, err)
	}

	cat, ok := resp.Result().(*itypes.AdminCatalog)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUpdateAdminCatalog, resp.Result())
	}

	return cat.ToModel(), nil
}

// DeleteAdminCatalog deletes a catalog by ID in the admin scope.
func (c *Client) DeleteAdminCatalog(ctx context.Context, id string) error {
	const opDeleteAdminCatalog = "Admin.DeleteCatalog"

	if id == "" {
		return fmt.Errorf("%s: id is required", opDeleteAdminCatalog)
	}
	if err := validators.New().Var(id, "urn=catalog"); err != nil {
		return fmt.Errorf("%s: id: %w", opDeleteAdminCatalog, err)
	}

	ep := endpoints.DeleteAdminCatalog()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], id),
		cav.SetCustomRestyOption(setXMLHeaders),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opDeleteAdminCatalog, err)
	}

	return nil
}

const (
	opGetAdminCatalogACL = "Admin.GetCatalogACL"
	opSetAdminCatalogACL = "Admin.SetCatalogACL"
)

// GetAdminCatalogACL retrieves the catalog ACL by ID or name in the admin scope.
func (c *Client) GetAdminCatalogACL(ctx context.Context, params types.ParamsGetAdminCatalogACL) (*types.ModelAdminCatalogACL, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetAdminCatalogACL, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	// If only the name is provided, resolve it against the admin catalog list so the
	// path parameter receives a value vCD accepts (name or id).
	if params.ID == "" && params.Name != "" {
		catalogs, err := c.ListAdminCatalogs(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s: list catalogs: %w", opGetAdminCatalogACL, err)
		}

		var matched *types.ModelAdminCatalog
		for _, cat := range catalogs {
			if cat.Name == params.Name {
				matched = cat
				break
			}
		}
		if matched == nil {
			return nil, fmt.Errorf("%s: catalog with name %q not found", opGetAdminCatalogACL, params.Name)
		}

		// Use the resolved ID for the path parameter when available.
		if matched.ID != "" {
			idOrName = matched.ID
		}
	}

	ep := endpoints.GetAdminCatalogACL()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.ControlAccessParams)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetAdminCatalogACL, err)
	}

	acl, ok := resp.Result().(*itypes.ControlAccessParams)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetAdminCatalogACL, resp.Result())
	}

	return acl.ToModel(), nil
}

// SetAdminCatalogACL sets the catalog ACL by ID in the admin scope.
func (c *Client) SetAdminCatalogACL(ctx context.Context, params types.ParamsSetAdminCatalogACL) (*types.ModelAdminCatalogACL, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opSetAdminCatalogACL, err)
	}

	req := &itypes.ControlAccessParams{
		Xmlns:               "http://www.vmware.com/vcloud/v1.5",
		IsSharedToEveryone:  params.IsSharedToEveryone,
		EveryoneAccessLevel: params.EveryoneAccessLevel,
	}

	if len(params.SharedWith) > 0 {
		settings := make([]*itypes.AccessSetting, 0, len(params.SharedWith))
		for _, item := range params.SharedWith {
			settings = append(settings, &itypes.AccessSetting{
				Subject: &itypes.LocalSubject{
					HREF: userHREFFromURN(item.UserID),
					Type: "urn:vcloud:orguser",
				},
				AccessLevel: item.AccessLevel,
			})
		}
		req.AccessSettings = &itypes.AccessSettingList{
			AccessSetting: settings,
		}
	}

	ep := endpoints.SetAdminCatalogACL()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.SetBody(req),
		cav.OverrideSetResult(new(itypes.ControlAccessParams)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opSetAdminCatalogACL, err)
	}

	acl, ok := resp.Result().(*itypes.ControlAccessParams)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opSetAdminCatalogACL, resp.Result())
	}

	return acl.ToModel(), nil
}

// userHREFFromURN extracts the UUID from a user URN (urn:vcloud:orguser:<uuid>)
// and builds the admin user HREF used by vCD access settings.
func userHREFFromURN(urn string) string {
	const prefix = "urn:vcloud:orguser:"
	if strings.HasPrefix(urn, prefix) {
		return fmt.Sprintf("https://api.cloudavenue.org/api/admin/user/%s", urn[len(prefix):])
	}
	return urn
}
