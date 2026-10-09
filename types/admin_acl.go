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

// ModelAdminCatalogACL represents the public model for a catalog ACL (admin scope).
type ModelAdminCatalogACL struct {
	// Indicates if the catalog is shared to everyone
	IsSharedToEveryone bool `documentation:"Indicates if the catalog is shared to everyone"`

	// Access level granted to everyone
	EveryoneAccessLevel *string `documentation:"Access level granted to everyone"`

	// List of subjects the catalog is shared with
	SharedWith []ModelAdminCatalogACLSubject `documentation:"List of subjects the catalog is shared with"`
}

// ModelAdminCatalogACLSubject represents a single subject shared with in a catalog ACL.
type ModelAdminCatalogACLSubject struct {
	// ID of the user the catalog is shared with
	UserID string `documentation:"ID of the user the catalog is shared with"`

	// ID of the group the catalog is shared with
	GroupID string `documentation:"ID of the group the catalog is shared with"`

	// Access level granted to the subject
	AccessLevel string `documentation:"Access level granted to the subject"`

	// Name of the subject
	SubjectName string `documentation:"Name of the subject"`
}

// ParamsGetAdminCatalogACL defines the parameters for getting a catalog ACL (admin scope).
type ParamsGetAdminCatalogACL struct {
	// ID of the catalog
	// Example: urn:vcloud:catalog:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string

	// Name of the catalog
	Name string
}

// Validate validates the parameters for getting a catalog ACL (admin scope).
func (p *ParamsGetAdminCatalogACL) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	if p.ID != "" {
		return validators.New().Var(p.ID, "urn=catalog")
	}
	return nil
}

// ParamsAdminCatalogACLSubject represents a single subject the catalog ACL is shared with.
type ParamsAdminCatalogACLSubject struct {
	// ID of the user the catalog is shared with
	// Example: urn:vcloud:orguser:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	UserID string

	// Access level granted to the subject
	AccessLevel string
}

// ParamsSetAdminCatalogACL defines the parameters for setting a catalog ACL (admin scope).
type ParamsSetAdminCatalogACL struct {
	// ID of the catalog
	// Example: urn:vcloud:catalog:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string

	// Indicates if the catalog is shared to everyone
	IsSharedToEveryone bool

	// Access level granted to everyone
	EveryoneAccessLevel *string

	// List of subjects the catalog is shared with
	SharedWith []ParamsAdminCatalogACLSubject
}

// Validate validates the parameters for setting a catalog ACL (admin scope).
func (p *ParamsSetAdminCatalogACL) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if err := validators.New().Var(p.ID, "urn=catalog"); err != nil {
		return fmt.Errorf("id: %w", err)
	}
	for i, s := range p.SharedWith {
		if s.UserID == "" {
			return fmt.Errorf("sharedWith[%d]: userId is required", i)
		}
		if s.AccessLevel == "" {
			return fmt.Errorf("sharedWith[%d]: accessLevel is required", i)
		}
	}
	return nil
}
