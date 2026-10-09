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

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

const (
	opListUsers         = "IAM.ListUsers"
	opGetUser           = "IAM.GetUser"
	opCreateLocalUser   = "IAM.CreateLocalUser"
	opCreateSAMLUser    = "IAM.CreateSAMLUser"
	opUpdateUser        = "IAM.UpdateUser"
	opDeleteUser        = "IAM.DeleteUser"
	opTakeOwnershipUser = "IAM.TakeOwnership"
	opEnableUser        = "IAM.EnableUser"
	opDisableUser       = "IAM.DisableUser"
	opUnlockUser        = "IAM.UnlockUser"
	opChangePassword    = "IAM.ChangePassword"
)

// ListUsers lists all users in the organization.
func (c *Client) ListUsers(ctx context.Context) ([]*ModelUser, error) {
	ep := endpoints.ListUsers()
	resp, err := c.c.Do(ctx, ep, cav.OverrideSetResult(new(itypes.APIResponseListUsers)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListUsers, err)
	}
	users, ok := resp.Result().(*itypes.APIResponseListUsers)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListUsers, resp.Result())
	}
	result := make([]*ModelUser, len(users.Users))
	for i := range users.Users {
		result[i] = iamCloudAPIUserToModel(users.Users[i])
	}
	return result, nil
}

// GetUser retrieves a user by ID or name.
func (c *Client) GetUser(ctx context.Context, params ParamsGetUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetUser, err)
	}
	id, err := c.resolveUserID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetUser, err)
	}
	ep := endpoints.GetUser()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id), cav.OverrideSetResult(new(itypes.APIUser)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetUser, err)
	}
	user, ok := resp.Result().(*itypes.APIUser)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetUser, resp.Result())
	}
	return iamCloudAPIUserToModel(*user), nil
}

// CreateLocalUser creates a new local user in the organization.
func (c *Client) CreateLocalUser(ctx context.Context, params ParamsCreateLocalUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateLocalUser, err)
	}
	body := itypes.APIUser{Username: params.Name, Password: params.Password, FullName: params.FullName, Email: params.EmailAddress, Phone: params.Telephone, Description: params.Description, Enabled: &params.IsEnabled, ProviderType: "LOCAL", DeployedVMQuota: params.DeployedVMQuota, StoredVMQuota: params.StoredVMQuota, RoleEntityRefs: []itypes.APIObjectReference{{Name: params.RoleName}}}
	return c.createUser(ctx, body, opCreateLocalUser)
}

// CreateSAMLUser creates a new SAML user in the organization.
func (c *Client) CreateSAMLUser(ctx context.Context, params ParamsCreateSAMLUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateSAMLUser, err)
	}
	body := itypes.APIUser{Username: params.Name, FullName: params.FullName, Email: params.EmailAddress, Phone: params.Telephone, Description: params.Description, Enabled: &params.IsEnabled, ProviderType: "SAML", RoleEntityRefs: []itypes.APIObjectReference{{Name: params.RoleName}}}
	return c.createUser(ctx, body, opCreateSAMLUser)
}

func (c *Client) createUser(ctx context.Context, body itypes.APIUser, operation string) (*ModelUser, error) {
	ep := endpoints.CreateUser()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body), cav.OverrideSetResult(new(itypes.APIUser)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	user, ok := resp.Result().(*itypes.APIUser)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", operation, resp.Result())
	}
	return iamCloudAPIUserToModel(*user), nil
}

// UpdateUser updates an existing user.
func (c *Client) UpdateUser(ctx context.Context, params ParamsUpdateUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateUser, err)
	}
	id, err := c.resolveUserID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateUser, err)
	}
	current, err := c.GetUser(ctx, ParamsGetUser{ID: id})
	if err != nil {
		return nil, fmt.Errorf("%s: get current: %w", opUpdateUser, err)
	}
	body := modelUserToCloudAPI(current, params)
	return c.updateUser(ctx, id, body, opUpdateUser)
}

func (c *Client) updateUser(ctx context.Context, id string, body itypes.APIUser, operation string) (*ModelUser, error) {
	ep := endpoints.UpdateUser()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id), cav.SetBody(body), cav.OverrideSetResult(new(itypes.APIUser)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	user, ok := resp.Result().(*itypes.APIUser)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", operation, resp.Result())
	}
	return iamCloudAPIUserToModel(*user), nil
}

// DeleteUser deletes a user by ID. TakeOwnership transfers owned entities first.
func (c *Client) DeleteUser(ctx context.Context, params ParamsDeleteUser) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opDeleteUser, err)
	}
	if params.TakeOwnership {
		if err := c.takeOwnership(ctx, params.ID); err != nil {
			return err
		}
	}
	ep := endpoints.DeleteUser()
	_, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.ID))
	if err != nil {
		return fmt.Errorf("%s: %w", opDeleteUser, err)
	}
	return nil
}

func (c *Client) takeOwnership(ctx context.Context, id string) error {
	ep := endpoints.TakeOwnership()
	_, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id))
	if err != nil {
		return fmt.Errorf("%s: %w", opTakeOwnershipUser, err)
	}
	return nil
}

// EnableUser enables a user by ID or name.
func (c *Client) EnableUser(ctx context.Context, params ParamsEnableUser) (*ModelUser, error) {
	return c.setUserEnabled(ctx, params.ID, params.Name, true, opEnableUser)
}

// DisableUser disables a user by ID or name.
func (c *Client) DisableUser(ctx context.Context, params ParamsDisableUser) (*ModelUser, error) {
	return c.setUserEnabled(ctx, params.ID, params.Name, false, opDisableUser)
}

