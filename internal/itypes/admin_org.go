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

// AdminOrg represents a VMware vCD AdminOrg XML element.
type AdminOrg struct {
	XMLName         xml.Name     `xml:"Org"`
	Name            string       `xml:"name,attr"`
	ID              string       `xml:"id,attr"`
	Href            string       `xml:"href,attr"`
	DisplayName     string       `xml:"DisplayName,omitempty"`
	Description     string       `xml:"Description,omitempty"`
	IsEnabled       bool         `xml:"IsEnabled,omitempty"`
	IsFullProtected bool         `xml:"IsFullProtected,omitempty"`
	OrgSettings     *OrgSettings `xml:"OrgSettings,omitempty"`
	Vdcs            []Reference  `xml:"Vdcs>Vdc"`
	Networks        []Reference  `xml:"Networks>Network"`
	Catalogs        []Reference  `xml:"Catalogs>Catalog"`
	StorageProfiles []Reference  `xml:"StorageProfiles>StorageProfile"`
	VDCGroups       []Reference  `xml:"VdcGroups>VdcGroup"`
	Rights          []Reference  `xml:"Rights>Right"`
	Users           []Reference  `xml:"Users>User"`
	Groups          []Reference  `xml:"Groups>Group"`
	Roles           []Reference  `xml:"Roles>Role"`
}

// OrgSettings represents the OrgSettings element of an AdminOrg.
type OrgSettings struct {
	XMLName xml.Name `xml:"OrgSettings"`
}

// AdminOrgs represents the wrapper for a list of organizations in VMware vCD AdminOrg API.
type AdminOrgs struct {
	XMLName xml.Name   `xml:"Orgs"`
	Orgs    []AdminOrg `xml:"Org"`
}

// ToModel converts the VMware vCD AdminOrg XML response to the public model.
func (r *AdminOrg) ToModel() *types.ModelAdminOrg {
	return &types.ModelAdminOrg{
		ID:              r.ID,
		Name:            r.Name,
		Href:            r.Href,
		DisplayName:     r.DisplayName,
		Description:     r.Description,
		IsEnabled:       r.IsEnabled,
		IsFullProtected: r.IsFullProtected,
	}
}
