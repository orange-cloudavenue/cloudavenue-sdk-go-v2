/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package draas

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/pspecs"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/validator"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

//go:generate command-generator -path draas_on_premise_commands.go

func init() {
	// * Draas OnPremise
	cmds.Register(commands.Command{
		Namespace: "Draas",
		Resource:  "OnPremiseIp",
		Verb:      "",
	})

	// * ListDraasOnPremise
	cmds.Register(commands.Command{
		Namespace:          "Draas",
		Resource:           "OnPremiseIp",
		Verb:               "List",
		ShortDocumentation: "List all OnPremise IPs allowed",
		LongDocumentation:  "List all OnPremise IPs allowed allowed for this organization's draas offer",
		AutoGenerate:       true,
		ModelType:          types.ModelListDraasOnPremise{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)

			ep := endpoints.ListDraasOnPremiseIp()

			resp, err := cc.c.Do(
				ctx,
				ep,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to list Draas OnPremiseIP: %w", err)
			}

			return resp.Result().(*itypes.ApiResponseListDraasOnPremise).ToModel(), nil
		},
	})

	// * AddDraasOnPremiseIp
	cmds.Register(commands.Command{
		Namespace:          "Draas",
		Resource:           "OnPremiseIp",
		Verb:               "Add",
		ShortDocumentation: "Add a new OnPremise IP",
		LongDocumentation:  "Add a new OnPremise IP (only IPV4) address to this organization's draas offer",
		AutoGenerate:       true,
		ParamsType:         types.ParamsAddDraasOnPremiseIP{},
		ParamsSpecs: pspecs.Params{
			&pspecs.String{
				Name:        "ip",
				Required:    true,
				Description: "Your OnPremise IP address to authorize to the DRaaS offer",
				Example:     "195.25.13.4",
				Validators: []validator.Validator{
					validator.ValidatorIPV4(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsAddDraasOnPremiseIP)
			ep := endpoints.AddDraasOnPremiseIp()

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.IP),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to add Draas OnPremiseIP: %w", err)
			}

			return nil, nil
		},
	})

	// * RemoveDraasOnPremiseIp
	cmds.Register(commands.Command{
		Namespace:          "Draas",
		Resource:           "OnPremiseIp",
		Verb:               "Remove",
		ShortDocumentation: "Remove an existing OnPremise IP",
		LongDocumentation:  "Remove an existing OnPremise IP address from this organization's draas offer",
		AutoGenerate:       true,
		ParamsType:         types.ParamsRemoveDraasOnPremiseIP{},
		ParamsSpecs: pspecs.Params{
			&pspecs.String{
				Name:        "ip",
				Required:    true,
				Description: "Your OnPremise IP address to remove from the DRaaS offer",
				Example:     "195.25.13.4",
				Validators: []validator.Validator{
					validator.ValidatorIPV4(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsRemoveDraasOnPremiseIP)
			ep := endpoints.RemoveDraasOnPremiseIp()

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.IP),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to remove Draas OnPremiseIP: %w", err)
			}

			return nil, nil
		},
	})
}
