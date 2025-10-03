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

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate command-generator -path publicip_commands.go

func init() {
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Resource:  "PublicIP",
	})

	// * Create
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Resource:  "PublicIP",
		Verb:      "Create",

		ShortDocumentation: "Create a new Public IP",
		LongDocumentation:  "This command allows you to create a new Public IP (IPv4) on the specified Edge Gateway.",
		AutoGenerate:       true,

		ModelType:  types.ModelEdgeGatewayPublicIP{},
		ParamsType: types.ParamsCreateEdgeGatewayPublicIP{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "edge_gateway_id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "edge_gateway_name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Example:     "tn01e02ocb0001234spt101",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
					commands.ValidatorResourceName("edgegateway"),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsCreateEdgeGatewayPublicIP)
			ep := endpoints.CreatePublicIp()
			logger := cc.logger.WithGroup("CreatePublicIP")

			// ID is required to request the API.
			if p.EdgeGatewayID == "" {
				var err error
				p.EdgeGatewayID, err = cc.retrieveEdgeGatewayIDByName(ctx, p.EdgeGatewayName)
				if err != nil {
					return nil, err
				}
			}

			// Job extractor to get the public IP name from the job response
			// This is used to retrieve the public IP after creation.
			// │ BODY         :
			// │ [
			// │    {
			// │       "actions": [
			// │          {
			// │             "details": "123.45.67.89",
			// │             "name": "reserve_ip for Org cav01ev01ocb0001234 for public ip",
			// │             "status": "DONE"
			// │          },
			var publicipCreated string
			ep.SetJobExtractorFunc(cav.ExtractorFunc(func(resp *resty.Response) {
				r := resp.Result().(*cav.CerberusJobAPIResponse)

				if len(*r) == 0 {
					logger.ErrorContext(ctx, "No job information returned")
					return
				}

				job := (*r)[0]
				for _, j := range job.Actions {
					if err := validators.New().Var(j.Details, "ip4_addr"); err == nil {
						logger.DebugContext(ctx, "Found public IP in job data", "ip", j.Details)
						publicipCreated = j.Details
						break
					}
				}
			}))

			edgeId, err := extractor.ExtractUUID(p.EdgeGatewayID)
			if err != nil {
				return nil, err
			}

			body := itypes.ApiRequestEdgegatewayPublicIP{
				NetworkType:   "internet",
				EdgeGatewayID: edgeId,
				Properties: itypes.ApiRequestEdgegatewayPublicIPProperties{
					Announced: true,
				},
			}

			_, err = cc.c.Do(
				ctx,
				ep,
				cav.SetBody(body),
			)
			if err != nil {
				return nil, fmt.Errorf("Failed to create public IP: %w", err)
			}

			return cc.GetPublicIP(ctx, types.ParamsGetEdgeGatewayPublicIP{
				EdgeGatewayID:   p.EdgeGatewayID,
				EdgeGatewayName: p.EdgeGatewayName,
				IP:              publicipCreated,
			})
		},
	})

	// * List
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Resource:  "PublicIP",
		Verb:      "List",

		ShortDocumentation: "List Public IPs",
		LongDocumentation:  "This command allows you to list all Public IPs in the Edge Gateway.",
		AutoGenerate:       true,

		ModelType:  types.ModelEdgeGatewayPublicIPs{},
		ParamsType: types.ParamsListEdgeGatewayPublicIP{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "edge_gateway_id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "edge_gateway_name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Example:     "tn01e02ocb0001234spt101",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
					commands.ValidatorResourceName("edgegateway"),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsListEdgeGatewayPublicIP)

			services, err := cc.GetServices(ctx, types.ParamsEdgeGateway{
				ID:   p.EdgeGatewayID,
				Name: p.EdgeGatewayName,
			})
			if err != nil {
				return nil, err
			}

			ips := &types.ModelEdgeGatewayPublicIPs{
				PublicIps: make([]types.ModelEdgeGatewayPublicIP, 0, len(services.PublicIP)),
			}

			for _, publicip := range services.PublicIP {
				ips.PublicIps = append(ips.PublicIps, types.ModelEdgeGatewayPublicIP{
					ID:              publicip.ID,
					EdgeGatewayID:   services.ID,
					EdgeGatewayName: services.Name,
					IP:              publicip.IP,
					Announced:       publicip.Announced,
				})
			}

			return ips, nil
		},
	})

	// * Get
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Resource:  "PublicIP",
		Verb:      "Get",

		ShortDocumentation: "Get a Public IP",
		LongDocumentation:  "This command allows you to retrieve information about a Public IP in the Edge Gateway.",
		AutoGenerate:       true,

		ModelType:  types.ModelEdgeGatewayPublicIP{},
		ParamsType: types.ParamsGetEdgeGatewayPublicIP{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "ip",
				Description: "The public IP address.",
				Required:    true,
				Example:     "195.25.13.4",
				Validators: []commands.Validator{
					commands.ValidatorIPV4(),
				},
			},
			commands.ParamsSpec{
				Name:        "edge_gateway_id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("edgegateway"),
				},
			},
			commands.ParamsSpec{
				Name:        "edge_gateway_name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Example:     "tn01e02ocb0001234spt101",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
					commands.ValidatorResourceName("edgegateway"),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsGetEdgeGatewayPublicIP)

			// ID is required to request the API.
			if p.EdgeGatewayID == "" {
				var err error
				p.EdgeGatewayID, err = cc.retrieveEdgeGatewayIDByName(ctx, p.EdgeGatewayName)
				if err != nil {
					return nil, err
				}
			}

			ep := endpoints.GetEdgeGatewayServices()
			resp, err := cc.c.Do(
				ctx,
				ep,
				cav.WithQueryParam(ep.QueryParams[0], p.EdgeGatewayID),   // Only for filtering the response
				cav.WithQueryParam(ep.QueryParams[1], p.EdgeGatewayName), // Only for filtering the response
				cav.WithQueryParam(ep.QueryParams[2], p.IP),              // Only for filtering the response
			)
			if err != nil {
				return nil, fmt.Errorf("error retrieving network services for edge gateway %s(%s): %w", p.EdgeGatewayName, p.EdgeGatewayID, err)
			}

			data := resp.Result().(*itypes.ApiResponseNetworkServices).ToModel(types.ParamsEdgeGateway{
				ID:   p.EdgeGatewayID,
				Name: p.EdgeGatewayName,
			})
			if data == nil || len(data.PublicIP) == 0 {
				return nil, fmt.Errorf("no public IPs found for edge gateway %s(%s)", p.EdgeGatewayName, p.EdgeGatewayID)
			}

			for _, publicip := range data.PublicIP {
				if publicip.IP == p.IP {
					return &types.ModelEdgeGatewayPublicIP{
						ID:              publicip.ID,
						EdgeGatewayID:   data.ID,
						EdgeGatewayName: data.Name,
						IP:              publicip.IP,
						Announced:       publicip.Announced,
					}, nil
				}
			}

			return nil, fmt.Errorf("public IP %s not found in edge gateway %s(%s)", p.IP, data.Name, data.ID)
		},
	})

	// * Delete
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Resource:  "PublicIP",
		Verb:      "Delete",

		ShortDocumentation: "Delete a Public IP",
		LongDocumentation:  "This command allows you to delete a Public IP in the Edge Gateway.",
		AutoGenerate:       true,

		ParamsType: types.ParamsDeleteEdgeGatewayPublicIP{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "ip",
				Description: "The public IP address.",
				Required:    true,
				Example:     "195.25.13.4",
				Validators: []commands.Validator{
					commands.ValidatorIPV4(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsDeleteEdgeGatewayPublicIP)
			ep := endpoints.DisableCloudavenueServices()

			ipId := fmt.Sprintf("ip-%s", strings.ReplaceAll(p.IP, ".", "-"))

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], ipId),
			)

			return nil, err
		},
	})
}
