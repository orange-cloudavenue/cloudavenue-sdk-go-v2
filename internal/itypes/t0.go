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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

type (
	// * APIResponse
	APIResponseT0Names []string

	APIResponseT0 struct {
		Tier0VRF      string `json:"tier0_vrf" fake:"{resource_name:t0}"`
		Tier0Provider string `json:"tier0_provider" fake:"pr01e02t0sp01"`
		// One of VRF_STANDARD, VRF_PREMIUM, VRF_DEDICATED_MEDIUM, VRF_DEDICATED_LARGE.
		ClassService string                 `json:"class_service" fake:"VRF_STANDARD"`
		Services     []APIResponseT0Service `json:"services" fakesize:"1"`
	}

	APIResponseT0Service struct {
		Service string `json:"service" fake:"INTERNET"`
		VLANID  string `json:"vlan_id,omitempty" fake:"1234"`
	}
)

func (t0Names APIResponseT0Names) ToModel() *types.ModelT0s {
	modelT0s := &types.ModelT0s{}
	for _, name := range t0Names {
		modelT0s.T0s = append(modelT0s.T0s, types.ModelT0{Name: name})
	}
	modelT0s.Count = len(modelT0s.T0s)
	return modelT0s
}

// edgeGateways returns the edge gateways attached to the T0 named name in the network hierarchy.
// The Tier-0 details endpoint does not expose edge gateways, they are only available in the network hierarchy.
func (hierarchy *APIResponseNetworkServices) edgeGateways(name string) []APIResponseNetworkServicesChildren {
	if hierarchy == nil {
		return nil
	}

	for _, t0 := range *hierarchy {
		if t0.Name != name {
			continue
		}

		var edgeGateways []APIResponseNetworkServicesChildren
		for _, child := range t0.Children {
			if child.Type == "edge-gateway" {
				edgeGateways = append(edgeGateways, child)
			}
		}
		return edgeGateways
	}

	return nil
}

// ToModel converts the Tier-0 details to a model.
// The hierarchy is the response of the network hierarchy endpoint and is used to
// retrieve the edge gateways attached to the T0 (nil means no edge gateway).
func (t0 APIResponseT0) ToModel(hierarchy *APIResponseNetworkServices) *types.ModelT0 {
	cos, hasCoS := classOfServices[t0.ClassService]

	model := &types.ModelT0{
		Name:            t0.Tier0VRF,
		ClassOfService:  t0.ClassService,
		MaxEdgeGateways: cos.MaxEdgeGateways,
		Bandwidth: types.ModelT0Bandwidth{
			Capacity:       cos.MaxBandwidth,
			AllowUnlimited: cos.AllowUnlimited,
		},
	}

	for _, child := range hierarchy.edgeGateways(t0.Tier0VRF) {
		if child.Properties.RateLimit > 0 {
			model.Bandwidth.Provisioned += child.Properties.RateLimit
		}

		model.EdgeGateways = append(model.EdgeGateways, types.ModelT0EdgeGateway{
			ID:        child.Properties.EdgeUUID,
			Name:      child.Name,
			Bandwidth: child.Properties.RateLimit,
		})
	}

	model.Bandwidth.Remaining = model.Bandwidth.Capacity - model.Bandwidth.Provisioned

	// Calculate the allowed bandwidth values, they must fit in the remaining bandwidth.
	if hasCoS && model.Bandwidth.Remaining > 0 {
		for _, bw := range cos.MaxEdgeGatewayBandwidth {
			if bw <= model.Bandwidth.Remaining {
				model.Bandwidth.AllowedBandwidthValues = append(model.Bandwidth.AllowedBandwidthValues, bw)
			}
		}

		for i := range model.EdgeGateways {
			model.EdgeGateways[i].AllowedBandwidthValues = append([]int(nil), model.Bandwidth.AllowedBandwidthValues...)
		}
	}

	return model
}
