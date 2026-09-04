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

//go:generate endpoint-generator -path edgegateway_bandwidth.go -output edgegateway_bandwidth

func init() {
	// * UpdateEdgeGatewayBandwidth
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Edge%20Gateways/put_api_customers_v2_0_edges__edge_id_",
		Name:             "UpdateEdgeGatewayBandwidth",
		Description:      "Update EdgeGateway Bandwidth",
		Method:           cav.MethodPUT,
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
		BodyRequestType: itypes.APIRequestBandwidth{},
		ResponseType:    cav.Job{},
	}.Register()
}
