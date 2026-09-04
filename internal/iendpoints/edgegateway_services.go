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
	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"

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
		QueryParams: []cav.QueryParam{
			// Query parameters are not used in this endpoint, but can be added
			// for the mock response if needed
			{
				Name:        pathParamEdgeID,
				Description: "The ID of the edge gateway to get network services for",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnEdgeGateway)
				},
				TransformFunc: extractor.ExtractUUID,
			},
			{
				Name:        "edgeName",
				Description: "The name of the edge gateway to get network services for",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, ruleResourceNameEdgeGateway)
				},
			},
			{
				Name:        "publicIP",
				Description: "The public IP address of the edge gateway",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "ipv4")
				},
			},
		},
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
