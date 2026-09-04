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
	"fmt"
	"regexp"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path edgegateway.go -output edgegateway

func init() {
	// GetEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/edgeGateways/gatewayId/get/",
		Name:             "GetEdgeGateway",
		Description:      "Get EdgeGateway",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/edgeGateways/{edgeId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamEdgeID,
				Description: descEdgeGatewayID,
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnEdgeGateway)
				},
			},
		},
		ResponseType: itypes.APIResponseEdgegateway{},
	}.Register()

	// QueryEdgeGateway
	cav.Endpoint{
		// "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-ExecuteQuery.html"
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/types/QueryResultEdgeGatewayRecordType.html",
		Name:             "QueryEdgeGateway",
		Description:      "Query EdgeGateway",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathQueryAPI,
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamType,
				Description: descTypeOfObjectQuery,
				Value:       typeEdgeGateway,
			},
			{
				Name:        queryParamFilter,
				Description: "The filter to apply to the query",
				Required:    false,
				ValidatorFunc: func(value string) error {
					// check if the value is a valid key==value pair
					x := regexp.MustCompile(`^[a-zA-Z0-9_]+==.*`)
					if !x.MatchString(value) {
						return fmt.Errorf("invalid filter format, expected key==value")
					}

					return nil
				},
			},
		},
		PathParams:      nil,
		BodyRequestType: nil,
		ResponseType:    itypes.APIResponseQueryEdgeGateway{},
	}.Register()

	// CreateEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/createVdcEdge",
		Name:             "CreateEdgeGateway",
		Description:      "Create EdgeGateway",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/{vdc-type}/{vdc-name}/edges",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-type",
				Description: "The type of the VDC where the edge gateway will be created.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "oneof=vdc vdcgroup")
				},
				TransformFunc: func(value string) (string, error) {
					switch value {
					case queryParamVDC:
						return "vdcs", nil
					case "vdcgroup":
						return "vdc-groups", nil
					}
					return "", fmt.Errorf("invalid vdc-type: %s", value)
				},
			},
			{
				Name:        pathParamVDCName,
				Description: "The name of the VDC where the edge gateway will be created.",
				Required:    true,
			},
		},
		QueryParams:     nil,
		BodyRequestType: itypes.APIRequestEdgeGateway{},
		ResponseType:    cav.CerberusJobCreatedAPIResponse{},
	}.Register()

	// DeleteEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/deleteEdge",
		Name:             "DeleteEdgeGateway",
		Description:      "Delete EdgeGateway",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/edges/{edgeId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamEdgeID,
				Description: descEdgeGatewayID,
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, ruleRequiredURNEdgeGateway)
				},
				TransformFunc: extractor.ExtractUUID,
			},
		},
		QueryParams:     nil,
		BodyRequestType: nil,
		ResponseType:    cav.CerberusJobCreatedAPIResponse{},
	}.Register()

	// ListEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/edgeGateways/get/",
		Name:             "ListEdgeGateway",
		Description:      "List EdgeGateways",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/edgeGateways",
		PathParams:       nil,
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamPageSize,
				Description: "The number of items to return per page.",
				Value:       pageSize128,
			},
		},
		ResponseType: itypes.APIResponseEdgegateways{},
	}.Register()
}
