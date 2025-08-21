/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

type (
	ModelEdgeGatewayPublicIPs struct {
		EdgegatewayID   string `documentation:"ID of the edge gateway"`
		EdgegatewayName string `documentation:"Name of the edge gateway"`

		PublicIPs []ModelEdgeGatewayPublicIP `documentation:"List of public IPs"`
	}

	ModelEdgeGatewayPublicIP struct {
		EdgegatewayID   string `documentation:"ID of the edge gateway"`
		EdgegatewayName string `documentation:"Name of the edge gateway"`

		ModelEdgeGatewayServicesPublicIP
	}

	ParamsGetEdgeGatewayPublicIP struct {
		ID   string `fake:"{urn:edgegateway}"`
		Name string `fake:"{resource_name:edgegateway}"`

		IP string `validate:"ipv4" fake:"{IPv4Address}"`
	}

	ParamsDeleteEdgeGatewayPublicIP struct {
		IP string `validate:"ipv4" fake:"{IPv4Address}"`
	}
)
