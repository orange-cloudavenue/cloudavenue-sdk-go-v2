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
		PublicIps []ModelEdgeGatewayPublicIP `documentation:"List of public IPs"`
	}

	ModelEdgeGatewayPublicIP struct {
		ID string `documentation:"Unique identifier of the public IP"`

		EdgeGatewayID   string `documentation:"Unique identifier of the edge gateway"`
		EdgeGatewayName string `documentation:"Name of the edge gateway"`

		IP        string `documentation:"Public IPv4 address"`
		Announced bool   `documentation:"True if the public IP is advertised via BGP"`
	}

	ParamsListEdgeGatewayPublicIP struct {
		EdgeGatewayID   string `fake:"{urn:edgegateway}"`
		EdgeGatewayName string `fake:"{resource_name:edgegateway}"`
	}

	ParamsCreateEdgeGatewayPublicIP struct {
		EdgeGatewayID   string `fake:"{urn:edgegateway}"`
		EdgeGatewayName string `fake:"{resource_name:edgegateway}"`
	}

	ParamsGetEdgeGatewayPublicIP struct {
		EdgeGatewayID   string `fake:"{urn:edgegateway}"`
		EdgeGatewayName string `fake:"{resource_name:edgegateway}"`

		IP string `validate:"ipv4" fake:"{IPv4Address}"`
	}

	ParamsDeleteEdgeGatewayPublicIP struct {
		IP string `validate:"ipv4" fake:"{IPv4Address}"`
	}
)
