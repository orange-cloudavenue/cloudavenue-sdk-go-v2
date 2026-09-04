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

// APIResponseListVApp contains VApp records returned by list operations.
type APIResponseListVApp struct {
	Records []APIResponseListVAppRecord `json:"record" fakesize:"2"`
}

// APIResponseListVAppRecord represents a VApp in list responses.
type APIResponseListVAppRecord struct {
	HREF        string `json:"href" fake:"{href_uuid}"`
	ID          string `json:"id" fake:"{urn:vapp}"`
	Name        string `json:"name" fake:"mockvapp-{word}"`
	Description string `json:"description" fake:"{sentence}"`
	Status      string `json:"status" fake:"{vapp_status}"`
	Deployed    bool   `json:"deployed" fake:"{bool}"`
	NumberOfVMs int    `json:"numberOfVMs" fake:"{number:0,10}"`
}

// APIResponseGetVApp contains a VApp returned by get operations.
type APIResponseGetVApp struct {
	HREF                 string                          `json:"href" fake:"{href_uuid}"`
	ID                   string                          `json:"id" fake:"{urn:vapp}"`
	Name                 string                          `json:"name" fake:"mockvapp-{word}"`
	Description          string                          `json:"description" fake:"{sentence}"`
	Status               string                          `json:"status" fake:"{vapp_status}"`
	Deployed             bool                            `json:"deployed" fake:"{bool}"`
	LeaseSettings        APIResponseLeaseSettings        `json:"leaseSettings"`
	GuestProperties      APIResponseGuestProperties      `json:"guestProperties"`
	NetworkConfigSection APIResponseNetworkConfigSection `json:"networkConfigSection"`
	VAppParent           string                          `json:"vAppParent" fake:"{urn:vapp}"`
	VDC                  string                          `json:"vdc" fake:"{urn:vdc}"`
	Org                  string                          `json:"org" fake:"{urn:org}"`
	Owner                string                          `json:"owner" fake:"{urn:user}"`
	Metadata             []APIResponseMetadata           `json:"metadataEntry" fakesize:"1"`
}

// APIResponseLeaseSettings contains VApp lease settings.
type APIResponseLeaseSettings struct {
	DeploymentLeaseInSeconds  int    `json:"deploymentLeaseInSeconds" fake:"{number:0,86400}"`
	StorageLeaseInSeconds     int    `json:"storageLeaseInSeconds" fake:"{number:0,86400}"`
	DeploymentLeaseExpiration string `json:"deploymentLeaseExpiration" fake:"{date}"`
	StorageLeaseExpiration    string `json:"storageLeaseExpiration" fake:"{date}"`
}

// APIResponseGuestProperties contains VApp guest properties.
type APIResponseGuestProperties struct {
	XMLNS                        string                      `json:"xmlns" fake:"http://www.vmware.com/vcloud/v1.5"`
	XMLNSXsi                     string                      `json:"xmlns:xsi" fake:"http://www.w3.org/2001/XMLSchema-instance"`
	XsiNoNamespaceSchemaLocation string                      `json:"xsi:noNamespaceSchemaLocation" fake:"http://www.vmware.com/vcloud/v1.5/vmguest.xsd"`
	ProductSectionList           []APIResponseProductSection `json:"ProductSectionList" fakesize:"1"`
}

// APIResponseProductSection contains product section info for guest properties.
type APIResponseProductSection struct {
	Info    string               `json:"Info" fake:"Product section"`
	Product []APIResponseProduct `json:"Product" fakesize:"1"`
}

// APIResponseProduct contains product info.
type APIResponseProduct struct {
	Info     string `json:"Info" fake:"Product info"`
	Class    string `json:"Class" fake:"{word}"`
	Instance string `json:"Instance" fake:"{word}"`
}

// APIResponseNetworkConfigSection contains network configuration.
type APIResponseNetworkConfigSection struct {
	XMLNS                        string                     `json:"xmlns" fake:"http://www.vmware.com/vcloud/v1.5"`
	XMLNSXsi                     string                     `json:"xmlns:xsi" fake:"http://www.w3.org/2001/XMLSchema-instance"`
	XsiNoNamespaceSchemaLocation string                     `json:"xsi:noNamespaceSchemaLocation" fake:"http://www.vmware.com/vcloud/v1.5/vAppNetworkConfig.xsd"`
	NetworkConfig                []APIResponseNetworkConfig `json:"networkConfig" fakesize:"1"`
}

// APIResponseNetworkConfig contains network configuration.
type APIResponseNetworkConfig struct {
	NetworkName   string                          `json:"networkName" fake:"{word}-net"`
	NetworkHREF   string                          `json:"network" fake:"{href_uuid}"`
	Configuration APIResponseNetworkConfiguration `json:"configuration"`
}

