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

// * Request / Response API

type (
	ApiRequestEdgeGateway struct {
		T0Name string `json:"tier0VrfId" fake:"{resource_name:t0}" validate:"required,resource_name=t0"`
	}

	ApiResponseEdgegateways struct {
		Values []ApiResponseEdgegateway `json:"values,omitempty" fakesize:"1"` // List of edge gateways.
	}
	ApiResponseEdgegateway struct {
		ID          string `json:"id" fake:"{urn:edgegateway}"`             // The ID of the edge gateway.
		Name        string `json:"name" fake:"{resource_name:edgegateway}"` // The name of the edge gateway.
		Description string `json:"description" fake:"{sentence}"`

		EdgeGatewayUplinks []struct {
			Connected bool `json:"connected" fake:"true"` // Indicates if the uplink is connected.
			Dedicated bool `json:"dedicated"`
			Subnets   struct {
				Values []struct {
					DNSServer1 string `json:"dnsServer1" fake:"{ipv4address}"`
					DNSServer2 string `json:"dnsServer2" fake:"{ipv4address}"`
					DNSSuffix  string `json:"dnsSuffix" fake:"{domainname}"`
					Enabled    bool   `json:"enabled" fake:"{bool}"` // Indicates if the subnet is enabled.
					Gateway    string `json:"gateway" fake:"{ipv4address}"`
					IPRanges   struct {
						Values []struct {
							EndAddress   string `json:"endAddress" fake:"{ipv4address}"`
							StartAddress string `json:"startAddress" fake:"{ipv4address}"`
						} `json:"values" fakesize:"1"`
					} `json:"ipRanges"`
					PrefixLength int64  `json:"prefixLength" fake:"{number:24,32}"` // The prefix length of the subnet.
					PrimaryIP    string `json:"primaryIp" fake:"{ipv4address}"`
					TotalIPCount int64  `json:"totalIpCount" fake:"{number:5,10}"` // The total number of IP addresses in the subnet.
					UsedIPCount  int64  `json:"usedIpCount" fake:"{number:6,8}"`   // The number of used IP addresses in the subnet.
				} `json:"values" fakesize:"1"`
			} `json:"subnets" fakesize:"1"`
			UplinkID   string `json:"uplinkId" fake:"{urn:network}"`
			UplinkName string `json:"uplinkName" fake:"{resource_name:t0}"` // The name of the uplink.
		} `json:"edgeGatewayUplinks" fakesize:"1"`

		OrgVDC *ApiObjectReference `json:"orgVdc"`

		// OwnerRef contains information about the owner of the edge gateway (VDC Or VDCGroup)
		OwnerRef *ApiObjectReference `json:"ownerRef"`

		// OrgVdcNetworkCount holds the number of Org VDC networks connected to the gateway.
		OrgVDCNetworkCount int64 `json:"orgVdcNetworkCount" fake:"{number:1,10}"`
	}

	ApiResponseQueryEdgeGateway struct {
		Record []ApiResponseQueryEdgeGatewayRecord `json:"record,omitempty" fakesize:"1"` // List of edge gateways.
	}

	ApiResponseQueryEdgeGatewayRecord struct {
		ID                  string `json:"-"`                                                 // The ID of the entity.
		HREF                string `json:"href,omitempty" fake:"{href_uuid}"`                 // The URI of the entity.
		Type                string `json:"type,omitempty"`                                    // The MIME type of the entity.
		Name                string `json:"name,omitempty" fake:"{resource_name:edgegateway}"` // EdgeGateway name.
		VDCID               string `json:"vdc,omitempty" fake:"{urn:vdc}"`                    // VDC Reference or ID
		VDCName             string `json:"orgVdcName,omitempty" fake:"{word}"`                // VDC name
		NumberOfExtNetworks int    `json:"numberOfExtNetworks,omitempty" fake:"{number:1,5}"` // Number of external networks connected to the edgeGateway.
		NumberOfOrgNetworks int    `json:"numberOfOrgNetworks,omitempty" fake:"{number:1,5}"` // Number of org VDC networks connected to the edgeGateway
		IsBusy              bool   `json:"isBusy,omitempty" fake:"{bool}"`                    // True if this Edge Gateway is busy.
		GatewayStatus       string `json:"gatewayStatus,omitempty" fake:"{word}"`             // Status of the edgeGateway
	}
)

// ToModel converts the ApiResponseEdgegateways to ModelEdgeGateways.
func (api *ApiResponseEdgegateways) ToModel() *types.ModelEdgeGateways {
	if api == nil {
		return nil
	}

	model := &types.ModelEdgeGateways{
		EdgeGateways: make([]types.ModelEdgeGateway, 0, len(api.Values)),
	}

	for _, v := range api.Values {
		model.EdgeGateways = append(model.EdgeGateways, *v.ToModel())
	}

	return model
}

// ToModel converts the edgegatewayAPIResponse to ModelEdgeGateway.
func (api *ApiResponseEdgegateway) ToModel() *types.ModelEdgeGateway {
	if api == nil {
		return nil
	}

	return &types.ModelEdgeGateway{
		ID:          api.ID,
		Name:        api.Name,
		Description: api.Description,
		OwnerRef: func() *types.ModelObjectReference {
			if api.OwnerRef != nil {
				return &types.ModelObjectReference{
					ID:   api.OwnerRef.ID,
					Name: api.OwnerRef.Name,
				}
			}
			return nil
		}(),
		UplinkT0: func() *types.ModelObjectReference {
			if len(api.EdgeGatewayUplinks) > 0 {
				return &types.ModelObjectReference{
					ID:   api.EdgeGatewayUplinks[0].UplinkID,
					Name: api.EdgeGatewayUplinks[0].UplinkName,
				}
			}
			return nil
		}(),
	}
}
