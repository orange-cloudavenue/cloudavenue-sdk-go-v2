/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

const (
	// pathCerberusDraasOnPremiseIP addresses one on-premise DRAAS IP.
	pathCerberusDraasOnPremiseIP = "/infrapicustomerproxy/v2.0/vcda/ips/{ip}"

	// pathCerberusDraasOnPremiseIPs lists DRAAS on-premise IPs.
	pathCerberusDraasOnPremiseIPs = "/infrapicustomerproxy/v2.0/vcda/ips"

	// pathCerberusEdgeGatewayByID addresses one edge gateway by ID.
	pathCerberusEdgeGatewayByID = "/infrapicustomerproxy/v2.0/edges/{edgeId}"

	// pathCerberusEdgeGateways lists edge gateways.
	pathCerberusEdgeGateways = "/infrapicustomerproxy/v2.0/edges"

	// pathCerberusEdgeGatewayCreate creates an edge gateway under a VDC or VDC group.
	pathCerberusEdgeGatewayCreate = "/infrapicustomerproxy/v2.0/{vdcType}/{vdcName}/edges"

	// pathCerberusCloudavenueServices lists and creates CloudAvenue services.
	pathCerberusCloudavenueServices = "/infrapicustomerproxy/v2.0/services"

	// pathCerberusCloudavenueServiceByID addresses one CloudAvenue service by ID.
	pathCerberusCloudavenueServiceByID = "/infrapicustomerproxy/v2.0/services/{serviceId}"

	// pathCerberusNetwork manages edge network-related operations.
	pathCerberusNetwork = "/infrapicustomerproxy/v2.0/network"

	// pathCerberusConfigurations gets or updates organization configuration.
	pathCerberusConfigurations = "/infrapicustomerproxy/v2.0/configurations"

	// pathCerberusVDCs lists and creates VDCs.
	pathCerberusVDCs = "/infrapicustomerproxy/v2.0/vdcs"

	// pathCerberusVDCByName addresses one VDC by name.
	pathCerberusVDCByName = "/infrapicustomerproxy/v2.0/vdcs/{vdcName}"

	// pathCerberusT0s lists Tier-0 VRFs.
	pathCerberusT0s = "/infrapicustomerproxy/v2.0/tier-0-vrfs"

	// pathCerberusT0ByName addresses one Tier-0 VRF by name.
	pathCerberusT0ByName = "/infrapicustomerproxy/v2.0/tier-0-vrfs/{tier0Name}"
)

const (
	pathParamEdgeID    = "edgeId"
	pathParamServiceID = "serviceId"
	pathParamVDCName   = "vdcName"
)
