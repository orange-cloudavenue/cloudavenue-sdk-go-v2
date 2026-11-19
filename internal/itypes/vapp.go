/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// APIResponseListVApp contains VApp records returned by list operations.
type APIResponseListVApp struct {
	Records []APIResponseListVAppRecord `json:"record" fakesize:"2"`
}

// APIResponseListVAppRecord represents a VApp in list responses.
type APIResponseListVAppRecord struct {
	HREF        string         `json:"href" fake:"{href_uuid}"`
	ID          string         `json:"id" fake:"{urn:vapp}"`
	Name        string         `json:"name" fake:"mockvapp-{word}"`
	Description string         `json:"description" fake:"{sentence}"`
	Status      VAppStatusCode `json:"status" fake:"{vapp_status}"`
	Deployed    bool           `json:"deployed" fake:"{bool}"`
	NumberOfVMs int            `json:"numberOfVMs" fake:"{number:0,10}"`
}

// VAppStatusCode accepts the legacy numeric wire encoding (bare JSON number)
// as well as a stringified code, and never fails on unknown values.
type VAppStatusCode string

// UnmarshalJSON accepts both JSON numbers and strings.
func (c *VAppStatusCode) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == typeNull {
		return nil
	}
	if data[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		*c = VAppStatusCode(value)
		return nil
	}
	*c = VAppStatusCode(data)
	return nil
}

// VAppOwnerRef accepts the legacy owner representation: either a plain
// string reference or the nested owner object returned by the JSON API.
type VAppOwnerRef string

