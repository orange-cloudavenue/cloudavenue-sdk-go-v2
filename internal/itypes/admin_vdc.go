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

// AdminVDC represents a VMware vCD AdminVDC XML element.
type AdminVDC struct {
	XMLName               xml.Name    `xml:"Vdc"`
	Name                  string      `xml:"name,attr"`
	ID                    string      `xml:"id,attr"`
	Href                  string      `xml:"href,attr"`
	DisplayName           string      `xml:"DisplayName,omitempty"`
	Description           string      `xml:"Description,omitempty"`
	IsEnabled             bool        `xml:"IsEnabled,omitempty"`
	IsFullProtected       bool        `xml:"IsFullProtected,omitempty"`
	VdcStorageProfiles    []Reference `xml:"StorageProfiles>StorageProfile"`
	ComputePolicies       []Reference `xml:"ComputePolicies>ComputePolicy"`
	Networks              []Reference `xml:"Networks>Network"`
	EdgeGatewayReferences []Reference `xml:"EdgeGateways>EdgeGateway"`
	VmCount               int         `xml:"VmCount"`
	VmRunningCount        int         `xml:"VmRunningCount"`
	VappCount             int         `xml:"VappCount"`
	VappTemplateCount     int         `xml:"VappTemplateCount"`
	Allocation            struct {
		CPU     int `xml:"CPU"`
		Memory  int `xml:"Memory"`
		Storage int `xml:"Storage"`
	} `xml:"Allocation"`
	ComputePolicyReferences []Reference `xml:"ComputePolicyReferences>ComputePolicyReference"`
}

// AdminVDCs represents the wrapper for a list of VDCs in VMware vCD AdminVDC API.
type AdminVDCs struct {
	XMLName xml.Name   `xml:"Vdcs"`
	VDCs    []AdminVDC `xml:"Vdc"`
}

// ToModel converts the VMware vCD AdminVDC XML response to the public model.
func (r *AdminVDC) ToModel() *types.ModelAdminVDC {
	return &types.ModelAdminVDC{
		ID:                r.ID,
		Name:              r.Name,
		Href:              r.Href,
		DisplayName:       r.DisplayName,
		Description:       r.Description,
		IsEnabled:         r.IsEnabled,
		IsFullProtected:   r.IsFullProtected,
		VmCount:           r.VmCount,
		VmRunningCount:    r.VmRunningCount,
		VappCount:         r.VappCount,
		VappTemplateCount: r.VappTemplateCount,
	}
}
