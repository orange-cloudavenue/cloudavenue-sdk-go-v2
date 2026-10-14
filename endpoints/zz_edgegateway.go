/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// GetEdgeGateway - Get EdgeGateway
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/getEdgeById
func GetEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("GetEdgeGateway")
}

// QueryEdgeGateway - List EdgeGateways (compatibility alias)
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/getEdges
func QueryEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("QueryEdgeGateway")
}

// CreateEdgeGateway - Create EdgeGateway in a VDC
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/createVdcEdge
func CreateEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateEdgeGateway")
}

// CreateEdgeGatewayVDCGroup - Create EdgeGateway in a VDC group
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/createVdcGroupEdge
func CreateEdgeGatewayVDCGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateEdgeGatewayVDCGroup")
}

// DeleteEdgeGateway - Delete EdgeGateway
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/deleteEdge
func DeleteEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteEdgeGateway")
}

// ListEdgeGateway - List EdgeGateways
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/getEdges
func ListEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("ListEdgeGateway")
}
