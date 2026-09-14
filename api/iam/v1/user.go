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
	"net/http"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

const (
	opListUsers       = "IAM.ListUsers"
	opGetUser         = "IAM.GetUser"
	opCreateLocalUser = "IAM.CreateLocalUser"
	opCreateSAMLUser  = "IAM.CreateSAMLUser"
	opUpdateUser      = "IAM.UpdateUser"
	opDeleteUser      = "IAM.DeleteUser"
	opEnableUser      = "IAM.EnableUser"
	opDisableUser     = "IAM.DisableUser"
	opUnlockUser      = "IAM.UnlockUser"
	opChangePassword  = "IAM.ChangePassword"
)

// withOrgID injects the organization ID from request context into the given path param.
func withOrgID(pp cav.PathParam) cav.EndpointRequestOption {
	return func(endpoint *cav.Endpoint, req *resty.Request) error {
		cd := cav.GetExtraDataFromContext(req.Context())
		if cd.OrganizationID == "" {
			return errors.New("organization ID not found in context")
		}
		return cav.WithPathParam(pp, cd.OrganizationID)(endpoint, req)
	}
}

// setXMLHeaders sets Accept and Content-Type to application/xml.
func setXMLHeaders(req *resty.Request) {
	req.SetHeader("Accept", "application/xml")
	req.SetHeader("Content-Type", "application/xml")
}

// iamUserToModel converts internal itypes.User to public ModelUser.
func iamUserToModel(u itypes.User) *ModelUser {
	return &ModelUser{
		ID:              u.ID,
		Name:            u.Name,
		FullName:        u.FullName,
		EmailAddress:    u.EmailAddress,
		Telephone:       u.Telephone,
		Description:     u.Description,
		IsEnabled:       u.IsEnabled,
		DeployedVmQuota: u.DeployedVmQuota,
		StoredVmQuota:   u.StoredVmQuota,
		ProviderType:    u.ProviderType,
		RoleName:        u.Role.Name,
		RoleHref:        u.Role.Href,
	}
}

// modelUserToIAMXML converts public params to internal itypes.User for XML marshaling.
func modelUserToIAMXML(params any) (itypes.User, error) {
	var u itypes.User
	switch p := params.(type) {
	case ParamsCreateLocalUser:
		u = itypes.User{
			Name:            p.Name,
			Password:        p.Password,
			FullName:        p.FullName,
			EmailAddress:    p.EmailAddress,
			Telephone:       p.Telephone,
			Description:     p.Description,
			IsEnabled:       p.IsEnabled,
			DeployedVmQuota: p.DeployedVmQuota,
			StoredVmQuota:   p.StoredVmQuota,
			ProviderType:    "INTEGRATED",
			Role: itypes.Reference{
				Name: p.RoleName,
			},
		}
	case ParamsCreateSAMLUser:
		u = itypes.User{
			Name:         p.Name,
			FullName:     p.FullName,
			EmailAddress: p.EmailAddress,
			Telephone:    p.Telephone,
			Description:  p.Description,
			IsEnabled:    p.IsEnabled,
			ProviderType: "SAML",
			Role: itypes.Reference{
				Name: p.RoleName,
			},
		}
	case ParamsUpdateUser:
		isEnabled := false
		if p.IsEnabled != nil {
			isEnabled = *p.IsEnabled
		}
		u = itypes.User{
			Name:            p.Name,
			Password:        p.Password,
			FullName:        p.FullName,
			EmailAddress:    p.EmailAddress,
			Telephone:       p.Telephone,
			Description:     p.Description,
			IsEnabled:       isEnabled,
			DeployedVmQuota: p.DeployedVmQuota,
			StoredVmQuota:   p.StoredVmQuota,
			Role: itypes.Reference{
				Name: p.RoleName,
			},
		}
	default:
		return itypes.User{}, fmt.Errorf("unsupported params type %T", params)
	}
	return u, nil
}

