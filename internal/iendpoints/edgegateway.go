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

const (
	pathCerberusEdgeGatewayVDCCreate      = "/infrapicustomerproxy/v2.0/vdcs/{vdcName}/edges"
	pathCerberusEdgeGatewayVDCGroupCreate = "/infrapicustomerproxy/v2.0/vdc-groups/{vdcGroupName}/edges"
)

//go:generate endpoint-generator -path edgegateway.go -output edgegateway

func init() {
	// GetEdgeGateway
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/getEdgeById",
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
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/getEdges",
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
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/createVdcEdge",
		Name:             "CreateEdgeGateway",
		Description:      "Create EdgeGateway in a VDC",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusEdgeGatewayVDCCreate,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCName,
				Description: "The name of the VDC where the edge gateway will be created.",
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestEdgeGateway{},
		ResponseType:    cav.CerberusJobCreatedAPIResponse{},
	}.Register()

	// CreateEdgeGatewayVDCGroup
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/createVdcGroupEdge",
		Name:             "CreateEdgeGatewayVDCGroup",
		Description:      "Create EdgeGateway in a VDC group",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusEdgeGatewayVDCGroupCreate,
		PathParams: []cav.PathParam{
			{
				Name:        "vdcGroupName",
				Description: "The name of the VDC group where the edge gateway will be created.",
				Required:    true,
			},
		},
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
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/getEdges",
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
