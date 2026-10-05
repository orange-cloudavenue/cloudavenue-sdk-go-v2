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

//go:generate endpoint-generator -path edgegateway_services.go -output edgegateway_services

func init() {
	// * GetEdgeGatewayServices
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/getNetworkHierarchy",
		Name:             "GetEdgeGatewayServices",
		Description:      "Get EdgeGateway Network Services",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/network",
		ResponseType:     itypes.APIResponseNetworkServices{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/addNetworkConnectivity",
		Name:             "EnableCloudavenueServices",
		Description:      "Enable Cloud Avenue Services",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/services",
		ResponseType:     cav.Job{},
		BodyRequestType:  itypes.APIRequestNetworkServicesCavSvc{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/deleteNetworkService",
		Name:             "DisableCloudavenueServices",
		Description:      "Disable Cloud Avenue Services",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/services/{serviceId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamServiceID,
				Description: "The ID of the service to delete",
				Required:    true,
			},
		},
		ResponseType: cav.Job{},
	}.Register()
}
