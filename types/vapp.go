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

	// VAppReference identifies a VCD resource used by a creation request.
	VAppReference struct {
		HREF string
		ID   string
		Type string
		Name string
	}

	// VAppSourcedItem describes the source item included by ComposeVApp.
	VAppSourcedItem struct {
		Source            VAppReference
		SourceDelete      bool
		VAppScopedLocalID string
		StorageProfile    *VAppReference
	}

	// VAppInstantiationParams contains optional OVF network sections. More
	// sections can be added without conflating compose and template requests.
	VAppInstantiationParams struct {
		NetworkConfigSection     *VAppNetworkConfigSection
		NetworkConnectionSection *VAppNetworkConnectionSection
	}

	VAppNetworkConfigSection struct {
		Info          string
		NetworkConfig []VAppNetworkConfig
	}

	VAppNetworkConfig struct {
		NetworkName string
		Network     string
	}

	VAppNetworkConnectionSection struct {
		Info                          string
		PrimaryNetworkConnectionIndex int
	}

	// ParamsComposeVApp is the parameters for the documented ComposeVApp
	// operation. Compose accepts a sourced item and optional OVF sections.
	ParamsComposeVApp struct {
		// Name is the name of the VApp to compose.
		Name string

		// Description is the description of the VApp to create.
		Description string

		// VDCID is the VDC ID where the VApp will be composed.
		VDCID string

		Deploy              bool
		PowerOn             bool
		LinkedClone         bool
		AllEULAsAccepted    bool
		SourcedItem         *VAppSourcedItem
		InstantiationParams *VAppInstantiationParams
	}

	// ParamsInstantiateVAppTemplate is the parameters for instantiating a VApp
	// from a catalog/template resource.
	ParamsInstantiateVAppTemplate struct {
		Name        string
		Description string
		VDCID       string
		TemplateID  string

		Deploy              bool
		PowerOn             bool
		LinkedClone         bool
		AllEULAsAccepted    bool
		InstantiationParams *VAppInstantiationParams
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

	// ModelVAppLeaseSettings contains org/VDC-level lease settings.
	ModelVAppLeaseSettings struct {
		DeploymentLeaseInSeconds int `documentation:"Deployment lease in seconds"`
		StorageLeaseInSeconds    int `documentation:"Storage lease in seconds"`
	}

	// ParamsGetVAppLeaseSettings defines the parameters for getting org/VDC-level lease settings.
	ParamsGetVAppLeaseSettings struct {
		// OrgID is the unique identifier of the organization.
		// Example: urn:vcloud:org:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		OrgID string

		// OrgName is the name of the organization.
		OrgName string
	}

	// ParamsUpdateVAppLeaseSettings defines the parameters for updating org/VDC-level lease settings.
	ParamsUpdateVAppLeaseSettings struct {
		// OrgID is the unique identifier of the organization.
		// Example: urn:vcloud:org:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		OrgID string

		// OrgName is the name of the organization.
		OrgName string

		// LeaseSettings are the lease settings to apply.
		LeaseSettings ModelVAppLeaseSettings
	}
)

// Validate checks ParamsGetVApp structural constraints.
func (p ParamsGetVApp) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsComposeVApp structural constraints.
func (p ParamsComposeVApp) Validate() error {
	if p.Name == "" || p.VDCID == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsInstantiateVAppTemplate structural constraints.
func (p ParamsInstantiateVAppTemplate) Validate() error {
	if p.Name == "" || p.VDCID == "" || p.TemplateID == "" {
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

// Validate checks ParamsGetVAppLeaseSettings structural constraints.
func (p ParamsGetVAppLeaseSettings) Validate() error {
	if p.OrgID == "" && p.OrgName == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsUpdateVAppLeaseSettings structural constraints.
func (p ParamsUpdateVAppLeaseSettings) Validate() error {
	if p.OrgID == "" && p.OrgName == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}
