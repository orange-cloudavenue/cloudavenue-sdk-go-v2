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

// ModelCatalogAccessControlGrant represents a single access control grant for a catalog.
type ModelCatalogAccessControlGrant struct {
	// ID of the access control grant.
	AccessControlID string `documentation:"ID of the access control grant"`

	// Name of the subject (org/user/group/serviceaccount).
	SubjectName string `documentation:"Name of the subject"`

	// Type of the subject (org/user/group/serviceaccount).
	SubjectType string `documentation:"Type of the subject"`

	// URN of the role granted to the subject.
	RoleURN string `documentation:"URN of the role granted to the subject"`

	// Name of the role granted to the subject.
	RoleName string `documentation:"Name of the role granted to the subject"`

	// ID of the organization the grant belongs to.
	OrganizationID string `documentation:"ID of the organization the grant belongs to"`

	// Name of the organization the grant belongs to.
	OrganizationName string `documentation:"Name of the organization the grant belongs to"`
}

// ModelListCatalogAccessControlGrant is the list of catalog access control grants.
type ModelListCatalogAccessControlGrant struct {
	Grants []ModelCatalogAccessControlGrant `documentation:"List of catalog access control grants"`
}

// ParamsGetCatalogAccessControl defines the parameters for getting catalog access control.
type ParamsGetCatalogAccessControl struct {
	// URN of the catalog.
	// Example: urn:vcloud:catalog:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	CatalogURN string
}

// Validate validates the parameters for getting catalog access control.
func (p *ParamsGetCatalogAccessControl) Validate() error {
	if p.CatalogURN == "" {
		return fmt.Errorf("catalogUrn is required")
	}
	return nil
}

// ParamsSetCatalogAccessControl defines the parameters for setting catalog access control.
type ParamsSetCatalogAccessControl struct {
	// URN of the catalog.
	// Example: urn:vcloud:catalog:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	CatalogURN string

	// Grants is the list of access control grants to apply.
	Grants []ModelCatalogAccessControlGrant
}

// Validate validates the parameters for setting catalog access control.
func (p *ParamsSetCatalogAccessControl) Validate() error {
	if p.CatalogURN == "" {
		return fmt.Errorf("catalogUrn is required")
	}
	for i, grant := range p.Grants {
		if grant.SubjectName == "" {
			return fmt.Errorf("grants[%d]: subjectName is required", i)
		}
		if grant.SubjectType == "" {
			return fmt.Errorf("grants[%d]: subjectType is required", i)
		}
		if grant.RoleURN == "" {
			return fmt.Errorf("grants[%d]: roleUrn is required", i)
		}
	}
	return nil
}
