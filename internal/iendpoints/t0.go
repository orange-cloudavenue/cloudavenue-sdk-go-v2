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
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Tier-0/getTier0Vrfs",
		Name:             "ListT0",
		Description:      "List T0",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/infrapicustomerproxy/v2.0/tier-0-vrfs",
		ResponseType:     itypes.APIResponseT0Names{},
	}.Register()

	// GET - T0 details
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Tier-0/getTier0VrfByName",
		Name:             "GetT0",
		Description:      "Get T0",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/infrapicustomerproxy/v2.0/tier-0-vrfs/{tier0_name}",
		PathParams: []cav.PathParam{
			{
				Name:        "tier0_name",
				Description: "The name of the T0",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "resource_name=t0")
				},
			},
		},
		ResponseType: itypes.APIResponseT0{},
	}.Register()
}
