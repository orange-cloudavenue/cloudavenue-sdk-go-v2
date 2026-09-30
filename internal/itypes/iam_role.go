/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

// * Request / Response API for Global Roles

type (
	// APIRequestGlobalRole represents the request body for global role operations.
	APIRequestGlobalRole struct {
		ID          string                      `json:"id,omitempty"`
		Name        string                      `json:"name,omitempty"`
		Description string                      `json:"description,omitempty"`
		Rights      []APIRequestGlobalRoleRight `json:"rights,omitempty"`
	}

	// APIRequestGlobalRoleRight represents a right reference in a global role request.
	APIRequestGlobalRoleRight struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	// APIResponseGlobalRole represents a global role response.
	APIResponseGlobalRole struct {
		ID              string                       `json:"id,omitempty"`
		Name            string                       `json:"name,omitempty"`
		Description     string                       `json:"description,omitempty"`
		Rights          []APIResponseGlobalRoleRight `json:"rights,omitempty"`
		NumberOfTenants int                          `json:"numberOfTenants,omitempty"`
		CanPublish      bool                         `json:"canPublish,omitempty"`
		IsInternalRole  bool                         `json:"isInternalRole,omitempty"`
	}

	// APIResponseGlobalRoleRight represents a right in a global role response.
	APIResponseGlobalRoleRight struct {
		ID               string `json:"id,omitempty"`
		Name             string `json:"name,omitempty"`
		Description      string `json:"description,omitempty"`
		BundleKey        string `json:"bundleKey,omitempty"`
		Category         string `json:"category,omitempty"`
		RightType        string `json:"rightType,omitempty"`
		IsPublishable    bool   `json:"isPublishable,omitempty"`
		ServiceNamespace string `json:"serviceNamespace,omitempty"`
	}

	// APIResponseListGlobalRoles represents the wrapper for a list of global roles.
	APIResponseListGlobalRoles struct {
		GlobalRoles []APIResponseGlobalRole `json:"globalRoles,omitempty"`
	}

	// APIResponseGlobalRoleRights represents the wrapper for a list of global role rights.
	APIResponseGlobalRoleRights struct {
		Rights []APIResponseGlobalRoleRight `json:"rights,omitempty"`
	}

	// APIRequestGlobalRoleRights represents the request body for adding/replacing global role rights.
	APIRequestGlobalRoleRights struct {
		Rights []APIRequestGlobalRoleRight `json:"rights,omitempty"`
	}

	// APIResponseGlobalRoleTenants represents the wrapper for a list of global role tenants.
	APIResponseGlobalRoleTenants struct {
		Tenants []APIResponseGlobalRoleTenant `json:"tenants,omitempty"`
	}

	// APIResponseGlobalRoleTenant represents a tenant in a global role response.
	APIResponseGlobalRoleTenant struct {
		ID    string `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		OrgID string `json:"orgId,omitempty"`
	}

	// APIRequestGlobalRoleTenants represents the request body for setting global role tenants.
	APIRequestGlobalRoleTenants struct {
		Tenants []APIRequestGlobalRoleTenant `json:"tenants,omitempty"`
	}

	// APIRequestGlobalRoleTenant represents a tenant reference in a global role request.
	APIRequestGlobalRoleTenant struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}
)

// ToModel converts the APIResponseGlobalRole to ModelGlobalRole.
func (api *APIResponseGlobalRole) ToModel() *types.ModelGlobalRole {
	if api == nil {
		return nil
	}

	rights := make([]*types.ModelRight, len(api.Rights))
	for i, r := range api.Rights {
		rights[i] = r.ToModel()
	}

	return &types.ModelGlobalRole{
		ID:              api.ID,
		Name:            api.Name,
		Description:     api.Description,
		Rights:          rights,
		NumberOfTenants: api.NumberOfTenants,
		CanPublish:      api.CanPublish,
		IsInternalRole:  api.IsInternalRole,
	}
}

// ToModel converts each global role in the list to ModelGlobalRole.
func (api *APIResponseListGlobalRoles) ToModel() []*types.ModelGlobalRole {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelGlobalRole, len(api.GlobalRoles))
	for i, r := range api.GlobalRoles {
		result[i] = r.ToModel()
	}

	return result
}

// ToModel converts the APIResponseGlobalRoleRight to ModelRight.
func (api *APIResponseGlobalRoleRight) ToModel() *types.ModelRight {
	if api == nil {
		return nil
	}

	return &types.ModelRight{
		ID:               api.ID,
		Name:             api.Name,
		Description:      api.Description,
		BundleKey:        api.BundleKey,
		Category:         api.Category,
		RightType:        api.RightType,
		IsPublishable:    api.IsPublishable,
		ServiceNamespace: api.ServiceNamespace,
	}
}

// ToModel converts each right in the list to ModelRight.
func (api *APIResponseGlobalRoleRights) ToModel() []*types.ModelRight {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelRight, len(api.Rights))
	for i, r := range api.Rights {
		result[i] = r.ToModel()
	}

	return result
}

// ToModel converts the APIResponseGlobalRoleTenant to ModelTenant.
func (api *APIResponseGlobalRoleTenant) ToModel() *types.ModelTenant {
	if api == nil {
		return nil
	}

	return &types.ModelTenant{
		ID:    api.ID,
		Name:  api.Name,
		OrgID: api.OrgID,
	}
}

// ToModel converts each tenant in the list to ModelTenant.
func (api *APIResponseGlobalRoleTenants) ToModel() []*types.ModelTenant {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelTenant, len(api.Tenants))
	for i, t := range api.Tenants {
		result[i] = t.ToModel()
	}

	return result
}

// * Request / Response API for Tokens

type (
	// APIRequestToken represents the request body for token operations.
	APIRequestToken struct {
		ID          string             `json:"id,omitempty"`
		Name        string             `json:"name,omitempty"`
		Description string             `json:"description,omitempty"`
		Role        APIObjectReference `json:"role,omitempty"`
		Enabled     bool               `json:"enabled,omitempty"`
	}

	// APIResponseToken represents a token response.
	APIResponseToken struct {
		ID          string             `json:"id,omitempty"`
		Name        string             `json:"name,omitempty"`
		Description string             `json:"description,omitempty"`
		Enabled     bool               `json:"enabled,omitempty"`
		Role        APIObjectReference `json:"role,omitempty"`
	}

	// APIResponseListTokens represents the wrapper for a list of tokens.
	APIResponseListTokens struct {
		Tokens []APIResponseToken `json:"tokens,omitempty"`
	}
)

// ToModel converts the APIResponseToken to ModelToken.
func (api *APIResponseToken) ToModel() *types.ModelToken {
	if api == nil {
		return nil
	}

	return &types.ModelToken{
		ID:          api.ID,
		Name:        api.Name,
		Description: api.Description,
		Enabled:     api.Enabled,
		RoleName:    api.Role.Name,
		RoleHref:    api.Role.ID,
	}
}

// ToModel converts each token in the list to ModelToken.
func (api *APIResponseListTokens) ToModel() []*types.ModelToken {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelToken, len(api.Tokens))
	for i, t := range api.Tokens {
		result[i] = t.ToModel()
	}

	return result
}