// APIResponseNetworkConfiguration contains network configuration details.
type APIResponseNetworkConfiguration struct {
	IPScope     APIResponseIPScope `json:"ipScope"`
	FenceMode   string             `json:"fenceMode" fake:"{vapp_fence_mode}"`
	IsConnected bool               `json:"isConnected" fake:"{bool}"`
}

// APIResponseIPScope contains IP scope configuration.
type APIResponseIPScope struct {
	IsInherited  bool                 `json:"isInherited" fake:"{bool}"`
	Gateway      string               `json:"gateway" fake:"{ipv4}"`
	Netmask      string               `json:"netmask" fake:"255.255.255.0"`
	PrefixLength int                  `json:"prefixLength" fake:"24"`
	DNSServer1   string               `json:"dns1" fake:"{ipv4}"`
	DNSServer2   string               `json:"dns2" fake:"{ipv4}"`
	DNSSuffix    string               `json:"dnsSuffix" fake:"example.com"`
	IPRanges     []APIResponseIPRange `json:"ipRanges" fakesize:"1"`
}

// APIResponseIPRange contains IP range configuration.
type APIResponseIPRange struct {
	StartAddress string `json:"startAddress" fake:"{ipv4}"`
	EndAddress   string `json:"endAddress" fake:"{ipv4}"`
}

// APIResponseMetadata contains metadata entries.
type APIResponseMetadata struct {
	Key   string `json:"key" fake:"{word}"`
	Value string `json:"value" fake:"{word}"`
	Type  string `json:"type" fake:"string"`
}

// APIRequestCreateVApp contains the request payload for creating a VApp.
type APIRequestCreateVApp struct {
	XMLNS       string `json:"xmlns" fake:"http://www.vmware.com/vcloud/v1.5"`
	Name        string `json:"name" fake:"mockvapp-{word}"`
	Description string `json:"description,omitempty" fake:"{sentence}"`
}

// APIRequestUpdateVApp contains the request payload for updating a VApp.
type APIRequestUpdateVApp struct {
	XMLNS         string                   `json:"xmlns" fake:"http://www.vmware.com/vcloud/v1.5"`
	Name          string                   `json:"name,omitempty" fake:"mockvapp-{word}"`
	Description   string                   `json:"description,omitempty" fake:"{sentence}"`
	LeaseSettings *APIRequestLeaseSettings `json:"leaseSettings,omitempty"`
}

// APIRequestLeaseSettings contains lease settings for update requests.
type APIRequestLeaseSettings struct {
	DeploymentLeaseInSeconds *int `json:"deploymentLeaseInSeconds,omitempty"`
	StorageLeaseInSeconds    *int `json:"storageLeaseInSeconds,omitempty"`
}

// APIRequestUndeployVApp contains the request payload for undeploying a VApp.
type APIRequestUndeployVApp struct {
	XMLNS            string `json:"xmlns" fake:"http://www.vmware.com/vcloud/v1.5"`
	UndeployPowerOff bool   `json:"undeployPowerOff" fake:"{bool}"`
}

// ToModel converts APIResponseListVAppRecord to types.ModelVApp.
func (r *APIResponseListVAppRecord) ToModel() types.ModelVApp {
	return types.ModelVApp{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		Deployed:    r.Deployed,
		NumberOfVMs: r.NumberOfVMs,
	}
}

// ToModel converts APIResponseGetVApp to types.ModelVApp.
func (r *APIResponseGetVApp) ToModel() types.ModelVApp {
	m := types.ModelVApp{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		Deployed:    r.Deployed,
		HREF:        r.HREF,
		VDC:         r.VDC,
		Org:         r.Org,
		Owner:       r.Owner,
		Parent:      r.VAppParent,
	}

	if r.LeaseSettings.DeploymentLeaseInSeconds > 0 {
		m.DeploymentLeaseInSeconds = &r.LeaseSettings.DeploymentLeaseInSeconds
	}
	if r.LeaseSettings.StorageLeaseInSeconds > 0 {
		m.StorageLeaseInSeconds = &r.LeaseSettings.StorageLeaseInSeconds
	}

	for _, meta := range r.Metadata {
		if meta.Key == "vappBillingModel" {
			m.Properties.BillingModel = meta.Value
		}
	}

	return m
}

// ToModel converts APIResponseListVApp to types.ModelListVApp.
func (r *APIResponseListVApp) ToModel() *types.ModelListVApp {
	model := &types.ModelListVApp{
		VApps: make([]types.ModelVApp, 0),
	}

	for _, vapp := range r.Records {
		model.VApps = append(model.VApps, vapp.ToModel())
	}

	return model
}