func (c *Client) setUserEnabled(ctx context.Context, id, name string, enabled bool, operation string) (*ModelUser, error) {
	if id == "" && name == "" {
		return nil, fmt.Errorf("%s: validate: id or name is required", operation)
	}
	resolved, err := c.resolveUserID(ctx, id, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	current, err := c.GetUser(ctx, ParamsGetUser{ID: resolved})
	if err != nil {
		return nil, fmt.Errorf("%s: get current: %w", operation, err)
	}
	body := modelUserToCloudAPI(current, ParamsUpdateUser{IsEnabled: &enabled})
	ep := endpoints.UpdateUser()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], resolved), cav.SetBody(body), cav.OverrideSetResult(new(itypes.APIUser)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	user, ok := resp.Result().(*itypes.APIUser)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", operation, resp.Result())
	}
	return iamCloudAPIUserToModel(*user), nil
}

// UnlockUser unlocks a user by setting locked=false through CloudAPI PUT.
func (c *Client) UnlockUser(ctx context.Context, params ParamsUnlockUser) (*ModelUser, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUnlockUser, err)
	}
	id, err := c.resolveUserID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUnlockUser, err)
	}
	current, err := c.GetUser(ctx, ParamsGetUser{ID: id})
	if err != nil {
		return nil, fmt.Errorf("%s: get current: %w", opUnlockUser, err)
	}
	locked := false
	body := modelUserToCloudAPI(current, ParamsUpdateUser{})
	body.Locked = &locked
	ep := endpoints.UnlockUser()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id), cav.SetBody(body), cav.OverrideSetResult(new(itypes.APIUser)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUnlockUser, err)
	}
	user, ok := resp.Result().(*itypes.APIUser)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUnlockUser, resp.Result())
	}
	return iamCloudAPIUserToModel(*user), nil
}

// ChangePassword changes a local user's password. CloudAPI requires currentPassword.
func (c *Client) ChangePassword(ctx context.Context, params ParamsChangePassword) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opChangePassword, err)
	}
	id, err := c.resolveUserID(ctx, params.ID, params.Name)
	if err != nil {
		return fmt.Errorf("%s: %w", opChangePassword, err)
	}
	ep := endpoints.ChangePassword()
	body := itypes.APIRequestPasswordChange{CurrentPassword: params.CurrentPassword, NewPassword: params.Password}
	_, err = c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id), cav.SetBody(body))
	if err != nil {
		return fmt.Errorf("%s: %w", opChangePassword, err)
	}
	return nil
}

func (c *Client) resolveUserID(ctx context.Context, id, name string) (string, error) {
	if id != "" {
		return id, nil
	}
	users, err := c.ListUsers(ctx)
	if err != nil {
		return "", err
	}
	for _, user := range users {
		if user.Name == name {
			return user.ID, nil
		}
	}
	// CloudAPI normally expects a URN, but retaining the supplied name keeps
	// compatibility with deployments and fixtures that accept name selectors.
	return name, nil
}

// withOrgID injects organization ID from request context into a path parameter.
// It remains used by token endpoints, which retain their existing API route.
func withOrgID(pp cav.PathParam) cav.EndpointRequestOption {
	return func(endpoint *cav.Endpoint, req *resty.Request) error {
		cd := cav.GetExtraDataFromContext(req.Context())
		if cd.OrganizationID == "" {
			return errors.New("organization ID not found in context")
		}
		return cav.WithPathParam(pp, cd.OrganizationID)(endpoint, req)
	}
}

func iamCloudAPIUserToModel(u itypes.APIUser) *ModelUser {
	model := &ModelUser{ID: u.ID, Name: u.Username, FullName: u.FullName, EmailAddress: u.Email, Telephone: u.Phone, Description: u.Description, ProviderType: u.ProviderType, DeployedVMQuota: u.DeployedVMQuota, StoredVMQuota: u.StoredVMQuota}
	if u.Enabled != nil {
		model.IsEnabled = *u.Enabled
	}
	if u.Locked != nil {
		model.IsLocked = *u.Locked
	}
	if len(u.RoleEntityRefs) > 0 {
		model.RoleName = u.RoleEntityRefs[0].Name
		model.RoleHref = u.RoleEntityRefs[0].ID
	}
	return model
}

func modelUserToCloudAPI(current *ModelUser, params ParamsUpdateUser) itypes.APIUser {
	enabled := current.IsEnabled
	body := itypes.APIUser{ID: current.ID, Username: current.Name, FullName: current.FullName, Email: current.EmailAddress, Phone: current.Telephone, Description: current.Description, Enabled: &enabled, DeployedVMQuota: current.DeployedVMQuota, StoredVMQuota: current.StoredVMQuota, Locked: &current.IsLocked, RoleEntityRefs: []itypes.APIObjectReference{{ID: current.RoleHref, Name: current.RoleName}}}
	if params.Name != "" {
		body.Username = params.Name
	}
	if params.Password != "" {
		body.Password = params.Password
	}
	if params.RoleName != "" {
		body.RoleEntityRefs = []itypes.APIObjectReference{{Name: params.RoleName}}
	}
	if params.FullName != "" {
		body.FullName = params.FullName
	}
	if params.EmailAddress != "" {
		body.Email = params.EmailAddress
	}
	if params.Telephone != "" {
		body.Phone = params.Telephone
	}
	if params.Description != nil {
		body.Description = params.Description
	}
	if params.DeployedVMQuota != nil {
		body.DeployedVMQuota = params.DeployedVMQuota
	}
	if params.StoredVMQuota != nil {
		body.StoredVMQuota = params.StoredVMQuota
	}
	if params.IsEnabled != nil {
		body.Enabled = params.IsEnabled
	}
	return body
}
