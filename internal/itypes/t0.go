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
	"github.com/orange-cloudavenue/common-go/urn"
)

type (
	// * ApiResponse
	ApiResponseT0s []ApiResponseT0

	ApiResponseT0 struct {
		Type       string                  `json:"type" fake:"tier-0-vrf"`
		Name       string                  `json:"name" fake:"{resource_name:t0}"`
		Properties ApiResponseT0Properties `json:"properties,omitempty"`
		Children   []ApiResponseT0Children `json:"children,omitempty" fakesize:"1"`
	}

	ApiResponseT0Properties struct {
		ClassOfService string `json:"classOfService,omitempty" fake:"SHARED_STANDARD"`
	}

	ApiResponseT0Children struct {
		Type       string `json:"type" fake:"edge-gateway"`
		Name       string `json:"name" fake:"{resource_name:edgegateway}"`
		Properties struct {
			RateLimit int    `json:"rateLimit,omitempty" fake:"5"`
			EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"` // The UUID of the edge gateway
		} `json:"properties,omitempty"`
	}
)

func (t0s ApiResponseT0s) ToModel() *types.ModelT0s {
	var modelT0s types.ModelT0s
	for _, t0 := range t0s {
		modelT0 := types.ModelT0{
			Name:           t0.Name,
			ClassOfService: t0.Properties.ClassOfService,
			Bandwidth: types.ModelT0Bandwidth{
				Capacity: func() int {
					if cof, ok := classOfServices[t0.Properties.ClassOfService]; ok {
						return cof.MaxBandwidth
					}
					return 0
				}(),
				Provisioned: func() int {
					if t0.Children == nil {
						return 0
					}

					var totalProvisioned int
					for _, child := range t0.Children {
						if child.Properties.RateLimit > 0 {
							totalProvisioned += child.Properties.RateLimit
						}
					}
					return totalProvisioned
				}(),
				Remaining:              0,   // Remaining is calculated later
				AllowedBandwidthValues: nil, // AllowedBandwidthValues is calculated later
				AllowUnlimited: func() bool {
					if cof, ok := classOfServices[t0.Properties.ClassOfService]; ok {
						return cof.AllowUnlimited
					}
					return false
				}(),
			},
			Summary: func() types.ModelT0Summary {
				t0s := types.ModelT0Summary{}
				if t0.Children != nil {
					t0s.EdgeGateways = len(t0.Children)
				}
				if cof, ok := classOfServices[t0.Properties.ClassOfService]; ok {
					t0s.MaxEdgeGateways = cof.MaxEdgeGateways
				}
				return t0s
			}(),
			EdgeGateways: func() []types.ModelT0EdgeGateway {
				var edgeGateways []types.ModelT0EdgeGateway

				for _, child := range t0.Children {
					if child.Type == "edge-gateway" {
						edgeGateway := types.ModelT0EdgeGateway{
							ID:                     urn.Normalize(urn.EdgeGateway, child.Properties.EdgeUUID).String(),
							Name:                   child.Name,
							Bandwidth:              child.Properties.RateLimit,
							AllowedBandwidthValues: nil, // This will be calculated later based on the remaining bandwidth
						}
						edgeGateways = append(edgeGateways, edgeGateway)
					}
				}
				return edgeGateways
			}(),
		}
		modelT0s.T0s = append(modelT0s.T0s, modelT0)
	}

	// Calculate the remaining bandwidth and allowed bandwidth values for each T0 router
	for i, t0 := range modelT0s.T0s {
		t0.Bandwidth.Remaining = t0.Bandwidth.Capacity - t0.Bandwidth.Provisioned

		// Calculate the allowed bandwidth values for the T0 router
		if cof, ok := classOfServices[t0.ClassOfService]; ok {
			// Use t0.Bandwidth.Remaining to ensure the values are within the allowed range
			if t0.Bandwidth.Remaining > 0 {
				for _, bw := range cof.MaxEdgeGatewayBandwidth {
					if bw <= t0.Bandwidth.Remaining {
						t0.Bandwidth.AllowedBandwidthValues = append(t0.Bandwidth.AllowedBandwidthValues, bw)
					}
				}
			}
		}
		modelT0s.T0s[i] = t0
	}

	// Calculate the allowed bandwidth values for each edge gateway
	for i, t0 := range modelT0s.T0s {
		for j, edgeGateway := range t0.EdgeGateways {
			// Use t0.Bandwidth.Remaining to ensure the values are within the allowed range
			if t0.Bandwidth.Remaining > 0 {
				for _, bw := range classOfServices[t0.ClassOfService].MaxEdgeGatewayBandwidth {
					if bw <= t0.Bandwidth.Remaining {
						edgeGateway.AllowedBandwidthValues = append(edgeGateway.AllowedBandwidthValues, bw)
					}
				}
			}
			modelT0s.T0s[i].EdgeGateways[j] = edgeGateway
		}
	}

	modelT0s.Count = len(modelT0s.T0s)
	return &modelT0s
}
