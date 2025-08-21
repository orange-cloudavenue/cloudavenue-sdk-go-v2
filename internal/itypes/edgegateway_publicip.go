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
	ApiRequestEdgegatewayPublicIP struct {
		// NetworkType
		NetworkType string `json:"networkType" default:"internet" validate:"required"` // The type of network service to create (load-balancer, service, internet)

		// EdgeGatewayID - The ID of the edge gateway is a UUID and not a URN.
		EdgeGatewayID string `json:"edgeGateway" validate:"required,uuid"`

		Properties ApiRequestEdgegatewayPublicIPProperties `json:"properties" validate:"omitempty"`
	}

	ApiRequestEdgegatewayPublicIPProperties struct {
		// Announced represents if the public IP address is announced
		Announced bool `json:"announced" validate:"omitempty"`
	}
)
