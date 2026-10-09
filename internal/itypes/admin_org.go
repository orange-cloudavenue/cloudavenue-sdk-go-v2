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
	"encoding/xml"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// AdminOrg represents a read-only VMware vCD AdminOrg XML response element.
type AdminOrg struct {
	// XMLName deliberately has no element constraint. VMware has returned both
	// AdminOrg and Org roots for this admin representation across API versions.
	XMLName           xml.Name
	Xmlns             string       `xml:"xmlns,attr,omitempty"`
	Type              string       `xml:"type,attr,omitempty"`
	Name              string       `xml:"name,attr"`
	ID                string       `xml:"id,attr"`
	Href              string       `xml:"href,attr"`
	DisplayName       string       `xml:"FullName,omitempty"`
	LegacyDisplayName string       `xml:"DisplayName,omitempty"`
	Description       string       `xml:"Description,omitempty"`
	IsEnabled         bool         `xml:"IsEnabled,omitempty"`
	IsFullProtected   bool         `xml:"IsFullProtected,omitempty"`
	OrgSettings       *OrgSettings `xml:"Settings,omitempty"`
	Vdcs              []Reference  `xml:"Vdcs>Vdc"`
	Networks          []Reference  `xml:"Networks>Network"`
	Catalogs          []Reference  `xml:"Catalogs>CatalogReference"`
	StorageProfiles   []Reference  `xml:"StorageProfiles>VdcStorageProfile"`
	VDCGroups         []Reference  `xml:"VdcGroups>VdcGroup"`
	Rights            []Reference  `xml:"RightReferences>RightReference"`
	Users             []Reference  `xml:"Users>UserReference"`
	Groups            []Reference  `xml:"Groups>GroupReference"`
	Roles             []Reference  `xml:"RoleReferences>RoleReference"`
}

// OrgSettings represents the OrgSettings element of an AdminOrg.
type OrgSettings struct {
	Href                      string                `xml:"href,attr,omitempty"`
	Type                      string                `xml:"type,attr,omitempty"`
	OrgGeneralSettings        *OrgGeneralSettings   `xml:"OrgGeneralSettings,omitempty"`
	OrgVAppLeaseSettings      *VAppLeaseSettings    `xml:"VAppLeaseSettings,omitempty"`
	OrgVAppTemplateSettings   *VAppTemplateSettings `xml:"VAppTemplateLeaseSettings,omitempty"`
	OrgLDAPSettings           *OrgLDAPSettings      `xml:"OrgLdapSettings,omitempty"`
	OrgPasswordPolicySettings *OrgPasswordPolicy    `xml:"OrgPasswordPolicySettings,omitempty"`
}

// OrgGeneralSettings contains organization-wide user and catalog settings.
type OrgGeneralSettings struct {
	CanPublishCatalogs       bool `xml:"CanPublishCatalogs,omitempty"`
	CanPublishExternally     bool `xml:"CanPublishExternally,omitempty"`
	CanSubscribe             bool `xml:"CanSubscribe,omitempty"`
	DeployedVMQuota          int  `xml:"DeployedVMQuota,omitempty"`
	StoredVMQuota            int  `xml:"StoredVmQuota,omitempty"`
	UseServerBootSequence    bool `xml:"UseServerBootSequence,omitempty"`
	DelayAfterPowerOnSeconds int  `xml:"DelayAfterPowerOnSeconds,omitempty"`
}

// VAppLeaseSettings contains organization vApp lease settings.
type VAppLeaseSettings struct {
	DeleteOnStorageLeaseExpiration   *bool `xml:"DeleteOnStorageLeaseExpiration,omitempty"`
	DeploymentLeaseSeconds           *int  `xml:"DeploymentLeaseSeconds,omitempty"`
	StorageLeaseSeconds              *int  `xml:"StorageLeaseSeconds,omitempty"`
	PowerOffOnRuntimeLeaseExpiration *bool `xml:"PowerOffOnRuntimeLeaseExpiration,omitempty"`
}

// VAppTemplateSettings contains organization vApp template lease settings.
type VAppTemplateSettings struct {
	DeleteOnStorageLeaseExpiration *bool `xml:"DeleteOnStorageLeaseExpiration,omitempty"`
	StorageLeaseSeconds            *int  `xml:"StorageLeaseSeconds,omitempty"`
}

// OrgLDAPSettings keeps the LDAP settings element available without
// pretending to model its provider-specific nested configuration.
type OrgLDAPSettings struct {
	Mode string `xml:"OrgLdapMode,omitempty"`
}

// OrgPasswordPolicy is the password-policy settings element.
type OrgPasswordPolicy struct{}

// AdminOrgs represents a VMware vCD AdminOrg XML response containing
// read-only organization views.
type AdminOrgs struct {
	XMLName xml.Name
	Orgs    []AdminOrg `xml:"Org"`
}

// ToModel converts the VMware vCD AdminOrg XML response to the public model.
func (r *AdminOrg) ToModel() *types.ModelAdminOrg {
	displayName := r.DisplayName
	if displayName == "" {
		displayName = r.LegacyDisplayName
	}
	if displayName == "" {
		displayName = r.Name
	}
	return &types.ModelAdminOrg{
		ID:              r.ID,
		Name:            r.Name,
		Href:            r.Href,
		DisplayName:     displayName,
		Description:     r.Description,
		IsEnabled:       r.IsEnabled,
		IsFullProtected: r.IsFullProtected,
	}
}
