/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path edgegateway_publicip.go -output edgegateway_publicip

func init() {
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/addNetworkConnectivity",
		Name:             "CreatePublicIp",
		Description:      "Create a new public IP",
		Method:           cav.MethodPOST,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/services",
		BodyResponseType: cav.Job{},
		BodyRequestType:  itypes.ApiRequestEdgegatewayPublicIP{},
	}.Register()
}
