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
		BundleKey       string                       `json:"bundleKey,omitempty"`
		ReadOnly        bool                         `json:"readOnly"`
		Rights          []APIResponseGlobalRoleRight `json:"rights,omitempty"`
		NumberOfTenants int                          `json:"numberOfTenants,omitempty"`
		CanPublish      bool                         `json:"canPublish"`
		IsInternalRole  bool                         `json:"isInternalRole"`
		PublishAll      bool                         `json:"publishAll"`
	}

	// APIResponseGlobalRoleRight represents a right in a global role response.
	APIResponseGlobalRoleRight struct {
		ID               string `json:"id,omitempty"`
		Name             string `json:"name,omitempty"`
		Description      string `json:"description,omitempty"`
		BundleKey        string `json:"bundleKey,omitempty"`
		Category         string `json:"category,omitempty"`
		RightType        string `json:"rightType,omitempty"`
		IsPublishable    bool   `json:"isPublishable"`
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
		BundleKey:       api.BundleKey,
		ReadOnly:        api.ReadOnly,
		Rights:          rights,
		NumberOfTenants: api.NumberOfTenants,
		CanPublish:      api.CanPublish,
		IsInternalRole:  api.IsInternalRole,
		PublishAll:      api.PublishAll,
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
		Type        string             `json:"type,omitempty"`
		ExtensionID string             `json:"extensionId,omitempty"`
		Role        APIObjectReference `json:"role,omitzero"`
		Enabled     bool               `json:"enabled"`
	}

	// APIResponseToken represents a token response.
	APIResponseToken struct {
		ID                string             `json:"id,omitempty"`
		Name              string             `json:"name,omitempty"`
		Description       string             `json:"description,omitempty"`
		Token             string             `json:"token,omitempty"`
		ExpirationTimeUTC string             `json:"expirationTimeUtc,omitempty"`
		ExtensionID       string             `json:"extensionId,omitempty"`
		Owner             APIObjectReference `json:"owner"`
		Org               APIObjectReference `json:"org"`
		Type              string             `json:"type,omitempty"`
		RequireRotation   *bool              `json:"requireRotation,omitempty"`
		Enabled           bool               `json:"enabled"`
		Role              APIObjectReference `json:"role"`
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
		ID:                api.ID,
		Name:              api.Name,
		Description:       api.Description,
		Enabled:           api.Enabled,
		RoleName:          api.Role.Name,
		RoleHref:          api.Role.ID,
		Type:              api.Type,
		ExtensionID:       api.ExtensionID,
		Token:             api.Token,
		ExpirationTimeUTC: api.ExpirationTimeUTC,
		OwnerID:           api.Owner.ID,
		OwnerName:         api.Owner.Name,
		OrgID:             api.Org.ID,
		OrgName:           api.Org.Name,
		RequireRotation:   api.RequireRotation,
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

// * Request / Response API for Org-scoped Roles

type (
	// APIRequestRole represents the request body for role operations.
	APIRequestRole struct {
		ID          string `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		BundleKey   string `json:"bundleKey,omitempty"`
	}

	// APIResponseRole represents a role response.
	APIResponseRole struct {
		ID          string `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		BundleKey   string `json:"bundleKey,omitempty"`
		ReadOnly    bool   `json:"readOnly"`
	}

	// APIResponseListRoles represents the pageable wrapper for a list of roles.
	APIResponseListRoles struct {
		Roles       []APIResponseRole `json:"values,omitempty"`
		ResultTotal int               `json:"resultTotal,omitempty"`
		PageCount   int               `json:"pageCount,omitempty"`
		Page        int               `json:"page,omitempty"`
		PageSize    int               `json:"pageSize,omitempty"`
	}

	// APIResponseRoleRights represents the pageable wrapper for a list of role rights.
	// CloudAPI serves role rights as plain entity references (id and name only).
	APIResponseRoleRights struct {
		Rights      []APIEntityReference `json:"values,omitempty"`
		ResultTotal int                  `json:"resultTotal,omitempty"`
		PageCount   int                  `json:"pageCount,omitempty"`
		Page        int                  `json:"page,omitempty"`
		PageSize    int                  `json:"pageSize,omitempty"`
	}

	// APIRequestRoleRights represents the request body for adding/replacing role rights.
	// Rights carries no omitempty so that removing all rights sends an explicit empty list.
	APIRequestRoleRights struct {
		Rights []APIEntityReference `json:"values"`
	}
)

// ToModel converts the APIResponseRole to ModelRole.
func (api *APIResponseRole) ToModel() *types.ModelRole {
	if api == nil {
		return nil
	}

	return &types.ModelRole{
		ID:          api.ID,
		Name:        api.Name,
		Description: api.Description,
		BundleKey:   api.BundleKey,
		ReadOnly:    api.ReadOnly,
	}
}

// ToModel converts each role in the list to ModelRole.
func (api *APIResponseListRoles) ToModel() []*types.ModelRole {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelRole, len(api.Roles))
	for i, r := range api.Roles {
		result[i] = r.ToModel()
	}

	return result
}

// ToModel converts each right reference in the list to ModelRight.
func (api *APIResponseRoleRights) ToModel() []*types.ModelRight {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelRight, len(api.Rights))
	for i, r := range api.Rights {
		result[i] = &types.ModelRight{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	return result
}

// * Request / Response API for Rights

type (
	// APIResponseRight represents a right response.
	APIResponseRight struct {
		ID               string               `json:"id,omitempty"`
		Name             string               `json:"name,omitempty"`
		Description      string               `json:"description,omitempty"`
		BundleKey        string               `json:"bundleKey,omitempty"`
		Category         string               `json:"category,omitempty"`
		RightType        string               `json:"rightType,omitempty"`
		IsPublishable    bool                 `json:"isPublishable"`
		ServiceNamespace string               `json:"serviceNamespace,omitempty"`
		ImpliedRights    []APIEntityReference `json:"impliedRights,omitempty"`
	}

	// APIResponseListRights represents the pageable wrapper for a list of rights.
	APIResponseListRights struct {
		Rights      []APIResponseRight `json:"values,omitempty"`
		ResultTotal int                `json:"resultTotal,omitempty"`
		PageCount   int                `json:"pageCount,omitempty"`
		Page        int                `json:"page,omitempty"`
		PageSize    int                `json:"pageSize,omitempty"`
	}
)

// ToModel converts the APIResponseRight to ModelRight.
func (api *APIResponseRight) ToModel() *types.ModelRight {
	if api == nil {
		return nil
	}

	impliedRights := make([]*types.ModelRight, len(api.ImpliedRights))
	for i, ir := range api.ImpliedRights {
		impliedRights[i] = &types.ModelRight{
			ID:   ir.ID,
			Name: ir.Name,
		}
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
		ImpliedRights:    impliedRights,
	}
}

// ToModel converts each right in the list to ModelRight.
func (api *APIResponseListRights) ToModel() []*types.ModelRight {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelRight, len(api.Rights))
	for i, r := range api.Rights {
		result[i] = r.ToModel()
	}

	return result
}
