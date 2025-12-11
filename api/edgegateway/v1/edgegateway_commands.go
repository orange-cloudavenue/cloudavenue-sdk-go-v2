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
	"slices"

	"golang.org/x/sync/errgroup"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/pspecs"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate command-generator -path edgegateway_commands.go

func init() {
	// * EdgeGateway
	// This command is a high-level command that allows you to manage documentation for the EdgeGateway resource.
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
	})

	// * GetEdgeGateway
	cmds.Register(commands.Command{
		Namespace: "EdgeGateway",
		Verb:      "Get",

		ShortDocumentation: "GetEdgeGateway retrieves an edge gateway",
		LongDocumentation:  "Get EdgeGateway performs a GET request to retrieve an edge gateway",

		ParamsType: types.ParamsEdgeGateway{},
		ParamsSpecs: pspecs.Params{
			&pspecs.String{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []validator.Validator{
					validator.ValidatorRequiredIfParamIsNull("name"),
					validator.ValidatorOmitempty(),
					validator.ValidatorURN("edgegateway"),
				},
			},
			&pspecs.String{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []validator.Validator{
					validator.ValidatorRequiredIfParamIsNull("id"),
					validator.ValidatorOmitempty(),
					validator.ValidatorResourceName("edgegateway"),
				},
			},
		},
		ModelType: types.ModelEdgeGateway{},

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsEdgeGateway)
			ep := endpoints.GetEdgeGateway()

			logger := cc.logger.WithGroup("GetEdgeGateway")

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
				cav.WithPathParam(ep.PathParams[0], p.ID),
			)
			if err != nil {
				logger.Error("Failed to get edge gateway", "error", err)
				return nil, err
			}

			return resp.Result().(*itypes.ApiResponseEdgegateway).ToModel(), nil
		},
		AutoGenerate: true,
	})

	// * ListEdgeGateway
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Verb:               "List",
		ShortDocumentation: "ListEdgeGateways retrieves a list of edge gateways",
		LongDocumentation:  "List EdgeGateways performs a GET request to retrieve a list of edge gateways",
		ModelType:          types.ModelEdgeGateways{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			ep := endpoints.ListEdgeGateway()

			logger := cc.logger.WithGroup("ListEdgeGateways")

			resp, err := cc.c.Do(
				ctx,
				ep,
			)
			if err != nil {
				logger.Error("Failed to list edge gateways", "error", err)
				return nil, err
			}

			return resp.Result().(*itypes.ApiResponseEdgegateways).ToModel(), nil
		},
		AutoGenerate: true,
	})

	// * CreateEdgeGateway
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Verb:               "Create",
		ShortDocumentation: "CreateEdgeGateway creates a new edge gateway",
		LongDocumentation:  "Create EdgeGateway performs a POST request to create a new edge gateway",
		AutoGenerate:       true,
		ParamsType:         types.ParamsCreateEdgeGateway{},
		ParamsSpecs: pspecs.Params{
			&pspecs.String{
				Name:        "owner_name",
				Description: "The name of the VDC or VDC Group that this edge gateway belongs to.",
				Example:     "my-vdc",
				Required:    true,
			},
			&pspecs.String{
				Name:        "t0_name",
				Description: "The name of the T0 router that this edge gateway will be connected to. If not provided and only if one T0 router is available, the first T0 router will be used.",
				Required:    false,
				Example:     "prvrf01eocb0001234allsp01",
				Validators: []validator.Validator{
					validator.ValidatorOmitempty(),
					validator.ValidatorResourceName("t0"),
				},
			},
			&pspecs.Int{
				Name:        "bandwidth",
				Description: "The bandwidth limit in Mbps for the edge gateway. If t0 is SHARED, it must be one of the available values for the T0 router (Default value:  5Mbps). If t0 is DEDICATED, unlimited bandwidth is allowed.",
				Required:    false,
				Validators: []validator.Validator{
					validator.ValidatorOmitempty(),
				},
			},
		},
		ModelType: types.ModelEdgeGateway{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsCreateEdgeGateway)
			ep := endpoints.CreateEdgeGateway()

			logger := cc.logger.WithGroup("CreateEdgeGateway")

			// Here determine if the owner_name is a VDC or VDCGroup

			var (
				vdcs      *itypes.ApiResponseListVDC
				vdcGroups *itypes.ApiResponseListVdcGroup

				t0s *types.ModelT0s

				errG = errgroup.Group{}
			)

			epListVdc := endpoints.ListVdc()
			epListVdcGroup := endpoints.ListVdcGroup()

			// * Call in parallel all API
			errG.Go(func() error {
				resp, err := cc.c.Do(
					ctx,
					epListVdcGroup,
					cav.WithQueryParam(epListVdcGroup.QueryParams[0], "name=="+p.OwnerName),
				)
				if err != nil {
					return fmt.Errorf("Failed to list VDC Groups: %w", err)
				}

				vdcGroups = resp.Result().(*itypes.ApiResponseListVdcGroup)
				return nil
			})

			errG.Go(func() error {
				resp, err := cc.c.Do(
					ctx,
					epListVdc,
					cav.WithQueryParam(epListVdc.QueryParams[0], "name=="+p.OwnerName),
				)
				if err != nil {
					return fmt.Errorf("Failed to list VDCs: %w", err)
				}

				vdcs = resp.Result().(*itypes.ApiResponseListVDC)
				return nil
			})

			errG.Go(func() error {
				var err error
				t0s, err = cc.ListT0(ctx)
				if err != nil {
					return fmt.Errorf("Failed to list T0 routers: %w", err)
				}

				return nil
			})

			if err := errG.Wait(); err != nil {
				return nil, err
			}

			var t0 types.ModelT0

			// * T0s
			if t0s.Count == 0 {
				logger.Error("No T0 routers available to connect the edge gateway")
				return nil, errors.New("No T0 routers available to connect the edge gateway")
			}

			// If T0Name is not provided, use the first available T0 router.
			if p.T0Name == "" {
				if t0s.Count > 1 {
					logger.Warn("Multiple T0 routers found, using the first one", "count", t0s.Count)
					return nil, errors.New("Multiple T0 routers found, please specify T0Name")
				}
				t0 = t0s.T0s[0]
			} else {
				// Find the T0 router by name
				for _, t0Model := range t0s.T0s {
					if t0Model.Name == p.T0Name {
						t0 = t0Model
						break
					}
				}
				if t0.Name == "" {
					logger.Error("T0 router not found", "t0Name", p.T0Name)
					return nil, errors.New("T0 router not found: " + p.T0Name)
				}
			}

			if len(t0.EdgeGateways) >= t0.MaxEdgeGateways {
				return nil, errors.New("Maximum number of edge gateways reached for T0: " + t0.Name)
			}

			// If the T0 is SHARED, validate the bandwidth.
			if !t0.Bandwidth.AllowUnlimited {
				if p.Bandwidth <= 0 {
					p.Bandwidth = 5 // Default bandwidth if not provided
				}
				if !slices.Contains(t0.Bandwidth.AllowedBandwidthValues, p.Bandwidth) {
					logger.Error("Invalid bandwidth value for SHARED T0", "bandwidth", p.Bandwidth, "allowedValues", t0.Bandwidth.AllowedBandwidthValues, "remaining", t0.Bandwidth.Remaining)
					return nil, errors.New("Invalid bandwidth value for SHARED T0")
				}
			}

			// Prepare the request body
			reqBody := itypes.ApiRequestEdgeGateway{
				T0Name: t0.Name,
			}

			// * OwnerName (VDC Or VDCGroup)

			if (vdcs == nil || len(vdcs.Records) == 0) && (vdcGroups == nil || len(vdcGroups.Values) == 0) {
				return nil, errors.New("No VDCs or VDC Groups found for owner: " + p.OwnerName)
			}

			if (vdcs != nil && len(vdcs.Records) >= 1) && (vdcGroups != nil && len(vdcGroups.Values) >= 1) {
				return nil, errors.New("Both VDCs and VDC Groups found for owner: " + p.OwnerName)
			}

			var ownerName, ownerType string

			switch {
			case vdcs != nil && len(vdcs.Records) == 1:
				// Single VDC found
				ownerName, ownerType = vdcs.Records[0].Name, "vdc"
			case vdcGroups != nil && len(vdcGroups.Values) == 1:
				// Single VDC Group found
				ownerName, ownerType = vdcGroups.Values[0].Name, "vdcgroup"
			default:
				return nil, errors.New("Ambiguous owner: " + p.OwnerName)
			}

			// Job extractor to get the edge gateway name from the job response
			// This is used to retrieve the edge gateway after creation.
			var edgeGatewayCreated string
			ep.SetJobExtractorFunc(cav.ExtractorFunc(func(resp *resty.Response) {
				r := resp.Result().(*cav.CerberusJobAPIResponse)

				if len(*r) == 0 {
					logger.Error("No job information returned")
					return
				}

				job := (*r)[0]
				for _, j := range job.Actions {
					if err := validators.New().Var(j.Details, "resource_name=edgegateway"); err == nil {
						edgeGatewayCreated = j.Details
						break
					}
				}
			}))

			// * Create the edge gateway
			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], ownerType),
				cav.WithPathParam(ep.PathParams[1], ownerName),
				cav.SetBody(reqBody),
			)
			if err != nil {
				logger.Error("Failed to create edge gateway", "error", err)
				return nil, err
			}

			// Get the edge gateway created by name
			edgeCreated, err := cc.GetEdgeGateway(ctx, types.ParamsEdgeGateway{
				Name: edgeGatewayCreated,
			})
			if err != nil {
				logger.Error("Failed to retrieve created edge gateway", "error", err)
				return nil, err
			}

			// After creation, update the edge gateway with the bandwidth if provided and the value is upper than 5Mbps.
			if p.Bandwidth > 5 {
				// Prepare the update request body
				updateReqBody := itypes.ApiRequestBandwidth{
					Bandwidth: p.Bandwidth,
				}
				epBandwidth := endpoints.UpdateEdgeGatewayBandwidth()
				_, err := cc.c.Do(
					ctx,
					epBandwidth,
					cav.WithPathParam(epBandwidth.PathParams[0], edgeCreated.ID),
					cav.SetBody(updateReqBody),
				)
				if err != nil {
					logger.Error("Failed to update edge gateway bandwidth", "error", err)
					return nil, err
				}
			}

			return edgeCreated, nil
		},
	})

	// * DeleteEdgeGateway
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Verb:               "Delete",
		ShortDocumentation: "DeleteEdgeGateway deletes an edge gateway",
		LongDocumentation:  "Delete EdgeGateway performs a DELETE request to delete an edge gateway",
		AutoGenerate:       true,
		ParamsType:         types.ParamsEdgeGateway{},
		ParamsSpecs: pspecs.Params{
			&pspecs.String{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []validator.Validator{
					validator.ValidatorRequiredIfParamIsNull("name"),
					validator.ValidatorOmitempty(),
					validator.ValidatorURN("edgegateway"),
				},
			},
			&pspecs.String{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []validator.Validator{
					validator.ValidatorRequiredIfParamIsNull("id"),
					validator.ValidatorOmitempty(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsEdgeGateway)
			ep := endpoints.DeleteEdgeGateway()

			logger := cc.logger.WithGroup("DeleteEdgeGateway")

			// ID is required to request the API.
			if p.ID == "" {
				var err error
				p.ID, err = cc.retrieveEdgeGatewayIDByName(ctx, p.Name)
				if err != nil {
					return nil, err
				}
			}

			if _, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.ID),
			); err != nil {
				logger.Error("Failed to delete edge gateway", "error", err)
				return nil, err
			}

			return nil, nil
		},
	})

	// * UpdateEdgeGateway
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Verb:               "Update",
		ShortDocumentation: "UpdateEdgeGateway updates an edge gateway",
		LongDocumentation:  "Update EdgeGateway performs a PUT request to update an edge gateway",
		AutoGenerate:       true,
		ParamsType:         types.ParamsUpdateEdgeGateway{},
		ParamsSpecs: pspecs.Params{
			&pspecs.String{
				Name:        "id",
				Description: "The unique identifier of the edge gateway.",
				Required:    false,
				Validators: []validator.Validator{
					validator.ValidatorRequiredIfParamIsNull("name"),
					validator.ValidatorOmitempty(),
					validator.ValidatorURN("edgegateway"),
				},
			},
			&pspecs.String{
				Name:        "name",
				Description: "The name of the edge gateway.",
				Required:    false,
				Validators: []validator.Validator{
					validator.ValidatorRequiredIfParamIsNull("id"),
					validator.ValidatorOmitempty(),
				},
			},
			&pspecs.Int{
				Name:        "bandwidth",
				Description: "The new bandwidth limit in Mbps for the edge gateway. If t0 is SHARED, it must be one of the available values for the T0 router (Default value: 5Mbps). If t0 is DEDICATED, unlimited bandwidth is allowed.",
				Required:    true,
			},
		},
		ModelType: types.ModelEdgeGateway{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsUpdateEdgeGateway)
			ep := endpoints.UpdateEdgeGatewayBandwidth()

			logger := cc.logger.WithGroup("UpdateEdgeGateway")

			// ID is required to request the API.
			if p.ID == "" {
				var err error
				p.ID, err = cc.retrieveEdgeGatewayIDByName(ctx, p.Name)
				if err != nil {
					return nil, err
				}
			}

			// Prepare the request body
			reqBody := itypes.ApiRequestBandwidth{
				Bandwidth: p.Bandwidth,
			}

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.ID),
				cav.SetBody(reqBody),
			)
			if err != nil {
				logger.Error("Failed to update edge gateway bandwidth", "error", err)
				return nil, err
			}

			return cc.GetEdgeGateway(ctx, types.ParamsEdgeGateway{
				ID: p.ID,
			})
		},
	})
}
