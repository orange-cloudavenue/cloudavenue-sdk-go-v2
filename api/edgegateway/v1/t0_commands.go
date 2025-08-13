/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

//go:generate command-generator -path t0_commands.go

func init() {
	// ! TO
	cmds.Register(commands.Command{
		Namespace: "T0",
	})

	// * ListT0
	cmds.Register(commands.Command{
		Namespace: "T0",
		Verb:      "List",

		ShortDocumentation: "List all T0s available in the organization.",
		LongDocumentation:  "List all T0s available in the organization. This command retrieves a list of all T0s, which are the top-level network services in the Edge Gateway architecture.",
		AutoGenerate:       true,

		ModelType: types.ModelT0s{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			ep := endpoints.ListT0()

			// Perform the request to list T0s
			resp, err := cc.c.Do(ctx, ep)
			if err != nil {
				return nil, fmt.Errorf("error listing T0s: %w", err)
			}

			return resp.Result().(*itypes.ApiResponseT0s).ToModel(), nil
		},
	})

	// * GetT0
	cmds.Register(commands.Command{
		Namespace:          "T0",
		Verb:               "Get",
		ShortDocumentation: "Get a specific T0 by name.",
		LongDocumentation:  "Retrieve a specific T0 directly by its name or by the edge gateway it is associated with. This command allows you to fetch detailed information about a specific T0.",
		AutoGenerate:       true,

		ParamsType: types.ParamsGetT0{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "t0_name",
				Description: "The name of the T0 to retrieve.",
				Required:    false,
				Example:     "prvrf01eocb0001234allsp01",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("edgegateway_id", "edgegateway_name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorResourceName("t0"),
				},
			},
			commands.ParamsSpec{
				Name:        "edgegateway_id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("t0_name", "edgegateway_name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "edgegateway_name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Example:     "tn01e02ocb0001234spt101",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("edgegateway_id", "t0_name"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		ModelType: types.ModelT0{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsGetT0)
			ep := endpoints.ListT0()

			// Perform the request to get the specific T0
			resp, err := cc.c.Do(
				ctx,
				ep,
				cav.WithQueryParam(ep.QueryParams[0], p.T0Name),          // Only for mock response
				cav.WithQueryParam(ep.QueryParams[1], p.EdgegatewayName), // Only for mock response
				cav.WithQueryParam(ep.QueryParams[2], p.EdgegatewayID),   // Only for mock response
			)
			if err != nil {
				return nil, fmt.Errorf("error getting T0: %w", err)
			}

			t0s := resp.Result().(*itypes.ApiResponseT0s).ToModel()
			var t0 *types.ModelT0

			for _, t := range t0s.T0s {
				if p.T0Name != "" && t.Name == p.T0Name {
					t0 = &t
					break
				}
				if p.EdgegatewayID != "" || p.EdgegatewayName != "" {
					for _, edgeGateway := range t.EdgeGateways {
						if p.EdgegatewayID == edgeGateway.ID || p.EdgegatewayName == edgeGateway.Name {
							t0 = &t
							break
						}
					}
				}
			}

			if t0 == nil {
				return nil, &errors.APIError{
					Operation:     "GetT0",
					StatusCode:    http.StatusNotFound,
					StatusMessage: http.StatusText(http.StatusNotFound),
					Message: func() string {
						if p.T0Name != "" {
							return fmt.Sprintf("T0 with name %s not found", p.T0Name)
						}
						if p.EdgegatewayID != "" {
							return fmt.Sprintf("T0 for edge gateway with ID %s not found", p.EdgegatewayID)
						}
						return fmt.Sprintf("T0 for edge gateway with name %s not found", p.EdgegatewayName)
					}(),
					Duration: resp.Duration(),
					Endpoint: resp.Request.URL,
					Method:   resp.Request.Method,
				}
			}

			return t0, nil
		},
	})
}
