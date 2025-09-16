/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
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
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/edgeGateways/gatewayId/get/
func GetEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("GetEdgeGateway")
}

// QueryEdgeGateway - Query EdgeGateway
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/38.1/doc/types/QueryResultEdgeGatewayRecordType.html
func QueryEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("QueryEdgeGateway")
}

// CreateEdgeGateway - Create EdgeGateway
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/createVdcEdge
func CreateEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateEdgeGateway")
}

// DeleteEdgeGateway - Delete EdgeGateway
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/deleteEdge
func DeleteEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteEdgeGateway")
}

// ListEdgeGateway - List EdgeGateways
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/edgeGateways/get/
func ListEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("ListEdgeGateway")
}
