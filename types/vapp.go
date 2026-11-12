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

		Properties      ModelVAppProperties      `documentation:"Properties of the VApp"`
		GuestProperties ModelVAppGuestProperties `documentation:"Guest properties of the VApp"`
	}

	// ModelVAppProperties contains VApp properties.
	ModelVAppProperties struct {
		BillingModel string `documentation:"Billing model of the VApp"`
	}

	// ModelVAppGuestProperties contains vApp ProductSection guest properties.
	ModelVAppGuestProperties struct {
		Sections []ModelVAppGuestPropertySection
	}

	// ModelVAppGuestPropertySection contains one ProductSection and its properties.
	ModelVAppGuestPropertySection struct {
		Class      string
		Instance   string
		Properties map[string]string
	}

	// ModelVAppRuntimeNetworkConfigSection contains the runtime networkConfigSection returned by VCD.
	ModelVAppRuntimeNetworkConfigSection struct {
		Info          string
		NetworkConfig []ModelVAppRuntimeNetworkConfig
	}

	// ModelVAppRuntimeNetworkConfig contains one runtime network definition of a vApp.
	ModelVAppRuntimeNetworkConfig struct {
		ID            string
		HREF          string
		Description   string
		Type          string
		NetworkName   string
		Configuration ModelVAppRuntimeNetworkConfiguration
	}

	// ModelVAppRuntimeNetworkConfiguration contains the runtime settings of one vApp network.
	ModelVAppRuntimeNetworkConfiguration struct {
		ParentNetwork                  *VAppReference
		FenceMode                      string
		RetainNetInfoAcrossDeployments *bool
		IsDeployed                     *bool
		IsShared                       *bool
		IsConnected                    bool
		IPScopes                       []ModelVAppRuntimeIPScope
	}

	// ModelVAppRuntimeIPScope contains one IP scope of a runtime vApp network.
	ModelVAppRuntimeIPScope struct {
		IsInherited  bool
		Gateway      string
		Netmask      string
		DNS1         string
		DNS2         string
		DNSSuffix    string
		PrefixLength *int
		IPRanges     []ModelVAppRuntimeIPRange
	}

	// ModelVAppRuntimeIPRange contains one IP range of a runtime vApp network.
	ModelVAppRuntimeIPRange struct {
		StartAddress string
		EndAddress   string
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

	// VAppUndeployMode defines the VMware undeploy action payload mode.
	VAppUndeployMode string

	// ParamsGetVAppGuestProperties identifies the vApp guest properties to read.
	ParamsGetVAppGuestProperties struct {
		ID   string
		Name string
	}

	// ParamsUpdateVAppGuestProperties identifies the vApp guest properties to replace.
	ParamsUpdateVAppGuestProperties struct {
		ID              string
		Name            string
		GuestProperties ModelVAppGuestProperties
	}

	// ParamsUndeployVApp defines the parameters for undeploying a vApp.
	ParamsUndeployVApp struct {
		ID   string
		Name string
		Mode VAppUndeployMode
	}

	// ParamsGetVAppNetworkConfig identifies the vApp runtime network section to read.
	ParamsGetVAppNetworkConfig struct {
		ID   string
		Name string
	}

	// ParamsCreateVAppNetwork defines a new isolated runtime vApp network.
	ParamsCreateVAppNetwork struct {
		VAppID   string
		VAppName string
		Network  ModelVAppRuntimeNetworkConfig
	}

	// ParamsAddOrgNetworkToVApp defines a new direct/routed runtime vApp network backed by an org network.
	ParamsAddOrgNetworkToVApp struct {
		VAppID   string
		VAppName string
		Network  ModelVAppRuntimeNetworkConfig
	}

	// ParamsUpdateVAppNetwork updates an existing runtime vApp network by name.
	ParamsUpdateVAppNetwork struct {
		VAppID      string
		VAppName    string
		NetworkName string
		Network     ModelVAppRuntimeNetworkConfig
	}

	// ParamsRemoveVAppNetwork removes one runtime vApp network by name.
	ParamsRemoveVAppNetwork struct {
		VAppID      string
		VAppName    string
		NetworkName string
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

const (
	// VAppUndeployModePowerOff powers off the vApp before undeploy.
	VAppUndeployModePowerOff VAppUndeployMode = "powerOff"
	// VAppUndeployModeSuspend suspends the vApp before undeploy.
	VAppUndeployModeSuspend VAppUndeployMode = "suspend"
	// VAppUndeployModeForce requests VMware force undeploy behavior.
	VAppUndeployModeForce VAppUndeployMode = "force"
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

// Validate checks ParamsGetVAppGuestProperties structural constraints.
func (p ParamsGetVAppGuestProperties) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsUpdateVAppGuestProperties structural constraints.
func (p ParamsUpdateVAppGuestProperties) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	return p.GuestProperties.Validate()
}

// Validate checks ParamsUndeployVApp structural constraints.
func (p ParamsUndeployVApp) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	return p.Mode.Validate()
}

// Validate checks ParamsGetVAppNetworkConfig structural constraints.
func (p ParamsGetVAppNetworkConfig) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsCreateVAppNetwork structural constraints.
func (p ParamsCreateVAppNetwork) Validate() error {
	if p.VAppID == "" && p.VAppName == "" {
		return fmt.Errorf("missing required parameters")
	}
	if p.Network.Configuration.ParentNetwork != nil {
		return fmt.Errorf("parent network is not allowed for create")
	}
	if p.Network.Configuration.FenceMode != "isolated" {
		return fmt.Errorf("create requires isolated fence mode")
	}
	return p.Network.Validate()
}

// Validate checks ParamsAddOrgNetworkToVApp structural constraints.
func (p ParamsAddOrgNetworkToVApp) Validate() error {
	if p.VAppID == "" && p.VAppName == "" {
		return fmt.Errorf("missing required parameters")
	}
	if err := p.Network.Validate(); err != nil {
		return err
	}
	if p.Network.Configuration.ParentNetwork == nil {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks ParamsUpdateVAppNetwork structural constraints.
func (p ParamsUpdateVAppNetwork) Validate() error {
	if p.VAppID == "" && p.VAppName == "" {
		return fmt.Errorf("missing required parameters")
	}
	if p.NetworkName == "" {
		return fmt.Errorf("missing required parameters")
	}
	if p.Network.NetworkName != "" && p.Network.NetworkName != p.NetworkName {
		return fmt.Errorf("network name mismatch")
	}
	return p.Network.Validate()
}

// Validate checks ParamsRemoveVAppNetwork structural constraints.
func (p ParamsRemoveVAppNetwork) Validate() error {
	if p.VAppID == "" && p.VAppName == "" {
		return fmt.Errorf("missing required parameters")
	}
	if p.NetworkName == "" {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks runtime vApp network structural constraints.
func (m ModelVAppRuntimeNetworkConfig) Validate() error {
	if m.NetworkName == "" {
		return fmt.Errorf("missing required parameters")
	}
	return m.Configuration.Validate()
}

// Validate checks guest properties structural constraints.
func (m ModelVAppGuestProperties) Validate() error {
	for _, section := range m.Sections {
		if err := section.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks guest property section structural constraints.
func (m ModelVAppGuestPropertySection) Validate() error {
	if m.Properties == nil {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

// Validate checks undeploy mode structural constraints.
func (m VAppUndeployMode) Validate() error {
	switch m {
	case "", VAppUndeployModePowerOff, VAppUndeployModeSuspend, VAppUndeployModeForce:
		return nil
	default:
		return fmt.Errorf("unsupported undeploy mode %q", m)
	}
}

// Validate checks runtime vApp network configuration structural constraints.
func (m ModelVAppRuntimeNetworkConfiguration) Validate() error {
	if m.FenceMode == "" {
		return fmt.Errorf("missing required parameters")
	}
	if len(m.IPScopes) == 0 && m.ParentNetwork == nil {
		return fmt.Errorf("missing required parameters")
	}
	for _, scope := range m.IPScopes {
		if err := scope.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks runtime IP scope structural constraints.
func (m ModelVAppRuntimeIPScope) Validate() error {
	if m.Gateway == "" || m.Netmask == "" {
		return fmt.Errorf("missing required parameters")
	}
	for _, r := range m.IPRanges {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks runtime IP range structural constraints.
func (m ModelVAppRuntimeIPRange) Validate() error {
	if m.StartAddress == "" || m.EndAddress == "" {
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
