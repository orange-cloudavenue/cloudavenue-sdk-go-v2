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
	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

type (

	// * APIResponse
	APIResponseNetworkServices []struct {
		Type     string                               `json:"type" fake:"tier-0-vrf"`
		Name     string                               `json:"name" fake:"{resource_name:t0}"`
		Children []APIResponseNetworkServicesChildren `json:"children,omitempty" fakesize:"1"`
	}

	APIResponseNetworkServicesChildren struct {
		Type        string `json:"type" fake:"edge-gateway"`
		Name        string `json:"name,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Properties  struct {
			// EdgeGateway
			RateLimit int    `json:"rateLimit,omitempty"`
			EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"` // The UUID of the edge gateway
		} `json:"properties"`
		Children  []APIResponseNetworkServicesSubChildren `json:"children,omitempty" fakesize:"6"`
		ServiceID string                                  `json:"serviceId,omitempty"`
	}
	APIResponseNetworkServicesSubChildren struct {
		Type        string `json:"type" fake:"{randomstring:[load-balancer,service]}"`
		Name        string `json:"name,omitempty" fake:"{randomstring:[cav-services,internet]}"`
		DisplayName string `json:"displayName,omitempty" fake:"{word}"`
		Properties  struct {
			// Load Balancer
			ClassOfService     string `json:"classOfService,omitempty"`
			MaxVirtualServices int    `json:"maxVirtualServices,omitempty"`

			// Public IP
			IP        string `json:"ip,omitempty" fake:"{ipv4address}"`
			Announced bool   `json:"announced,omitempty" fake:"true"`

			// Service
			Ranges []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"` // The network in ip/cidr format
		} `json:"properties"`
		ServiceID string `json:"serviceId,omitempty"`
	}

	// * APIRequest

	APIRequestNetworkServicesCavSvc struct {
		// NetworkType
		NetworkType string `json:"networkType" default:"cav-services" validate:"required"` // The type of network service to create (load-balancer, service, internet)

		// EdgeGatewayID - The ID of the edge gateway is a UUID and not a URN.
		EdgeGatewayID string `json:"edgeGateway" validate:"required,uuid"`

		// Properties
		Properties struct {
			PrefixLength int `json:"prefixLength,omitempty" validate:"omitempty,min=25,max=28" default:"27"` // The prefix length of the network in CIDR notation
		}
	}
)

func (ap *APIResponseNetworkServices) ToModel(params types.ParamsEdgeGateway) *types.ModelEdgeGatewayServices {
	if ap == nil || len(*ap) == 0 {
		return nil
	}

	data := &types.ModelEdgeGatewayServices{}
	for _, ns := range *ap {
		for _, child := range ns.Children {
			if child.Type != "edge-gateway" || (child.Properties.EdgeUUID != urn.ExtractUUID(params.ID) && child.Name != params.Name) {
				continue
			}

			data.ID = urn.Normalize(urn.EdgeGateway, child.Properties.EdgeUUID).String()
			data.Name = child.Name
			populateEdgeGatewayServices(data, child.Children)
		}
	}

	return data
}
