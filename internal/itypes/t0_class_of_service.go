/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

type (
	classOfService struct {
		// MaxBandwidth defines the maximum bandwidth in Mbps for this class of service.
		MaxBandwidth int

		// MaxEdgeGateways defines the maximum number of edge gateways allowed for this class of service.
		MaxEdgeGateways int

		// MaxEdgeGatewayBandwidth defines the values of bandwidth that can be allocated to each edge gateway.
		// This is a list of integers representing the allowed bandwidth values in Mbps.
		MaxEdgeGatewayBandwidth []int

		// Allow unlimited bandwidth for edge gateways.
		AllowUnlimited bool
	}
)

// Source : https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/network/
var classOfServices = map[string]classOfService{
	"SHARED_STANDARD": {
		MaxBandwidth:            300,
		MaxEdgeGateways:         4,
		MaxEdgeGatewayBandwidth: []int{5, 25, 50, 75, 100, 150, 200, 250, 300},
		AllowUnlimited:          false,
	},
	"SHARED_PREMIUM": {
		MaxBandwidth:            1000,
		MaxEdgeGateways:         8,
		MaxEdgeGatewayBandwidth: []int{5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000},
		AllowUnlimited:          false,
	},
	"DEDICATED_MEDIUM": {
		MaxBandwidth:            3500,
		MaxEdgeGateways:         6,
		MaxEdgeGatewayBandwidth: []int{5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000, 2000},
		AllowUnlimited:          true, // Unlimited bandwidth is allowed for edge gateway with this class of service.
	},
	"DEDICATED_LARGE": {
		MaxBandwidth:            10000,
		MaxEdgeGateways:         12,
		MaxEdgeGatewayBandwidth: []int{5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 6000},
		AllowUnlimited:          true, // Unlimited bandwidth is allowed for edge gateway with this class of service.
	},
}
