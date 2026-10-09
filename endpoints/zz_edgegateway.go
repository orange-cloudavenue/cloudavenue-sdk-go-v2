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
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func GetEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("GetEdgeGateway")
}
// QueryEdgeGateway - List EdgeGateways (compatibility alias)
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func QueryEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("QueryEdgeGateway")
}
// CreateEdgeGateway - Create EdgeGateway
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func CreateEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateEdgeGateway")
}
// DeleteEdgeGateway - Delete EdgeGateway
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func DeleteEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteEdgeGateway")
}
// ListEdgeGateway - List EdgeGateways
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func ListEdgeGateway() *cav.Endpoint {
	return cav.MustGetEndpoint("ListEdgeGateway")
}