// UnmarshalJSON normalizes string and owner-object payloads down to a
// single reference, preferring the owning user identifier.
func (o *VAppOwnerRef) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		*o = VAppOwnerRef(value)
		return nil
	}
	var payload struct {
		HREF string `json:"href"`
		User struct {
			HREF string `json:"href"`
			ID   string `json:"id"`
		} `json:"user"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	switch {
	case payload.User.ID != "":
		*o = VAppOwnerRef(payload.User.ID)
	case payload.User.HREF != "":
		*o = VAppOwnerRef(payload.User.HREF)
	default:
		*o = VAppOwnerRef(payload.HREF)
	}
	return nil
}

// vAppStatusCodes maps the deprecated numeric vApp status codes to their
// semantic tokens, matching the VMware enumeration (govcd types.VAppStatuses).
var vAppStatusCodes = map[string]string{
	"-1": "FAILED_CREATION", "0": "UNRESOLVED", "1": "RESOLVED", "2": "DEPLOYED",
	"3": "SUSPENDED", "4": "POWERED_ON", "5": "WAITING_FOR_INPUT", "6": "UNKNOWN",
	"7": "UNRECOGNIZED", "8": "POWERED_OFF", "9": "INCONSISTENT_STATE", "10": "MIXED",
	"11": "DESCRIPTOR_PENDING", "12": "COPYING_CONTENTS", "13": "DISK_CONTENTS_PENDING",
	"14": "QUARANTINED", "15": "QUARANTINE_EXPIRED", "16": "REJECTED",
	"17": "TRANSFER_TIMEOUT", "18": "VAPP_UNDEPLOYED", "19": "VAPP_PARTIALLY_DEPLOYED",
	"20": "PARTIALLY_POWERED_OFF", "21": "PARTIALLY_SUSPENDED",
}

// normalizeVAppStatus maps a numeric wire code to its semantic token.
// Unknown and already-semantic values are returned unchanged.
func normalizeVAppStatus(raw string) string {
	if token, ok := vAppStatusCodes[raw]; ok {
		return token
	}
	return raw
}

// APIResponseGetVApp contains a VApp returned by get operations.
type APIResponseGetVApp struct {
	HREF            string                     `json:"href" fake:"{href_uuid}"`
	ID              string                     `json:"id" fake:"{urn:vapp}"`
	Name            string                     `json:"name" fake:"mockvapp-{word}"`
	Description     string                     `json:"description" fake:"{sentence}"`
	Status          VAppStatusCode             `json:"status" fake:"{vapp_status}"`
	Deployed        bool                       `json:"deployed" fake:"{bool}"`
	LeaseSettings   APIResponseLeaseSettings   `json:"leaseSettings"`
	GuestProperties APIResponseGuestProperties `json:"guestProperties"`
	VAppParent      string                     `json:"vAppParent" fake:"{urn:vapp}"`
	VDC             string                     `json:"vdc" fake:"{urn:vdc}"`
	Org             string                     `json:"org" fake:"{urn:org}"`
	Owner           VAppOwnerRef               `json:"owner" fake:"{urn:user}"`
	Metadata        []APIResponseMetadata      `json:"metadataEntry" fakesize:"1"`
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
	ProductSectionList           []APIResponseProductSection `json:"ProductSectionList" fakesize:"1"` //nolint:tagliatelle
}

// APIResponseProductSection contains product section info for guest properties.
type APIResponseProductSection struct {
	Info     string                       `json:"Info" fake:"Product section"` //nolint:tagliatelle
	Class    string                       `json:"Class" fake:"{word}"`         //nolint:tagliatelle
	Instance string                       `json:"Instance" fake:"{word}"`      //nolint:tagliatelle
	Property []APIResponseProductProperty `json:"Property" fakesize:"1"`       //nolint:tagliatelle
}

// APIResponseProductProperty contains one guest property from a ProductSection.
type APIResponseProductProperty struct {
	Key              string `json:"key" fake:"{word}"`
	Label            string `json:"label" fake:"{word}"`
	Type             string `json:"type" fake:"{word}"`
	Value            string `json:"value" fake:"{word}"`
	UserConfigurable bool   `json:"userConfigurable" fake:"{bool}"`
}

// APIResponseVAppNetworkConfigSection contains the runtime XML networkConfigSection of a vApp.
type APIResponseVAppNetworkConfigSection struct {
	XMLName       xml.Name                       `xml:"NetworkConfigSection"`
	XMLNS         string                         `xml:"xmlns,attr,omitempty"`
	OVF           string                         `xml:"xmlns:ovf,attr,omitempty"`
	Info          string                         `xml:"ovf:Info,omitempty"`
	NetworkConfig []APIResponseVAppNetworkConfig `xml:"NetworkConfig,omitempty"`
}

type APIResponseVAppNetworkConfig struct {
	NetworkID     string                              `xml:"id,attr,omitempty"`
	NetworkHREF   string                              `xml:"href,attr,omitempty"`
	Description   string                              `xml:"Description,omitempty"`
	NetworkType   string                              `xml:"type,attr,omitempty"`
	NetworkName   string                              `xml:"networkName,attr,omitempty"`
	Configuration APIResponseVAppRuntimeConfiguration `xml:"Configuration"`
}

type APIResponseVAppRuntimeConfiguration struct {
	// Field order mirrors the XSD sequence of NetworkConfigurationType as
	// enforced by Cloud Director (govcd types.NetworkConfiguration is the
	// production reference): IpScopes first, then ParentNetwork, FenceMode,
	// RetainNetInfoAcrossDeployments, Features, RouterInfo, GuestVlanAllowed
	// and IsDeployed. Out-of-order elements are rejected with HTTP 400
	// cvc-complex-type.2.4.a.
	IPScopes                       *APIResponseVAppRuntimeIPScopes   `xml:"IpScopes,omitempty"`
	ParentNetwork                  *APIRequestVAppReference          `xml:"ParentNetwork,omitempty"`
	FenceMode                      string                            `xml:"FenceMode,omitempty"`
	RetainNetInfoAcrossDeployments *bool                             `xml:"RetainNetInfoAcrossDeployments,omitempty"`
	Features                       *APIResponseVAppRuntimeFeatures   `xml:"Features,omitempty"`
	RouterInfo                     *APIResponseVAppRuntimeRouterInfo `xml:"RouterInfo,omitempty"`
	GuestVLANAllowed               *bool                             `xml:"GuestVlanAllowed,omitempty"`
	IsDeployed                     *bool                             `xml:"IsDeployed,omitempty"`
}

type APIResponseVAppRuntimeFeatures struct {
	FirewallService      *APIResponseVAppRuntimeInnerXML `xml:"FirewallService,omitempty"`
	NATService           *APIResponseVAppRuntimeInnerXML `xml:"NatService,omitempty"`
	DhcpService          *APIResponseVAppRuntimeInnerXML `xml:"DhcpService,omitempty"`
	StaticRoutingService *APIResponseVAppRuntimeInnerXML `xml:"StaticRoutingService,omitempty"`
}

type APIResponseVAppRuntimeInnerXML struct {
	InnerXML string `xml:",innerxml"`
}

type APIResponseVAppRuntimeRouterInfo struct {
	ExternalIP               string `xml:"ExternalIp,omitempty"`
	ExternalIPAllocationMode string `xml:"ExternalIpAllocationMode,omitempty"`
	ConnectionType           string `xml:"ConnectionType,omitempty"`
	ParentNetwork            string `xml:"ParentNetwork,omitempty"`
	ParentNetworkName        string `xml:"ParentNetworkName,omitempty"`
}

type APIResponseVAppRuntimeIPScopes struct {
	// The vApp networkConfigSection wraps IpScope elements in <IpScopes>,
	// as govcd types.IPScopes does. The wire constraints are the enclosing
	// element order (IpScopes first in Configuration) and the IpScope child
	// order, not removal of the wrapper.
	IPScope []APIResponseVAppRuntimeIPScope `xml:"IpScope,omitempty"`
}

type APIResponseVAppRuntimeIPScope struct {
	// Field order and element names follow the IpScopeType XSD sequence
	// (govcd types.IPScope): the legacy wire element is SubnetPrefixLength.
	IsInherited  *bool                          `xml:"IsInherited,omitempty"`
	Gateway      string                         `xml:"Gateway,omitempty"`
	Netmask      string                         `xml:"Netmask,omitempty"`
	PrefixLength *int                           `xml:"SubnetPrefixLength,omitempty"`
	DNS1         string                         `xml:"Dns1,omitempty"`
	DNS2         string                         `xml:"Dns2,omitempty"`
	DNSSuffix    string                         `xml:"DnsSuffix,omitempty"`
	IPRanges     APIResponseVAppRuntimeIPRanges `xml:"IpRanges,omitempty"`
}

type APIResponseVAppRuntimeIPRanges struct {
	IPRange []APIResponseVAppRuntimeIPRange `xml:"IpRange,omitempty"`
}

type APIResponseVAppRuntimeIPRange struct {
	StartAddress string `xml:"StartAddress,omitempty"`
	EndAddress   string `xml:"EndAddress,omitempty"`
}

// APIResponseMetadata contains metadata entries.
type APIResponseMetadata struct {
	Key   string `json:"key" fake:"{word}"`
	Value string `json:"value" fake:"{word}"`
	Type  string `json:"type" fake:"string"`
}

// APIRequestVAppReference identifies a VCD resource in an XML request.
type APIRequestVAppReference struct {
	HREF string `xml:"href,attr,omitempty"`
	ID   string `xml:"id,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	Name string `xml:"name,attr,omitempty"`
}

