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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

//go:generate command-generator -path bandwidth_commands.go

func init() {
	// Register the GetEdgeGatewayBandwidth command
	cmds.Register(commands.Command{
		Namespace:          "EdgeGateway",
		Resource:           "Bandwidth",
		Verb:               "Get",
		ShortDocumentation: "Get the bandwidth of an edge gateway.",
		LongDocumentation:  "Get the bandwidth of an edge gateway. This command retrieves the bandwidth information for a specific edge gateway.",
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

		ModelType: types.ModelEdgeGatewayBandwidth{},

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsEdgeGateway)

			pT0 := types.ParamsGetT0{
				EdgeGatewayID:   p.ID,
				EdgeGatewayName: p.Name,
			}

			t0, err := cc.GetT0(ctx, pT0)
			if err != nil {
				return nil, err
			}

			var edgeGateway *types.ModelT0EdgeGateway
			for _, eg := range t0.EdgeGateways {
				if eg.ID == p.ID || eg.Name == p.Name {
					edgeGateway = &eg
					break
				}
			}

			// edgeGateway is never nil because GetT0FromEdgeGateway ensures the edge gateway exists

			// Transform ModelT0EdgeGateway to types.ModelEdgeGatewayBandwidth
			// Create a new types.ModelEdgeGatewayBandwidth from ModelT0EdgeGateway
			bandwidth := &types.ModelEdgeGatewayBandwidth{
				ID:                     edgeGateway.ID,
				Name:                   edgeGateway.Name,
				Bandwidth:              edgeGateway.Bandwidth,
				AllowedBandwidthValues: edgeGateway.AllowedBandwidthValues,
			}

			return bandwidth, nil
		},
	})
}
