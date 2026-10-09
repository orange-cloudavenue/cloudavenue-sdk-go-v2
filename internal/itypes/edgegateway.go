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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// * Request / Response API

type (
	APIRequestEdgeGateway struct {
		T0Name string `json:"tier0VrfId" fake:"{resource_name:t0}" validate:"required,resource_name=t0"`
	}

	APIResponseEdgegateways struct {
		Values []APIResponseEdgegateway `json:"values,omitempty" fakesize:"1"` // List of edge gateways.
	}
	APIResponseEdgegateway struct {
		ID          string `json:"id" fake:"{urn:edgegateway}"`             // The ID of the edge gateway.
		Name        string `json:"name" fake:"{resource_name:edgegateway}"` // The name of the edge gateway.
		Description string `json:"description" fake:"{sentence}"`
		EdgeID      string `json:"edgeId,omitempty"`
		EdgeName    string `json:"edgeName,omitempty"`
		OwnerType   string `json:"ownerType,omitempty"`
		OwnerName   string `json:"ownerName,omitempty"`
		Tier0VRFID  string `json:"tier0VrfId,omitempty"`

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

		OrgVDC *APIObjectReference `json:"orgVdc"`

		// OwnerRef contains information about the owner of the edge gateway (VDC Or VDCGroup)
		OwnerRef *APIObjectReference `json:"ownerRef"`

		// OrgVdcNetworkCount holds the number of Org VDC networks connected to the gateway.
		OrgVDCNetworkCount int64 `json:"orgVdcNetworkCount" fake:"{number:1,10}"`
	}

	APIResponseQueryEdgeGateway struct {
		Record []APIResponseQueryEdgeGatewayRecord `json:"record,omitempty" fakesize:"1"` // List of edge gateways.
	}

	APIResponseQueryEdgeGatewayRecord struct {
		ID                  string `json:"id,omitempty" fake:"{urn:edgegateway}"`             // The ID of the entity.
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

// UnmarshalJSON accepts Cerberus's array response and the former CloudAPI
// values envelope.
func (r *APIResponseEdgegateways) UnmarshalJSON(data []byte) error {
	var values []APIResponseEdgegateway
	if err := json.Unmarshal(data, &values); err == nil {
		r.Values = values
		return nil
	}
	var envelope struct {
		Values []APIResponseEdgegateway            `json:"values"`
		Record []APIResponseQueryEdgeGatewayRecord `json:"record"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	r.Values = envelope.Values
	if len(r.Values) == 0 && len(envelope.Record) > 0 {
		r.Values = make([]APIResponseEdgegateway, 0, len(envelope.Record))
		for _, record := range envelope.Record {
			r.Values = append(r.Values, APIResponseEdgegateway{
				ID:     record.ID,
				Name:   record.Name,
				OrgVDC: &APIObjectReference{ID: record.VDCID, Name: record.VDCName},
			})
		}
	}
	return nil
}

// UnmarshalJSON accepts both the former query response (record) and the
// CloudAPI edge gateway collection (values). QueryEdgeGateway remains a
// compatibility alias, so older fixtures and integrations continue to decode.
func (r *APIResponseQueryEdgeGateway) UnmarshalJSON(data []byte) error {
	var values []APIResponseEdgegateway
	if err := json.Unmarshal(data, &values); err == nil {
		r.Record = make([]APIResponseQueryEdgeGatewayRecord, 0, len(values))
		for _, value := range values {
			r.Record = append(r.Record, queryEdgeGatewayRecord(value))
		}
		return nil
	}

	var response struct {
		Record []APIResponseQueryEdgeGatewayRecord `json:"record"`
		Values []APIResponseEdgegateway            `json:"values"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	r.Record = response.Record
	for _, value := range response.Values {
		r.Record = append(r.Record, queryEdgeGatewayRecord(value))
	}
	return nil
}

func queryEdgeGatewayRecord(value APIResponseEdgegateway) APIResponseQueryEdgeGatewayRecord {
	id, name := value.ID, value.Name
	if id == "" {
		id = value.EdgeID
	}
	if name == "" {
		name = value.EdgeName
	}
	record := APIResponseQueryEdgeGatewayRecord{ID: id, Name: name}
	if value.OwnerRef != nil {
		record.VDCID = value.OwnerRef.ID
		record.VDCName = value.OwnerRef.Name
	}
	if record.VDCName == "" {
		record.VDCName = value.OwnerName
	}
	if value.OrgVDC != nil {
		record.VDCID = value.OrgVDC.ID
		record.VDCName = value.OrgVDC.Name
	}
	return record
}

// ToModel converts the APIResponseEdgegateways to ModelEdgeGateways.
func (api *APIResponseEdgegateways) ToModel() *types.ModelEdgeGateways {
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
func (api *APIResponseEdgegateway) ToModel() *types.ModelEdgeGateway {
	if api == nil {
		return nil
	}

	id, name := api.ID, api.Name
	if id == "" {
		id = api.EdgeID
	}
	if name == "" {
		name = api.EdgeName
	}
	ownerRef := api.OwnerRef
	if ownerRef == nil && (api.OwnerName != "" || api.OwnerType != "") {
		ownerRef = &APIObjectReference{Name: api.OwnerName}
	}
	uplinkT0 := func() *types.ModelObjectReference {
		if len(api.EdgeGatewayUplinks) > 0 {
			return &types.ModelObjectReference{
				ID:   api.EdgeGatewayUplinks[0].UplinkID,
				Name: api.EdgeGatewayUplinks[0].UplinkName,
			}
		}
		if api.Tier0VRFID != "" {
			return &types.ModelObjectReference{ID: api.Tier0VRFID}
		}
		return nil
	}()

	return &types.ModelEdgeGateway{
		ID:          id,
		Name:        name,
		Description: api.Description,
		OwnerRef: func() *types.ModelObjectReference {
			if ownerRef != nil {
				return &types.ModelObjectReference{
					ID:   ownerRef.ID,
					Name: ownerRef.Name,
				}
			}
			return nil
		}(),
		UplinkT0: uplinkT0,
	}
}
