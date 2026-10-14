/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

import (
	"fmt"

	"github.com/orange-cloudavenue/common-go/validators"
)

// ModelAdminCatalog represents the public model for a catalog (admin scope).
type ModelAdminCatalog struct {
	// ID of the catalog
	// Example: urn:vcloud:catalog:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string `documentation:"ID of the catalog"`

	// Name of the catalog
	Name string `documentation:"Name of the catalog"`

	// Href of the catalog
	Href string `documentation:"Href of the catalog"`

	// Description of the catalog
	Description string `documentation:"Description of the catalog"`

	// Indicates if the catalog is enabled
	IsEnabled bool `documentation:"Indicates if the catalog is enabled"`

	// Indicates if the catalog is published
	IsPublished bool `documentation:"Indicates if the catalog is published"`

	// Indicates if the catalog is trusted
	IsTrusted bool `documentation:"Indicates if the catalog is trusted"`
}

// ParamsGetAdminCatalog defines the parameters for getting a catalog (admin scope).
type ParamsGetAdminCatalog struct {
	// ID of the catalog
	// Example: urn:vcloud:catalog:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string

	// Name of the catalog
	Name string
}

// Validate validates the parameters for getting a catalog (admin scope).
func (p *ParamsGetAdminCatalog) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	if p.ID != "" {
		return validators.New().Var(p.ID, "urn=catalog")
	}
	return nil
}

// ParamsCreateAdminCatalog defines the parameters for creating a catalog (admin scope).
type ParamsCreateAdminCatalog struct {
	// ID of the organization
	// Example: urn:vcloud:org:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	OrgID string

	// Name of the catalog
	Name string

	// Description of the catalog
	Description string

	// Storage profile IDs of the catalog
	StorageProfileIDs []string
}

// Validate validates the parameters for creating a catalog (admin scope).
func (p *ParamsCreateAdminCatalog) Validate() error {
	if p.OrgID == "" {
		return fmt.Errorf("org id is required")
	}
	if err := validators.New().Var(p.OrgID, "urn=org"); err != nil {
		return fmt.Errorf("org id: %w", err)
	}
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// ParamsUpdateAdminCatalog defines the parameters for updating a catalog (admin scope).
type ParamsUpdateAdminCatalog struct {
	// ID of the catalog
	// Example: urn:vcloud:catalog:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string

	// Name of the catalog
	Name string

	// Description of the catalog
	Description string

	// Storage profile IDs of the catalog
	StorageProfileIDs []string
}

// Validate validates the parameters for updating a catalog (admin scope).
func (p *ParamsUpdateAdminCatalog) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if err := validators.New().Var(p.ID, "urn=catalog"); err != nil {
		return fmt.Errorf("id: %w", err)
	}
	return nil
}
