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

// UpdateEdgeGatewayBandwidth - Update EdgeGateway Bandwidth
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/put_api_customers_v2_0_edges__edge_id_ 
func UpdateEdgeGatewayBandwidth() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateEdgeGatewayBandwidth")
}
