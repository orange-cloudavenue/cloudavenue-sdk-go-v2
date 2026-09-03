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

type (
	// ModelListVApp contains a list of VApps.
	ModelListVApp struct {
		VApps []ModelVApp `documentation:"List of VApps"`
	}

	// ModelVApp represents a VApp.
	ModelVApp struct {
		ID          string `documentation:"ID of the VApp"`
		Name        string `documentation:"Name of the VApp"`
		Description string `documentation:"Description of the VApp"`
		Status      string `documentation:"Status of the VApp"`
		Deployed    bool   `documentation:"Whether the VApp is deployed"`
		HREF        string `documentation:"HREF of the VApp"`
		VDC         string `documentation:"VDC ID where the VApp resides"`
		Org         string `documentation:"Organization ID where the VApp resides"`
		Owner       string `documentation:"Owner of the VApp"`
		Parent      string `documentation:"Parent VApp ID if nested"`
		NumberOfVMs int    `documentation:"Number of VMs in the VApp"`

		DeploymentLeaseInSeconds *int `documentation:"Deployment lease in seconds"`
		StorageLeaseInSeconds    *int `documentation:"Storage lease in seconds"`

		Properties ModelVAppProperties `documentation:"Properties of the VApp"`
	}

	// ModelVAppProperties contains VApp properties.
	ModelVAppProperties struct {
		BillingModel string `documentation:"Billing model of the VApp"`
	}
)

type (
	// ParamsListVApp is the parameters for the ListVApp command.
	ParamsListVApp struct {
		// VDCID is the VDC ID to filter VApps by.
		VDCID string `documentation:"VDC ID to filter VApps by"`
	}

	// ParamsGetVApp is the parameters for the GetVApp command.
	ParamsGetVApp struct {
		// ID is the unique identifier of the VApp to get.
		ID string

		// Name is the name of the VApp to get.
		Name string
	}

	// ParamsCreateVApp is the parameters for the CreateVApp command.
	ParamsCreateVApp struct {
		// Name is the name of the VApp to create.
		Name string

		// Description is the description of the VApp to create.
		Description string

		// VDCID is the VDC ID where the VApp will be created.
		VDCID string
	}

	// ParamsUpdateVApp is the parameters for the UpdateVApp command.
	ParamsUpdateVApp struct {
		// ID is the unique identifier of the VApp to update.
		ID string

		// Name is the name of the VApp to update.
		Name string

		// Description is the new description of the VApp.
		Description *string

		// DeploymentLeaseInSeconds is the deployment lease in seconds.
		DeploymentLeaseInSeconds *int

		// StorageLeaseInSeconds is the storage lease in seconds.
		StorageLeaseInSeconds *int
	}

	// ParamsDeleteVApp is the parameters for the DeleteVApp command.
	ParamsDeleteVApp struct {
		// ID is the unique identifier of the VApp to delete.
		ID string

		// Name is the name of the VApp to delete.
		Name string
	}
)

// Validate checks ParamsGetVApp structural constraints.
func (p ParamsGetVApp) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsCreateVApp structural constraints.
func (p ParamsCreateVApp) Validate() error {
	if p.Name == "" || p.VDCID == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsUpdateVApp structural constraints.
func (p ParamsUpdateVApp) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	if p.Description == nil && p.DeploymentLeaseInSeconds == nil && p.StorageLeaseInSeconds == nil {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsDeleteVApp structural constraints.
func (p ParamsDeleteVApp) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}