// APIRequestVAppInstantiationParams carries optional OVF instantiation
// sections. Sections are deliberately represented explicitly instead of using
// the old name/description-only JSON payload.
type APIRequestVAppInstantiationParams struct {
	NetworkConfigSection     *APIRequestVAppNetworkConfigSection     `xml:"NetworkConfigSection,omitempty"`
	NetworkConnectionSection *APIRequestVAppNetworkConnectionSection `xml:"NetworkConnectionSection,omitempty"`
}

type APIRequestVAppNetworkConfigSection struct {
	Info          string                        `xml:"ovf:Info,omitempty"`
	NetworkConfig []APIRequestVAppNetworkConfig `xml:"NetworkConfig,omitempty"`
}

// APIRequestVAppNetworkConfig mirrors the XSD NetworkConfig: only the
// networkName attribute plus a mandatory Configuration block (ParentNetwork
// and FenceMode) are wire-valid; a bare network attribute is rejected.
type APIRequestVAppNetworkConfig struct {
	NetworkName   string                             `xml:"networkName,attr,omitempty"`
	Description   string                             `xml:"Description,omitempty"`
	Configuration APIRequestVAppNetworkConfiguration `xml:"Configuration"`
}

type APIRequestVAppNetworkConfiguration struct {
	ParentNetwork *APIRequestVAppReference `xml:"ParentNetwork,omitempty"`
	FenceMode     string                   `xml:"FenceMode"`
}

