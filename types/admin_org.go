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

// ModelAdminOrg represents the public, read-only model for an organization in
// the XML-backed AdminOrg view.
type ModelAdminOrg struct {
	// ID of the organization
	// Example: urn:vcloud:org:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string `documentation:"ID of the organization"`

	// Name of the organization
	Name string `documentation:"Name of the organization"`

	// Href of the organization
	Href string `documentation:"Href of the organization"`

	// Display name of the organization
	DisplayName string `documentation:"Display name of the organization"`

	// Description of the organization
	Description string `documentation:"Description of the organization"`

	// Indicates if the organization is enabled
	IsEnabled bool `documentation:"Indicates if the organization is enabled"`

	// Indicates if the organization is full protected
	IsFullProtected bool `documentation:"Indicates if the organization is full protected"`
}

// ParamsGetAdminOrg defines parameters for reading an organization from the
// XML-backed AdminOrg view.
type ParamsGetAdminOrg struct {
	// ID of the organization
	// Example: urn:vcloud:org:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string

	// Name of the organization
	Name string
}

// Validate validates parameters for reading an AdminOrg view.
func (p *ParamsGetAdminOrg) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	if p.ID != "" {
		return validators.New().Var(p.ID, "urn=org")
	}
	return nil
}
