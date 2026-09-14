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
	ApiResponseListVdcNetwork struct {
		Values []ApiResponseVdcNetwork `json:"values" fakesize:"3"`
	}

	// * Get / Create / Update / shared model
	ApiResponseVdcNetwork struct {
		ID                      string                   `json:"id,omitempty" fake:"{urn:network}"`
		Name                    string                   `json:"name" fake:"mockvdcnetwork-{word}"`
		Description             string                   `json:"description,omitempty" fake:"{sentence}"`
		Status                  string                   `json:"status,omitempty"`
		OwnerRef                *ApiObjectReference      `json:"ownerRef,omitempty"`
		NetworkType             string                   `json:"networkType"`
		Connection              *ApiVdcNetworkConnection `json:"connection,omitempty"`
		GuestVlanTaggingAllowed *bool                    `json:"guestVlanTaggingAllowed"`
		Subnets                 ApiVdcNetworkSubnets     `json:"subnets"`
		Shared                  *bool                    `json:"shared,omitempty"`
	}

	ApiVdcNetworkConnection struct {
		RouterRef           ApiObjectReference `json:"routerRef"`
		ConnectionTypeValue string             `json:"connectionTypeValue,omitempty"`
	}

	ApiVdcNetworkSubnets struct {
		Values []ApiVdcNetworkSubnetValue `json:"values"`
	}

	ApiVdcNetworkSubnetValue struct {
		Gateway      string                `json:"gateway"`
		PrefixLength int                   `json:"prefixLength"`
		DNSServer1   string                `json:"dnsServer1,omitempty"`
		DNSServer2   string                `json:"dnsServer2,omitempty"`
		DNSSuffix    string                `json:"dnsSuffix,omitempty"`
		IPRanges     ApiVdcNetworkIPRanges `json:"ipRanges"`
	}

	ApiVdcNetworkIPRanges struct {
		Values []ApiVdcNetworkIPRangeValue `json:"values"`
	}

	ApiVdcNetworkIPRangeValue struct {
		StartAddress string `json:"startAddress"`
		EndAddress   string `json:"endAddress"`
	}

	// * Create / Update request (same shape as ApiResponseVdcNetwork)
	ApiRequestVdcNetwork struct {
		ID                      string                   `json:"id,omitempty" fake:"{urn:network}"`
		Name                    string                   `json:"name" fake:"mockvdcnetwork-{word}"`
		Description             string                   `json:"description,omitempty" fake:"{sentence}"`
		Status                  string                   `json:"status,omitempty"`
		OwnerRef                *ApiObjectReference      `json:"ownerRef,omitempty"`
		NetworkType             string                   `json:"networkType"`
		Connection              *ApiVdcNetworkConnection `json:"connection,omitempty"`
		GuestVlanTaggingAllowed *bool                    `json:"guestVlanTaggingAllowed"`
		Subnets                 ApiVdcNetworkSubnets     `json:"subnets"`
		Shared                  *bool                    `json:"shared,omitempty"`
	}
)

func (r *ApiResponseListVdcNetwork) ToModel() *types.ModelListVdcNetwork {
	model := &types.ModelListVdcNetwork{
		VdcNetworks: make([]types.ModelGetVdcNetwork, 0),
	}

	for _, network := range r.Values {
		model.VdcNetworks = append(model.VdcNetworks, network.ToModel())
	}

	return model
}

func (r *ApiResponseVdcNetwork) ToModel() types.ModelGetVdcNetwork {
	m := types.ModelGetVdcNetwork{
		ID:                      r.ID,
		Name:                    r.Name,
		Description:             r.Description,
		Status:                  r.Status,
		NetworkType:             r.NetworkType,
		GuestVlanTaggingAllowed: r.GuestVlanTaggingAllowed,
		Shared:                  r.Shared,
	}

	if r.OwnerRef != nil {
		m.OwnerID = r.OwnerRef.ID
		m.OwnerName = r.OwnerRef.Name
	}

	if len(r.Subnets.Values) > 0 {
		subnet := r.Subnets.Values[0]
		m.Subnet = types.ModelVdcNetworkSubnet{
			Gateway:      subnet.Gateway,
			PrefixLength: subnet.PrefixLength,
			DNSServer1:   subnet.DNSServer1,
			DNSServer2:   subnet.DNSServer2,
			DNSSuffix:    subnet.DNSSuffix,
		}
		for _, ipRange := range subnet.IPRanges.Values {
			m.Subnet.IPRanges = append(m.Subnet.IPRanges, types.ModelVdcNetworkIPRange{
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

func (r *ApiRequestVdcNetwork) ToModel() types.ModelGetVdcNetwork {
	resp := ApiResponseVdcNetwork{
		ID:                      r.ID,
		Name:                    r.Name,
		Description:             r.Description,
		Status:                  r.Status,
		OwnerRef:                r.OwnerRef,
		NetworkType:             r.NetworkType,
		Connection:              r.Connection,
		GuestVlanTaggingAllowed: r.GuestVlanTaggingAllowed,
		Subnets:                 r.Subnets,
		Shared:                  r.Shared,
	}

	return resp.ToModel()
}
