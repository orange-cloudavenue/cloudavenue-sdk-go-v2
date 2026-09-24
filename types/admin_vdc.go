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

// ModelAdminVDC represents the public model for a VDC (admin scope).
type ModelAdminVDC struct {
	// ID of the VDC
	// Example: urn:vcloud:vdc:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string `documentation:"ID of the VDC"`

	// Name of the VDC
	Name string `documentation:"Name of the VDC"`

	// Href of the VDC
	Href string `documentation:"Href of the VDC"`

	// Display name of the VDC
	DisplayName string `documentation:"Display name of the VDC"`

	// Description of the VDC
	Description string `documentation:"Description of the VDC"`

	// Indicates if the VDC is enabled
	IsEnabled bool `documentation:"Indicates if the VDC is enabled"`

	// Indicates if the VDC is full protected
	IsFullProtected bool `documentation:"Indicates if the VDC is full protected"`

	// Number of VMs in the VDC
	VmCount int `documentation:"Number of VMs in the VDC"`

	// Number of running VMs in the VDC
	VmRunningCount int `documentation:"Number of running VMs in the VDC"`

	// Number of deployed vApps in the VDC
	VappCount int `documentation:"Number of deployed vApps in the VDC"`

	// Number of vApp templates in the VDC
	VappTemplateCount int `documentation:"Number of vApp templates in the VDC"`
}

// ParamsGetAdminVDC defines the parameters for getting a VDC (admin scope).
type ParamsGetAdminVDC struct {
	// ID of the VDC
	// Example: urn:vcloud:vdc:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string

	// Name of the VDC
	Name string
}

// Validate validates the parameters for getting a VDC (admin scope).
func (p *ParamsGetAdminVDC) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	if p.ID != "" {
		return validators.New().Var(p.ID, "urn=vdc")
	}
	return nil
}
