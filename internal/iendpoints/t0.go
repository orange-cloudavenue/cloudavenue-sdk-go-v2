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
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path t0.go -output t0

func init() {
	// GET - List all T0
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/getNetworkHierarchy",
		Name:             "ListT0",
		Description:      "List T0",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/network",
		ResponseType:     itypes.ApiResponseT0s{},

		QueryParams: []cav.QueryParam{
			// Query parameters are not used in this endpoint, but can be added
			// for the mock response if needed
			{
				Name:        "t0Name",
				Description: "The name of the T0",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "resource_name=t0")
				},
			},
			{
				Name:        "edgeGatewayName",
				Description: "The name of the Edge Gateway",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "resource_name=edgegateway")
				},
			},
			{
				Name:        "edgeGatewayID",
				Description: "The ID of the Edge Gateway",
				Required:    false,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=edgegateway")
				},
			},
		},
	}.Register()
}
