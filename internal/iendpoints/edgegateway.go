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

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path edgegateway.go -output edgegateway

func init() {
	// GetEdgeGateway
	cav.Endpoint{
		DocumentationURL: docURLCerberus,
		Name:             "GetEdgeGateway",
		Description:      "Get EdgeGateway",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusEdgeGatewayByID,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamEdgeID,
				Description: descEdgeGatewayID,
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnEdgeGateway)
				},
				TransformFunc: extractor.ExtractUUID,
			},
		},
		ResponseType: itypes.APIResponseEdgegateway{},
	}.Register()

	// QueryEdgeGateway
	cav.Endpoint{
		DocumentationURL: docURLCerberus,
		Name:             "QueryEdgeGateway",
		Description:      "List EdgeGateways (compatibility alias)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusEdgeGateways,
		PathParams:       nil,
		QueryParams:      nil,
		BodyRequestType:  nil,
		ResponseType:     itypes.APIResponseQueryEdgeGateway{},
	}.Register()

	// CreateEdgeGateway
	cav.Endpoint{
		DocumentationURL: docURLCerberus,
		Name:             "CreateEdgeGateway",
		Description:      "Create EdgeGateway",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusEdgeGatewayCreate,
		PathParams: []cav.PathParam{
			{
				Name:        "vdcType",
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
					return "", fmt.Errorf("invalid vdcType: %s", value)
				},
			},
			{
				Name:        pathParamVdcName,
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
		DocumentationURL: docURLCerberus,
		Name:             "DeleteEdgeGateway",
		Description:      "Delete EdgeGateway",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusEdgeGatewayByID,
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
		DocumentationURL: docURLCerberus,
		Name:             "ListEdgeGateway",
		Description:      "List EdgeGateways",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusEdgeGateways,
		PathParams:       nil,
		QueryParams:      nil,
		ResponseType:     itypes.APIResponseEdgegateways{},
	}.Register()
}
