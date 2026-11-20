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
	"maps"
	"sort"
	"strings"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// vAppLeaseSettingsSectionType is the _type of the lease section of the legacy vApp JSON.
const vAppLeaseSettingsSectionType = "LeaseSettingsSectionType"

const (
	vcloudXMLNS      = "http://www.vmware.com/vcloud/v1.5"
	ovfEnvelopeXMLNS = "http://schemas.dmtf.org/ovf/envelope/1"
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
		// The legacy payload may carry only the user href (/api/admin/user/<uuid>):
		// return the user URN to stay consistent with the other identifiers.
		if uuid := urn.ExtractUUID(payload.User.HREF); uuid != "" {
			*o = VAppOwnerRef(urn.Normalize(urn.User, uuid).String())
		} else {
			*o = VAppOwnerRef(payload.User.HREF)
		}
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

	// Links and Sections are only decoded from the real legacy payload: it has
	// no top-level vdc nor leaseSettings, the VDC is the rel="up" link and the
	// lease lives in the LeaseSettingsSectionType section.
	Links    []APIResponseVAppLink    `json:"link" fake:"skip"`
	Sections []APIResponseVAppSection `json:"section" fake:"skip"`
}

// APIResponseVAppLink is one entry of the `link` array of the legacy vApp JSON.
type APIResponseVAppLink struct {
	HREF string `json:"href"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

// APIResponseVAppSection is one entry of the `section` array of the legacy vApp JSON.
// Only the lease section fields are decoded.
type APIResponseVAppSection struct {
	Type                     string `json:"_type"`
	DeploymentLeaseInSeconds *int   `json:"deploymentLeaseInSeconds"`
	StorageLeaseInSeconds    *int   `json:"storageLeaseInSeconds"`
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

// productSectionInfo is the OVF Info of the guest properties section.
const productSectionInfo = "Custom properties"

// APIProductSectionList is the XML ProductSectionList of a vApp (GET and PUT
// of /api/vApp/{id}/productSections). The JSON representation returned by the
// API uses other key names than the SDK models, the XML one is the reference
// (same wire contract as go-vcloud-director).
type APIProductSectionList struct {
	XMLName        xml.Name              `xml:"ProductSectionList"`
	OVF            string                `xml:"xmlns:ovf,attr,omitempty"`
	XMLNS          string                `xml:"xmlns,attr"`
	ProductSection *APIProductSectionXML `xml:"http://schemas.dmtf.org/ovf/envelope/1 ProductSection,omitempty"`
}

// APIProductSectionXML is one OVF ProductSection.
type APIProductSectionXML struct {
	Class    string                  `xml:"http://schemas.dmtf.org/ovf/envelope/1 class,attr,omitempty"`
	Instance string                  `xml:"http://schemas.dmtf.org/ovf/envelope/1 instance,attr,omitempty"`
	Info     string                  `xml:"Info,omitempty"`
	Property []APIProductPropertyXML `xml:"http://schemas.dmtf.org/ovf/envelope/1 Property,omitempty"`
}

// APIProductPropertyXML is one OVF Property of a ProductSection.
type APIProductPropertyXML struct {
	Key              string              `xml:"http://schemas.dmtf.org/ovf/envelope/1 key,attr,omitempty"`
	Label            string              `xml:"http://schemas.dmtf.org/ovf/envelope/1 Label,omitempty"`
	Description      string              `xml:"http://schemas.dmtf.org/ovf/envelope/1 Description,omitempty"`
	DefaultValue     string              `xml:"http://schemas.dmtf.org/ovf/envelope/1 value,attr"`
	Value            *APIProductValueXML `xml:"http://schemas.dmtf.org/ovf/envelope/1 Value,omitempty"`
	Type             string              `xml:"http://schemas.dmtf.org/ovf/envelope/1 type,attr,omitempty"`
	UserConfigurable bool                `xml:"http://schemas.dmtf.org/ovf/envelope/1 userConfigurable,attr"`
}

// APIProductValueXML is the value element of an OVF Property.
type APIProductValueXML struct {
	Value string `xml:"http://schemas.dmtf.org/ovf/envelope/1 value,attr,omitempty"`
}

// ToModel converts the XML ProductSectionList to the public guest property model.
func (r *APIProductSectionList) ToModel() types.ModelVAppGuestProperties {
	model := types.ModelVAppGuestProperties{Sections: make([]types.ModelVAppGuestPropertySection, 0, 1)}
	if r == nil || r.ProductSection == nil {
		return model
	}

	section := types.ModelVAppGuestPropertySection{
		Class:      r.ProductSection.Class,
		Instance:   r.ProductSection.Instance,
		Properties: make(map[string]types.ModelVAppGuestProperty, len(r.ProductSection.Property)),
	}
	for _, property := range r.ProductSection.Property {
		value := property.DefaultValue
		if property.Value != nil {
			value = property.Value.Value
		}
		section.Properties[property.Key] = types.ModelVAppGuestProperty{
			Key:              property.Key,
			Label:            property.Label,
			Type:             property.Type,
			Value:            value,
			UserConfigurable: property.UserConfigurable,
		}
	}
	model.Sections = append(model.Sections, section)

	return model
}

// ToAPIProductSectionList converts the public guest property model to the XML
// ProductSectionList. The properties are sorted by key (deterministic payload).
// Only one ProductSection is sent: the vApp guest properties are one custom section.
func ToAPIProductSectionList(model types.ModelVAppGuestProperties) APIProductSectionList {
	result := APIProductSectionList{OVF: ovfEnvelopeXMLNS, XMLNS: vcloudXMLNS}

	section := APIProductSectionXML{Info: productSectionInfo}
	if len(model.Sections) > 0 {
		section.Class = model.Sections[0].Class
		section.Instance = model.Sections[0].Instance
	}

	merged := make(map[string]types.ModelVAppGuestProperty)
	for _, s := range model.Sections {
		maps.Copy(merged, s.Properties)
	}

	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	section.Property = make([]APIProductPropertyXML, 0, len(keys))
	for _, key := range keys {
		property := merged[key]
		propertyType := property.Type
		if propertyType == "" {
			propertyType = "string"
		}
		section.Property = append(section.Property, APIProductPropertyXML{
			Key:              key,
			Label:            cmpOr(property.Label, key),
			Type:             propertyType,
			Value:            &APIProductValueXML{Value: property.Value},
			UserConfigurable: property.UserConfigurable,
		})
	}

	result.ProductSection = &section
	return result
}

func cmpOr(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
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
	NetworkID   string `xml:"id,attr,omitempty"`
	NetworkHREF string `xml:"href,attr,omitempty"`
	// Link is the first child element of the network config. The id and href
	// attributes are not always present on the real payload, the Link href is then
	// the only place with the network UUID (go-vcloud-director reads it the same way).
	Link          *APIResponseVAppNetworkLink         `xml:"Link,omitempty"`
	Description   string                              `xml:"Description,omitempty"`
	NetworkType   string                              `xml:"type,attr,omitempty"`
	NetworkName   string                              `xml:"networkName,attr,omitempty"`
	Configuration APIResponseVAppRuntimeConfiguration `xml:"Configuration"`
	// IsDeployed is a child of NetworkConfig (after Configuration), not of Configuration.
	IsDeployed *bool `xml:"IsDeployed,omitempty"`
}

// APIResponseVAppNetworkLink is the Link child of a vApp NetworkConfig.
type APIResponseVAppNetworkLink struct {
	Rel  string `xml:"rel,attr,omitempty"`
	HREF string `xml:"href,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

type APIResponseVAppRuntimeConfiguration struct {
	// Field order mirrors the XSD sequence of NetworkConfigurationType (VMware
	// Cloud Director API, NetworkConfigSection): BackwardCompatibilityMode,
	// IpScopes, ParentNetwork, FenceMode, RetainNetInfoAcrossDeployments,
	// Features, SyslogServerSettings, RouterInfo, SubInterface,
	// DistributedInterface, ServiceInterface, GuestVlanAllowed, Connected,
	// DualStackNetwork. Out-of-order elements are rejected with HTTP 400
	// cvc-complex-type.2.4.a. IsDeployed is NOT part of it: it is a child of
	// NetworkConfig.
	//
	// The elements that are not exposed by the public model are kept here only so
	// that the read-modify-write of the whole section does not drop them.
	BackwardCompatibilityMode      *bool                             `xml:"BackwardCompatibilityMode,omitempty"`
	IPScopes                       *APIResponseVAppRuntimeIPScopes   `xml:"IpScopes,omitempty"`
	ParentNetwork                  *APIRequestVAppReference          `xml:"ParentNetwork,omitempty"`
	FenceMode                      string                            `xml:"FenceMode,omitempty"`
	RetainNetInfoAcrossDeployments *bool                             `xml:"RetainNetInfoAcrossDeployments,omitempty"`
	Features                       *APIResponseVAppRuntimeFeatures   `xml:"Features,omitempty"`
	SyslogServerSettings           *APIResponseVAppRuntimeInnerXML   `xml:"SyslogServerSettings,omitempty"`
	RouterInfo                     *APIResponseVAppRuntimeRouterInfo `xml:"RouterInfo,omitempty"`
	SubInterface                   *bool                             `xml:"SubInterface,omitempty"`
	DistributedInterface           *bool                             `xml:"DistributedInterface,omitempty"`
	ServiceInterface               *bool                             `xml:"ServiceInterface,omitempty"`
	GuestVLANAllowed               *bool                             `xml:"GuestVlanAllowed,omitempty"`
	Connected                      *bool                             `xml:"Connected,omitempty"`
	DualStackNetwork               *bool                             `xml:"DualStackNetwork,omitempty"`
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

// APIRequestRecomposeVApp is the XML body of the recomposeVApp action used to
// update the name and the description of a vApp.
type APIRequestRecomposeVApp struct {
	XMLName     xml.Name `xml:"RecomposeVAppParams"`
	OVF         string   `xml:"xmlns:ovf,attr"`
	XSI         string   `xml:"xmlns:xsi,attr"`
	XMLNS       string   `xml:"xmlns,attr"`
	Name        string   `xml:"name,attr"`
	Deploy      bool     `xml:"deploy,attr"`
	Description string   `xml:"Description"`
}

// APIRequestLeaseSettingsSection is the XML body replacing the lease settings
// section of a vApp.
type APIRequestLeaseSettingsSection struct {
	XMLName                  xml.Name `xml:"LeaseSettingsSection"`
	OVF                      string   `xml:"xmlns:ovf,attr"`
	XMLNS                    string   `xml:"xmlns,attr"`
	Info                     string   `xml:"ovf:Info"`
	DeploymentLeaseInSeconds *int     `xml:"DeploymentLeaseInSeconds,omitempty"`
	StorageLeaseInSeconds    *int     `xml:"StorageLeaseInSeconds,omitempty"`
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

	// The real payload carries the lease in the LeaseSettingsSectionType section.
	// A lease of 0 means "never expires" and must not be confused with "absent".
	for _, section := range r.Sections {
		if section.Type != vAppLeaseSettingsSectionType {
			continue
		}
		if section.DeploymentLeaseInSeconds != nil {
			v := *section.DeploymentLeaseInSeconds
			m.DeploymentLeaseInSeconds = &v
		}
		if section.StorageLeaseInSeconds != nil {
			v := *section.StorageLeaseInSeconds
			m.StorageLeaseInSeconds = &v
		}
		break
	}

	if m.VDC == "" {
		m.VDC = vdcURNFromLinks(r.Links)
	}

	for _, meta := range r.Metadata {
		if meta.Key == "vappBillingModel" {
			m.Properties.BillingModel = meta.Value
		}
	}
	m.GuestProperties = r.GuestProperties.ToModel()

	return m
}

// vdcURNFromLinks returns the URN of the parent VDC from the rel="up" link
// (https://<console>/api/vdc/<uuid>), or an empty string.
func vdcURNFromLinks(links []APIResponseVAppLink) string {
	for _, link := range links {
		if link.Rel != "up" || !strings.Contains(link.HREF, "/api/vdc/") {
			continue
		}
		if uuid := urn.ExtractUUID(link.HREF); uuid != "" {
			return urn.Normalize(urn.VDC, uuid).String()
		}
	}
	return ""
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
	href := r.NetworkHREF
	if href == "" && r.Link != nil {
		href = r.Link.HREF
	}
	configuration := r.Configuration.ToModel()
	// IsDeployed is a child of NetworkConfig; the public model keeps it in the configuration.
	configuration.IsDeployed = r.IsDeployed

	return types.ModelVAppRuntimeNetworkConfig{
		Description:   r.Description,
		ID:            r.NetworkID,
		HREF:          href,
		Type:          r.NetworkType,
		NetworkName:   r.NetworkName,
		Configuration: configuration,
	}
}

// ToModel converts APIResponseVAppRuntimeConfiguration to the public runtime model.
func (r *APIResponseVAppRuntimeConfiguration) ToModel() types.ModelVAppRuntimeNetworkConfiguration {
	model := types.ModelVAppRuntimeNetworkConfiguration{
		FenceMode:                      r.FenceMode,
		RetainNetInfoAcrossDeployments: r.RetainNetInfoAcrossDeployments,
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
// The round trip with ToModel is not symmetric: the href of the Link child of a NetworkConfig is
// exposed as HREF by ToModel and written back as the href attribute of the NetworkConfig, the
// Link element itself is not rebuilt. Updates that must keep the Link patch the raw section
// instead (patchRawVAppRuntimeNetwork).
func ToAPIResponseVAppNetworkConfigSection(model *types.ModelVAppRuntimeNetworkConfigSection) *APIResponseVAppNetworkConfigSection {
	if model == nil {
		return nil
	}
	section := &APIResponseVAppNetworkConfigSection{
		XMLNS:         vcloudXMLNS,
		OVF:           ovfEnvelopeXMLNS,
		Info:          model.Info,
		NetworkConfig: make([]APIResponseVAppNetworkConfig, 0, len(model.NetworkConfig)),
	}
	for _, network := range model.NetworkConfig {
		section.NetworkConfig = append(section.NetworkConfig, toAPIResponseVAppNetworkConfig(network))
	}
	return section
}

func toAPIResponseVAppNetworkConfig(model types.ModelVAppRuntimeNetworkConfig) APIResponseVAppNetworkConfig {
	return APIResponseVAppNetworkConfig{NetworkID: model.ID, NetworkHREF: model.HREF, Description: model.Description, NetworkType: model.Type, NetworkName: model.NetworkName, Configuration: toAPIResponseVAppRuntimeConfiguration(model.Configuration), IsDeployed: model.Configuration.IsDeployed}
}

// ToAPIRequestUndeployVApp converts the public undeploy mode to the VMware XML payload.
func ToAPIRequestUndeployVApp(mode types.VAppUndeployMode) APIRequestUndeployVApp {
	switch mode {
	case types.VAppUndeployModeSuspend:
		return APIRequestUndeployVApp{XMLNS: vcloudXMLNS, UndeployAction: "suspend"}
	case types.VAppUndeployModeForce:
		return APIRequestUndeployVApp{XMLNS: vcloudXMLNS, UndeployAction: "force"}
	default:
		return APIRequestUndeployVApp{XMLNS: vcloudXMLNS, UndeployAction: "powerOff"}
	}
}

func toAPIResponseVAppRuntimeConfiguration(model types.ModelVAppRuntimeNetworkConfiguration) APIResponseVAppRuntimeConfiguration {
	result := APIResponseVAppRuntimeConfiguration{
		FenceMode:                      model.FenceMode,
		RetainNetInfoAcrossDeployments: model.RetainNetInfoAcrossDeployments,
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
