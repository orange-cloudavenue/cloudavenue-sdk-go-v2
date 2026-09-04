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

type (
	// * List
	APIResponseListVDCNetwork struct {
		Values []APIResponseVDCNetwork `json:"values" fakesize:"3"`
	}

	// * Get / Create / Update / shared model
	APIResponseVDCNetwork struct {
		ID                      string                   `json:"id,omitempty" fake:"{urn:network}"`
		Name                    string                   `json:"name" fake:"mockvdcnetwork-{word}"`
		Description             string                   `json:"description,omitempty" fake:"{sentence}"`
		Status                  string                   `json:"status,omitempty"`
		OwnerRef                *APIObjectReference      `json:"ownerRef,omitempty"`
		NetworkType             string                   `json:"networkType"`
		Connection              *APIVDCNetworkConnection `json:"connection,omitempty"`
		GuestVLANTaggingAllowed *bool                    `json:"guestVLANTaggingAllowed"`
		Subnets                 APIVDCNetworkSubnets     `json:"subnets"`
		Shared                  *bool                    `json:"shared,omitempty"`
	}

	APIVDCNetworkConnection struct {
		RouterRef           APIObjectReference `json:"routerRef"`
		ConnectionTypeValue string             `json:"connectionTypeValue,omitempty"`
	}

	APIVDCNetworkSubnets struct {
		Values []APIVDCNetworkSubnetValue `json:"values"`
	}

	APIVDCNetworkSubnetValue struct {
		Gateway      string                `json:"gateway"`
		PrefixLength int                   `json:"prefixLength"`
		DNSServer1   string                `json:"dnsServer1,omitempty"`
		DNSServer2   string                `json:"dnsServer2,omitempty"`
		DNSSuffix    string                `json:"dnsSuffix,omitempty"`
		IPRanges     APIVDCNetworkIPRanges `json:"ipRanges"`
	}

	APIVDCNetworkIPRanges struct {
		Values []APIVDCNetworkIPRangeValue `json:"values"`
	}

	APIVDCNetworkIPRangeValue struct {
		StartAddress string `json:"startAddress"`
		EndAddress   string `json:"endAddress"`
	}

	// * Create / Update request (same shape as APIResponseVDCNetwork)
	APIRequestVDCNetwork struct {
		ID                      string                   `json:"id,omitempty" fake:"{urn:network}"`
		Name                    string                   `json:"name" fake:"mockvdcnetwork-{word}"`
		Description             string                   `json:"description,omitempty" fake:"{sentence}"`
		Status                  string                   `json:"status,omitempty"`
		OwnerRef                *APIObjectReference      `json:"ownerRef,omitempty"`
		NetworkType             string                   `json:"networkType"`
		Connection              *APIVDCNetworkConnection `json:"connection,omitempty"`
		GuestVLANTaggingAllowed *bool                    `json:"guestVLANTaggingAllowed"`
		Subnets                 APIVDCNetworkSubnets     `json:"subnets"`
		Shared                  *bool                    `json:"shared,omitempty"`
	}
)

func (r *APIResponseListVDCNetwork) ToModel() *types.ModelListVDCNetwork {
	model := &types.ModelListVDCNetwork{
		VDCNetworks: make([]types.ModelGetVDCNetwork, 0),
	}

	for _, network := range r.Values {
		model.VDCNetworks = append(model.VDCNetworks, network.ToModel())
	}

	return model
}

func (r *APIResponseVDCNetwork) ToModel() types.ModelGetVDCNetwork {
	m := types.ModelGetVDCNetwork{
		ID:                      r.ID,
		Name:                    r.Name,
		Description:             r.Description,
		Status:                  r.Status,
		NetworkType:             r.NetworkType,
		GuestVLANTaggingAllowed: r.GuestVLANTaggingAllowed,
		Shared:                  r.Shared,
	}

	if r.OwnerRef != nil {
		m.OwnerID = r.OwnerRef.ID
		m.OwnerName = r.OwnerRef.Name
	}

	if len(r.Subnets.Values) > 0 {
		subnet := r.Subnets.Values[0]
		m.Subnet = types.ModelVDCNetworkSubnet{
			Gateway:      subnet.Gateway,
			PrefixLength: subnet.PrefixLength,
			DNSServer1:   subnet.DNSServer1,
			DNSServer2:   subnet.DNSServer2,
			DNSSuffix:    subnet.DNSSuffix,
		}
		for _, ipRange := range subnet.IPRanges.Values {
			m.Subnet.IPRanges = append(m.Subnet.IPRanges, types.ModelVDCNetworkIPRange{
				StartAddress: ipRange.StartAddress,
				EndAddress:   ipRange.EndAddress,
			})
		}
	}

	if r.Connection != nil {
		m.EdgeGatewayID = r.Connection.RouterRef.ID
		m.EdgeGatewayName = r.Connection.RouterRef.Name
	}

	return m
}

func (r *APIRequestVDCNetwork) ToModel() types.ModelGetVDCNetwork {
	resp := APIResponseVDCNetwork{
		ID:                      r.ID,
		Name:                    r.Name,
		Description:             r.Description,
		Status:                  r.Status,
		OwnerRef:                r.OwnerRef,
		NetworkType:             r.NetworkType,
		Connection:              r.Connection,
		GuestVLANTaggingAllowed: r.GuestVLANTaggingAllowed,
		Subnets:                 r.Subnets,
		Shared:                  r.Shared,
	}

	return resp.ToModel()
}