type APIRequestVAppNetworkConnectionSection struct {
	Info                          string `xml:"ovf:Info,omitempty"`
	PrimaryNetworkConnectionIndex int    `xml:"PrimaryNetworkConnectionIndex,omitempty"`
}

// APIRequestVAppSourcedItem describes the source included by ComposeVApp.
type APIRequestVAppSourcedItem struct {
	Source            APIRequestVAppReference  `xml:"Source"`
	SourceDelete      bool                     `xml:"sourceDelete,attr,omitempty"`
	VAppScopedLocalID string                   `xml:"VAppScopedLocalId,omitempty"`
	StorageProfile    *APIRequestVAppReference `xml:"StorageProfile,omitempty"`
}

// APIRequestComposeVApp is the documented ComposeVApp XML request.
type APIRequestComposeVApp struct {
	XMLName             xml.Name                           `xml:"ComposeVAppParams"`
	OVF                 string                             `xml:"xmlns:ovf,attr"`
	XSI                 string                             `xml:"xmlns:xsi,attr"`
	XMLNS               string                             `xml:"xmlns,attr"`
	Name                string                             `xml:"name,attr"`
	Deploy              bool                               `xml:"deploy,attr"`
	PowerOn             bool                               `xml:"powerOn,attr"`
	LinkedClone         bool                               `xml:"linkedClone,attr,omitempty"`
	Description         string                             `xml:"Description,omitempty"`
	InstantiationParams *APIRequestVAppInstantiationParams `xml:"InstantiationParams,omitempty"`
	SourcedItem         *APIRequestVAppSourcedItem         `xml:"SourcedItem,omitempty"`
	// Element required last by the ComposeVAppParams sequence; the same-name
	// attribute form is XSD-rejected. Omitted for empty vApps.
	AllEULAsAccepted bool `xml:"AllEULAsAccepted,omitempty"`
}

// APIRequestInstantiateVAppTemplate is the documented
// InstantiateVAppTemplate XML request.
type APIRequestInstantiateVAppTemplate struct {
	XMLName     xml.Name `xml:"InstantiateVAppTemplateParams"`
	OVF         string   `xml:"xmlns:ovf,attr"`
	XSI         string   `xml:"xmlns:xsi,attr,omitempty"`
	XMLNS       string   `xml:"xmlns,attr"`
	Name        string   `xml:"name,attr"`
	Deploy      bool     `xml:"deploy,attr"`
	PowerOn     bool     `xml:"powerOn,attr"`
	LinkedClone bool     `xml:"linkedClone,attr,omitempty"`
	Description string   `xml:"Description,omitempty"`
	// Field order is wire-significant: the schema sequence is
	// Description, InstantiationParams, Source (per the official
	// instantiateVAppTemplate example); InstantiationParams after Source
	// fails XSD validation.
	InstantiationParams *APIRequestVAppInstantiationParams `xml:"InstantiationParams,omitempty"`
	Source              APIRequestVAppReference            `xml:"Source"`
	IsSourceDelete      bool                               `xml:"IsSourceDelete,omitempty"`
	AllEULAsAccepted    bool                               `xml:"AllEULAsAccepted,omitempty"`
}