// ListUsers lists all users in the organization.
func (c *Client) ListUsers(ctx context.Context) ([]*ModelUser, error) {
	ep := endpoints.ListUsers()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.Users)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListUsers, err)
	}

	users, ok := resp.Result().(*itypes.Users)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListUsers, resp.Result())
	}

	result := make([]*ModelUser, len(users.Users))
	for i, u := range users.Users {
		result[i] = iamUserToModel(u)
	}

	return result, nil
}

// GetUser retrieves a user by ID or name.
func (c *Client) GetUser(ctx context.Context, params ParamsGetUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetUser, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	ep := endpoints.GetUser()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.User)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetUser, err)
	}

	user, ok := resp.Result().(*itypes.User)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetUser, resp.Result())
	}

	return iamUserToModel(*user), nil
}

// CreateLocalUser creates a new local user in the organization.
func (c *Client) CreateLocalUser(ctx context.Context, params ParamsCreateLocalUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateLocalUser, err)
	}

	body, err := modelUserToIAMXML(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateLocalUser, err)
	}

	ep := endpoints.CreateUser()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.User)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateLocalUser, err)
	}

	user, ok := resp.Result().(*itypes.User)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opCreateLocalUser, resp.Result())
	}

	return iamUserToModel(*user), nil
}

// CreateSAMLUser creates a new SAML user in the organization.
func (c *Client) CreateSAMLUser(ctx context.Context, params ParamsCreateSAMLUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateSAMLUser, err)
	}

	body, err := modelUserToIAMXML(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateSAMLUser, err)
	}

	ep := endpoints.CreateUser()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.User)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateSAMLUser, err)
	}

	user, ok := resp.Result().(*itypes.User)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opCreateSAMLUser, resp.Result())
	}

	return iamUserToModel(*user), nil
}

// UpdateUser updates an existing user.
func (c *Client) UpdateUser(ctx context.Context, params ParamsUpdateUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateUser, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	body, err := modelUserToIAMXML(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateUser, err)
	}

	ep := endpoints.UpdateUser()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.User)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateUser, err)
	}

	user, ok := resp.Result().(*itypes.User)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUpdateUser, resp.Result())
	}

	return iamUserToModel(*user), nil
}

// DeleteUser deletes a user by ID or name.
func (c *Client) DeleteUser(ctx context.Context, params ParamsDeleteUser) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opDeleteUser, err)
	}

	ep := endpoints.DeleteUser()
	_, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], params.ID),
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("%t", params.TakeOwnership)),
		cav.SetCustomRestyOption(setXMLHeaders),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opDeleteUser, err)
	}

	return nil
}

// EnableUser enables a user by ID or name.
func (c *Client) EnableUser(ctx context.Context, params ParamsEnableUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opEnableUser, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	ep := endpoints.EnableUser()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.User)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opEnableUser, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	user, ok := resp.Result().(*itypes.User)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opEnableUser, resp.Result())
	}

	return iamUserToModel(*user), nil
}

// DisableUser disables a user by ID or name.
func (c *Client) DisableUser(ctx context.Context, params ParamsDisableUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opDisableUser, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	ep := endpoints.DisableUser()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.User)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opDisableUser, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	user, ok := resp.Result().(*itypes.User)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opDisableUser, resp.Result())
	}

	return iamUserToModel(*user), nil
}

// UnlockUser unlocks a user by ID or name.
func (c *Client) UnlockUser(ctx context.Context, params ParamsUnlockUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUnlockUser, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	ep := endpoints.UnlockUser()
	resp, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.OverrideSetResult(new(itypes.User)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUnlockUser, err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	user, ok := resp.Result().(*itypes.User)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUnlockUser, resp.Result())
	}

	return iamUserToModel(*user), nil
}

// ChangePassword changes a user's password.
func (c *Client) ChangePassword(ctx context.Context, params ParamsChangePassword) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opChangePassword, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	body := itypes.NewPassword{
		Password: params.Password,
	}

	ep := endpoints.ChangePassword()
	_, err := c.c.Do(
		ctx,
		ep,
		withOrgID(ep.PathParams[0]),
		cav.WithPathParam(ep.PathParams[1], idOrName),
		cav.SetCustomRestyOption(setXMLHeaders),
		cav.SetBody(body),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opChangePassword, err)
	}

	return nil
}
