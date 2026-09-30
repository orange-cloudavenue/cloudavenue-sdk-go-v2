/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

import "fmt"

// ModelToken represents a CloudAvenue IAM token.
type ModelToken struct {
	// ID of the token in URN format
	ID string `documentation:"URN of the token in the format urn:vcloud:token:<UUID>"`

	// Name of the token
	Name string `documentation:"Name of the token"`

	// Description of the token
	Description string `documentation:"Description of the token"`

	// Indicates if the token is enabled
	Enabled bool `documentation:"Indicates if the token is enabled"`

	// Name of the token's role
	RoleName string `documentation:"Name of the token's role"`

	// Href of the token's role
	RoleHref string `documentation:"Href of the token's role"`
}

// ParamsGetToken defines parameters for getting a token.
type ParamsGetToken struct {
	ID string
}

// Validate checks ParamsGetToken structural constraints.
func (p ParamsGetToken) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsCreateToken defines parameters for creating a token.
type ParamsCreateToken struct {
	Name        string
	RoleID      string
	RoleName    string
	Description string
	IsEnabled   bool
}

// Validate checks ParamsCreateToken structural constraints.
func (p ParamsCreateToken) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.RoleName == "" {
		return fmt.Errorf("role_name is required")
	}
	return nil
}

// ParamsUpdateToken defines parameters for updating a token.
type ParamsUpdateToken struct {
	ID          string
	Name        string
	RoleID      string
	RoleName    string
	Description string
	IsEnabled   bool
}

// Validate checks ParamsUpdateToken structural constraints.
func (p ParamsUpdateToken) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsDeleteToken defines parameters for deleting a token.
type ParamsDeleteToken struct {
	ID string
}

// Validate checks ParamsDeleteToken structural constraints.
func (p ParamsDeleteToken) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsEntityReference represents an entity reference (ID + Name) used in global role params.
type ParamsEntityReference struct {
	ID   string
	Name string
}

// ModelRight represents a right granted to a global role.
type ModelRight struct {
	// ID of the right
	ID string `documentation:"ID of the right"`

	// Name of the right
	Name string `documentation:"Name of the right"`

	// Description of the right
	Description string `documentation:"Description of the right"`

	// BundleKey of the right
	BundleKey string `documentation:"Bundle key of the right"`

	// Category of the right
	Category string `documentation:"Category of the right"`

	// RightType of the right
	RightType string `documentation:"Type of the right"`

	// IsPublishable indicates if the right is publishable
	IsPublishable bool `documentation:"Indicates if the right is publishable"`

	// ServiceNamespace of the right
	ServiceNamespace string `documentation:"Service namespace of the right"`
}

// ModelTenant represents a tenant associated with a global role.
type ModelTenant struct {
	// ID of the tenant
	ID string `documentation:"ID of the tenant"`

	// Name of the tenant
	Name string `documentation:"Name of the tenant"`

	// OrgID of the tenant
	OrgID string `documentation:"Organization ID of the tenant"`
}

// ModelGlobalRole represents a CloudAvenue IAM global role.
type ModelGlobalRole struct {
	// ID of the global role
	ID string `documentation:"ID of the global role"`

	// Name of the global role
	Name string `documentation:"Name of the global role"`

	// Description of the global role
	Description string `documentation:"Description of the global role"`

	// Rights of the global role
	Rights []*ModelRight `documentation:"Rights of the global role"`

	// NumberOfTenants is the number of tenants associated with the global role
	NumberOfTenants int `documentation:"Number of tenants associated with the global role"`

	// CanPublish indicates if the global role can be published
	CanPublish bool `documentation:"Indicates if the global role can be published"`

	// IsInternalRole indicates if the global role is an internal role
	IsInternalRole bool `documentation:"Indicates if the global role is an internal role"`
}

// ParamsGetGlobalRole defines parameters for getting a global role.
type ParamsGetGlobalRole struct {
	ID string
}

// Validate checks ParamsGetGlobalRole structural constraints.
func (p ParamsGetGlobalRole) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsCreateGlobalRole defines parameters for creating a global role.
type ParamsCreateGlobalRole struct {
	Name        string
	Description string
	Rights      []ParamsEntityReference
}