// APIRequestUpdateVApp contains the request payload for updating a VApp.
type APIRequestUpdateVApp struct {
	XMLNS         string                   `json:"xmlns" fake:"http://www.vmware.com/vcloud/v1.5"`
	Name          string                   `json:"name,omitempty" fake:"mockvapp-{word}"`
	Description   *string                  `json:"description,omitempty" fake:"{sentence}"`
	LeaseSettings *APIRequestLeaseSettings `json:"leaseSettings,omitempty"`
}

// APIRequestLeaseSettings contains lease settings for update requests.
type APIRequestLeaseSettings struct {
	DeploymentLeaseInSeconds *int `json:"deploymentLeaseInSeconds,omitempty"`
	StorageLeaseInSeconds    *int `json:"storageLeaseInSeconds,omitempty"`
}

// APIResponseOrgLeaseSettings contains org/VDC-level lease settings.
type APIResponseOrgLeaseSettings struct {
	DeploymentLeaseInSeconds int `json:"deploymentLeaseInSeconds"`
	StorageLeaseInSeconds    int `json:"storageLeaseInSeconds"`
}

// APIRequestOrgLeaseSettings is the request payload for updating org/VDC-level lease settings.
type APIRequestOrgLeaseSettings struct {
	DeploymentLeaseInSeconds int `json:"deploymentLeaseInSeconds"`
	StorageLeaseInSeconds    int `json:"storageLeaseInSeconds"`
}

// APIRequestUndeployVApp contains the request payload for undeploying a VApp.
type APIRequestUndeployVApp struct {
	XMLName xml.Name `xml:"UndeployVAppParams"`
	XMLNS   string   `xml:"xmlns,attr,omitempty" fake:"http://www.vmware.com/vcloud/v1.5"`
	// The 39.1 schema replaced the legacy UndeployPowerOff boolean with this.
	UndeployAction string `xml:"UndeployPowerAction,omitempty"`
}

