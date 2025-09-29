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
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/urn"
)

//go:generate command-generator -path services_commands.go

func init() {
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Resource:  "Services",
	})

	// * GetServices
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Resource:  "Services",
		Verb:      "Get",

		ShortDocumentation: "Retrieve services information about a specific EdgeGateway.",
		LongDocumentation:  "Retrieve services information about a specific EdgeGateway. This command retrieves the network services available on the EdgeGateway, such as load balancers, public IPs, and Cloud Avenue services.",
		AutoGenerate:       true,

		ParamsType: types.ParamsEdgeGateway{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		ModelType: types.ModelEdgeGatewayServices{},

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsEdgeGateway)
			ep := endpoints.GetEdgeGatewayServices()

			// ID is required to request the API.
			if p.ID == "" {
				var err error
				p.ID, err = cc.retrieveEdgeGatewayIDByName(ctx, p.Name)
				if err != nil {
					return nil, err
				}
			}

			resp, err := cc.c.Do(
				ctx,
				ep,
				cav.WithQueryParam(ep.QueryParams[0], p.ID),   // Only for filtering the response
				cav.WithQueryParam(ep.QueryParams[1], p.Name), // Only for filtering the response
			)
			if err != nil {
				return nil, fmt.Errorf("error retrieving network services for edge gateway %s: %w", p.ID, err)
			}

			data := resp.Result().(*itypes.ApiResponseNetworkServices).ToModel(p)
			if data == nil {
				return nil, fmt.Errorf("no network services found for edge gateway %s", p.ID)
			}

			return data, nil
		},
	})

	// ! CloudavenueServices
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Resource:  "CloudavenueServices",

		MarkdownDocumentation: "Manage Cloud Avenue services on an EdgeGateway. Cloudavenue services is a network setting that allows the EdgeGateway to connect to the mutualized Cloud Avenue services (DNS, DHCP, etc.).",
	})

	// * GetCloudavenueServices
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Resource:           "CloudavenueServices",
		Verb:               "Get",
		ShortDocumentation: "Retrieve Cloud Avenue services on an EdgeGateway.",
		LongDocumentation:  "Retrieve Cloud Avenue services on an EdgeGateway. This command returns the Cloud Avenue services available on the EdgeGateway, such as DNS, DHCP, and others.",
		AutoGenerate:       true,

		ParamsType: types.ParamsEdgeGateway{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Example:     generator.MustGenerate("{urn:edgegateway}"),
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
		},

		ModelType: types.ModelCloudavenueServices{},

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsEdgeGateway)

			svcs, err := cc.GetServices(ctx, p)
			if err != nil {
				return nil, err
			}

			return svcs.Services, nil
		},
	})

	// * EnableCloudavenueServices
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Resource:           "CloudavenueServices",
		Verb:               "Enable",
		ShortDocumentation: "Enable Cloud Avenue services on an EdgeGateway.",
		LongDocumentation:  "Enable Cloud Avenue services on an EdgeGateway. ",
		AutoGenerate:       true,

		ParamsType: types.ParamsEdgeGateway{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Example:     generator.MustGenerate("{urn:edgegateway}"),
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
		},

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsEdgeGateway)
			ep := endpoints.EnableCloudavenueServices()

			// ID is required to request the API.
			if p.ID == "" {
				var err error
				p.ID, err = cc.retrieveEdgeGatewayIDByName(ctx, p.Name)
				if err != nil {
					return nil, err
				}
			}

			// Prepare the request body
			requestBody := &itypes.ApiRequestNetworkServicesCavSvc{
				NetworkType:   "cav-services",
				EdgeGatewayID: urn.ExtractUUID(p.ID),
				Properties: struct {
					PrefixLength int "json:\"prefixLength,omitempty\" validate:\"omitempty,min=25,max=28\" default:\"27\""
				}{
					PrefixLength: 27,
				},
			}

			// Enable network services
			_, err := cc.c.Do(
				ctx,
				ep,
				cav.SetBody(requestBody),
			)
			if err != nil {
				// Check if the error is due to a service already enabled
				// The API returns a 400 Bad Request with a specific message in this case
				if !strings.Contains(err.Error(), "subnet not fully consumed") {
					return nil, fmt.Errorf("error enabling network services: %w", err)
				}
			}
			return nil, nil // No response is expected for this command.
		},
	})

	// * DisableCloudavenueServices
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Resource:           "CloudavenueServices",
		Verb:               "Disable",
		ShortDocumentation: "Disable Cloud Avenue services on an EdgeGateway.",
		LongDocumentation:  "Disable Cloud Avenue services on an EdgeGateway. Cloudavenue services is a network setting that allows the EdgeGateway to connect to the mutualized Cloud Avenue services (DNS, DHCP, etc.).",
		AutoGenerate:       true,

		ParamsType: types.ParamsEdgeGateway{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsEdgeGateway)
			ep := endpoints.DisableCloudavenueServices()

			// Ensure the edge gateway exists and retrieve its services
			// The GetNetworkServices method retrieves the edge gateway ID from the name if not provided.
			nSvc, err := cc.GetServices(ctx, p)
			if err != nil {
				return nil, fmt.Errorf("failed to get network services: %w", err)
			}

			// Disable network services
			_, err = cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], nSvc.Services.ID),
			)
			if err != nil {
				return nil, fmt.Errorf("error disabling network services: %w", err)
			}
			return nil, nil // No response is expected for this command.
		},
	})
}