// Validate checks ParamsCreateGlobalRole structural constraints.
func (p ParamsCreateGlobalRole) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// ParamsUpdateGlobalRole defines parameters for updating a global role.
type ParamsUpdateGlobalRole struct {
	ID          string
	Name        string
	Description string
	Rights      []ParamsEntityReference
}

// Validate checks ParamsUpdateGlobalRole structural constraints.
func (p ParamsUpdateGlobalRole) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsDeleteGlobalRole defines parameters for deleting a global role.
type ParamsDeleteGlobalRole struct {
	ID string
}

// Validate checks ParamsDeleteGlobalRole structural constraints.
func (p ParamsDeleteGlobalRole) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsListGlobalRoleRights defines parameters for listing global role rights.
type ParamsListGlobalRoleRights struct {
	ID string
}

// Validate checks ParamsListGlobalRoleRights structural constraints.
func (p ParamsListGlobalRoleRights) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsAddGlobalRoleRights defines parameters for adding rights to a global role.
type ParamsAddGlobalRoleRights struct {
	ID     string
	Rights []ParamsEntityReference
}

// Validate checks ParamsAddGlobalRoleRights structural constraints.
func (p ParamsAddGlobalRoleRights) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if len(p.Rights) == 0 {
		return fmt.Errorf("rights is required")
	}
	return nil
}

// ParamsReplaceGlobalRoleRights defines parameters for replacing global role rights.
type ParamsReplaceGlobalRoleRights struct {
	ID     string
	Rights []ParamsEntityReference
}

// Validate checks ParamsReplaceGlobalRoleRights structural constraints.
func (p ParamsReplaceGlobalRoleRights) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if len(p.Rights) == 0 {
		return fmt.Errorf("rights is required")
	}
	return nil
}

// ParamsListGlobalRoleTenants defines parameters for listing global role tenants.
type ParamsListGlobalRoleTenants struct {
	ID string
}

// Validate checks ParamsListGlobalRoleTenants structural constraints.
func (p ParamsListGlobalRoleTenants) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsSetGlobalRoleTenants defines parameters for setting global role tenants.
type ParamsSetGlobalRoleTenants struct {
	ID      string
	Tenants []ParamsEntityReference
}

// Validate checks ParamsSetGlobalRoleTenants structural constraints.
func (p ParamsSetGlobalRoleTenants) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if len(p.Tenants) == 0 {
		return fmt.Errorf("tenants is required")
	}
	return nil
}

// ParamsPublishGlobalRoleTenants defines parameters for publishing global role tenants.
type ParamsPublishGlobalRoleTenants struct {
	ID      string
	Tenants []ParamsEntityReference
}

// Validate checks ParamsPublishGlobalRoleTenants structural constraints.
func (p ParamsPublishGlobalRoleTenants) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if len(p.Tenants) == 0 {
		return fmt.Errorf("tenants is required")
	}
	return nil
}

// ParamsUnpublishGlobalRoleTenants defines parameters for unpublishing global role tenants.
type ParamsUnpublishGlobalRoleTenants struct {
	ID      string
	Tenants []ParamsEntityReference
}

// Validate checks ParamsUnpublishGlobalRoleTenants structural constraints.
func (p ParamsUnpublishGlobalRoleTenants) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if len(p.Tenants) == 0 {
		return fmt.Errorf("tenants is required")
	}
	return nil
}

// ParamsPublishAllGlobalRoleTenants defines parameters for publishing all global role tenants.
type ParamsPublishAllGlobalRoleTenants struct {
	ID string
}

// Validate checks ParamsPublishAllGlobalRoleTenants structural constraints.
func (p ParamsPublishAllGlobalRoleTenants) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsUnpublishAllGlobalRoleTenants defines parameters for unpublishing all global role tenants.
type ParamsUnpublishAllGlobalRoleTenants struct {
	ID string
}

// Validate checks ParamsUnpublishAllGlobalRoleTenants structural constraints.
func (p ParamsUnpublishAllGlobalRoleTenants) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}