// ToModel converts APIResponseListVAppRecord to types.ModelVApp.
func (r *APIResponseListVAppRecord) ToModel() types.ModelVApp {
	return types.ModelVApp{
		ID:          vAppRecordID(r.ID, r.HREF),
		Name:        r.Name,
		Description: r.Description,
		Status:      normalizeVAppStatus(string(r.Status)),
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
		Status:      normalizeVAppStatus(string(r.Status)),
		Deployed:    r.Deployed,
		HREF:        r.HREF,
		VDC:         r.VDC,
		Org:         r.Org,
		Owner:       string(r.Owner),
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
	m.GuestProperties = r.GuestProperties.ToModel()

	return m
}

// ToModel converts APIResponseGuestProperties to the public guest property model.
func (r *APIResponseGuestProperties) ToModel() types.ModelVAppGuestProperties {
	model := types.ModelVAppGuestProperties{Sections: make([]types.ModelVAppGuestPropertySection, 0, len(r.ProductSectionList))}
	for _, section := range r.ProductSectionList {
		model.Sections = append(model.Sections, section.ToModel())
	}
	return model
}

// ToModel converts APIResponseProductSection to the public guest property section model.
func (r *APIResponseProductSection) ToModel() types.ModelVAppGuestPropertySection {
	model := types.ModelVAppGuestPropertySection{
		Class:      r.Class,
		Instance:   r.Instance,
		Properties: make(map[string]types.ModelVAppGuestProperty, len(r.Property)),
	}
	for _, property := range r.Property {
		model.Properties[property.Key] = types.ModelVAppGuestProperty{
			Key:              property.Key,
			Label:            property.Label,
			Type:             property.Type,
			Value:            property.Value,
			UserConfigurable: property.UserConfigurable,
		}
	}
	return model
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

// ToModel converts APIResponseVAppNetworkConfigSection to the public runtime model.
func (r *APIResponseVAppNetworkConfigSection) ToModel() *types.ModelVAppRuntimeNetworkConfigSection {
	if r == nil {
		return nil
	}
	model := &types.ModelVAppRuntimeNetworkConfigSection{Info: r.Info, NetworkConfig: make([]types.ModelVAppRuntimeNetworkConfig, 0, len(r.NetworkConfig))}
	for _, network := range r.NetworkConfig {
		model.NetworkConfig = append(model.NetworkConfig, network.ToModel())
	}
	return model
}

// ToModel converts APIResponseVAppNetworkConfig to the public runtime model.
func (r *APIResponseVAppNetworkConfig) ToModel() types.ModelVAppRuntimeNetworkConfig {
	return types.ModelVAppRuntimeNetworkConfig{
		Description:   r.Description,
		ID:            r.NetworkID,
		HREF:          r.NetworkHREF,
		Type:          r.NetworkType,
		NetworkName:   r.NetworkName,
		Configuration: r.Configuration.ToModel(),
	}
}

// ToModel converts APIResponseVAppRuntimeConfiguration to the public runtime model.
func (r *APIResponseVAppRuntimeConfiguration) ToModel() types.ModelVAppRuntimeNetworkConfiguration {
	model := types.ModelVAppRuntimeNetworkConfiguration{
		FenceMode:                      r.FenceMode,
		RetainNetInfoAcrossDeployments: r.RetainNetInfoAcrossDeployments,
		IsDeployed:                     r.IsDeployed,
		GuestVLANAllowed:               r.GuestVLANAllowed,
		IPScopes:                       make([]types.ModelVAppRuntimeIPScope, 0),
	}
	if r.ParentNetwork != nil {
		ref := types.VAppReference{HREF: r.ParentNetwork.HREF, ID: r.ParentNetwork.ID, Type: r.ParentNetwork.Type, Name: r.ParentNetwork.Name}
		model.ParentNetwork = &ref
	}
	if r.IPScopes != nil && len(r.IPScopes.IPScope) > 0 {
		model.IPScopes = make([]types.ModelVAppRuntimeIPScope, 0, len(r.IPScopes.IPScope))
		for _, scope := range r.IPScopes.IPScope {
			model.IPScopes = append(model.IPScopes, scope.ToModel())
		}
	}
	return model
}

func (r *APIResponseVAppRuntimeIPScope) ToModel() types.ModelVAppRuntimeIPScope {
	model := types.ModelVAppRuntimeIPScope{
		IsInherited:  r.IsInherited != nil && *r.IsInherited,
		Gateway:      r.Gateway,
		Netmask:      r.Netmask,
		DNS1:         r.DNS1,
		DNS2:         r.DNS2,
		DNSSuffix:    r.DNSSuffix,
		PrefixLength: r.PrefixLength,
		IPRanges:     make([]types.ModelVAppRuntimeIPRange, 0, len(r.IPRanges.IPRange)),
	}
	for _, ipRange := range r.IPRanges.IPRange {
		model.IPRanges = append(model.IPRanges, ipRange.ToModel())
	}
	return model
}

func (r *APIResponseVAppRuntimeIPRange) ToModel() types.ModelVAppRuntimeIPRange {
	return types.ModelVAppRuntimeIPRange{StartAddress: r.StartAddress, EndAddress: r.EndAddress}
}

// ToAPIResponseVAppNetworkConfigSection converts the public runtime model to the runtime XML payload.
func ToAPIResponseVAppNetworkConfigSection(model *types.ModelVAppRuntimeNetworkConfigSection) *APIResponseVAppNetworkConfigSection {
	if model == nil {
		return nil
	}
	section := &APIResponseVAppNetworkConfigSection{
		XMLNS:         "http://www.vmware.com/vcloud/v1.5",
		OVF:           "http://schemas.dmtf.org/ovf/envelope/1",
		Info:          model.Info,
		NetworkConfig: make([]APIResponseVAppNetworkConfig, 0, len(model.NetworkConfig)),
	}
	for _, network := range model.NetworkConfig {
		section.NetworkConfig = append(section.NetworkConfig, toAPIResponseVAppNetworkConfig(network))
	}
	return section
}

func toAPIResponseVAppNetworkConfig(model types.ModelVAppRuntimeNetworkConfig) APIResponseVAppNetworkConfig {
	return APIResponseVAppNetworkConfig{NetworkID: model.ID, NetworkHREF: model.HREF, Description: model.Description, NetworkType: model.Type, NetworkName: model.NetworkName, Configuration: toAPIResponseVAppRuntimeConfiguration(model.Configuration)}
}

// ToAPIRequestUndeployVApp converts the public undeploy mode to the VMware XML payload.
func ToAPIRequestUndeployVApp(mode types.VAppUndeployMode) APIRequestUndeployVApp {
	switch mode {
	case types.VAppUndeployModeSuspend:
		return APIRequestUndeployVApp{XMLNS: "http://www.vmware.com/vcloud/v1.5", UndeployAction: "suspend"}
	case types.VAppUndeployModeForce:
		return APIRequestUndeployVApp{XMLNS: "http://www.vmware.com/vcloud/v1.5", UndeployAction: "force"}
	default:
		return APIRequestUndeployVApp{XMLNS: "http://www.vmware.com/vcloud/v1.5", UndeployAction: "powerOff"}
	}
}

func toAPIResponseVAppRuntimeConfiguration(model types.ModelVAppRuntimeNetworkConfiguration) APIResponseVAppRuntimeConfiguration {
	result := APIResponseVAppRuntimeConfiguration{
		FenceMode:                      model.FenceMode,
		RetainNetInfoAcrossDeployments: model.RetainNetInfoAcrossDeployments,
		IsDeployed:                     model.IsDeployed,
		GuestVLANAllowed:               model.GuestVLANAllowed,
		IPScopes:                       nil,
	}
	if model.ParentNetwork != nil {
		ref := APIRequestVAppReference{HREF: model.ParentNetwork.HREF, ID: model.ParentNetwork.ID, Type: model.ParentNetwork.Type, Name: model.ParentNetwork.Name}
		result.ParentNetwork = &ref
	}
	if len(model.IPScopes) > 0 {
		result.IPScopes = &APIResponseVAppRuntimeIPScopes{IPScope: make([]APIResponseVAppRuntimeIPScope, 0, len(model.IPScopes))}
		for _, scope := range model.IPScopes {
			result.IPScopes.IPScope = append(result.IPScopes.IPScope, toAPIResponseVAppRuntimeIPScope(scope))
		}
	}
	return result
}

func toAPIResponseVAppRuntimeIPScope(model types.ModelVAppRuntimeIPScope) APIResponseVAppRuntimeIPScope {
	isInherited := model.IsInherited
	result := APIResponseVAppRuntimeIPScope{
		IsInherited:  &isInherited,
		Gateway:      model.Gateway,
		Netmask:      model.Netmask,
		DNS1:         model.DNS1,
		DNS2:         model.DNS2,
		DNSSuffix:    model.DNSSuffix,
		PrefixLength: model.PrefixLength,
		IPRanges:     APIResponseVAppRuntimeIPRanges{IPRange: make([]APIResponseVAppRuntimeIPRange, 0, len(model.IPRanges))},
	}
	for _, ipRange := range model.IPRanges {
		result.IPRanges.IPRange = append(result.IPRanges.IPRange, APIResponseVAppRuntimeIPRange{StartAddress: ipRange.StartAddress, EndAddress: ipRange.EndAddress})
	}
	return result
}

// vAppRecordID returns the vApp URN for a query record. Records can carry a
// null id, so fall back to deriving the URN from the entity HREF
// (/api/vApp/vapp-<uuid>).
func vAppRecordID(id, href string) string {
	if id != "" {
		return id
	}
	if href == "" {
		return ""
	}
	last := href
	if i := strings.LastIndex(href, "/"); i >= 0 {
		last = href[i+1:]
	}
	if uuid := strings.TrimPrefix(last, "vapp-"); uuid != last && uuid != "" {
		return "urn:vcloud:vapp:" + uuid
	}
	return ""
}
