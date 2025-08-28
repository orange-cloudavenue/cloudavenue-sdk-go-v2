/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package organization

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

//go:generate command-generator -path organization_commands.go

func init() {
	// * ORGANIZATION
	// This command is a high-level command that allows you to manage documentation for Orgs resources.
	cmds.Register(commands.Command{
		Namespace: "Organization",
	})

	// * GetOrganization
	cmds.Register(commands.Command{
		Namespace:          "Organization",
		Verb:               "Get",
		ShortDocumentation: "Get an organization information.",
		LongDocumentation:  "Retrieve detailed information about your organization.",
		AutoGenerate:       true,
		ModelType:          types.ModelGetOrganization{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)

			logger := cc.logger.WithGroup("GetOrganization")

			// ----------------------------------------------------------
			// Parallelization of infraAPI and VMware Cloud Director calls
			var (
				errG       errgroup.Group
				org        *types.ModelGetOrganization
				orgDetails *types.ModelGetOrganization
			)

			// infraAPI goroutine
			errG.Go(func() error {
				resp, err := cc.c.Do(ctx, endpoints.GetOrganization())
				if err != nil {
					return fmt.Errorf("error getting organization information: %w", err)
				}
				org = resp.Result().(*itypes.ApiResponseGetOrg).ToModel()
				return nil
			})

			// VMware Cloud Director goroutine
			errG.Go(func() error {
				resp, err := cc.c.Do(ctx, endpoints.GetOrganizationDetails())
				if err != nil {
					return fmt.Errorf("error getting organizations details: %w", err)
				}
				orgDetails = resp.Result().(*itypes.ApiResponseGetOrgs).ToModel()
				if nil == orgDetails {
					return fmt.Errorf("error: organization not found")
				}
				return nil
			})

			// Wait for both goroutines to finish and check for errors
			if err := errG.Wait(); err != nil {
				return nil, err
			}
			// ----------------------------------------------------------

			logger.DebugContext(ctx, "Successfully retrieved organization information")

			return &types.ModelGetOrganization{
				ID:                  orgDetails.ID,
				Name:                org.Name,
				DisplayName:         org.DisplayName,
				Description:         org.Description,
				IsEnabled:           org.IsEnabled,
				CustomerMail:        org.CustomerMail,
				InternetBillingMode: org.InternetBillingMode,
				Resources: types.ModelGetOrganizationResources{
					Vdc:       orgDetails.Resources.Vdc,
					Catalog:   orgDetails.Resources.Catalog,
					Vapp:      orgDetails.Resources.Vapp,
					RunningVM: orgDetails.Resources.RunningVM,
					User:      orgDetails.Resources.User,
					Disk:      orgDetails.Resources.Disk,
				},
			}, nil
		},
	})

	// * UpdateOrganization
	cmds.Register(commands.Command{
		Namespace:          "Organization",
		Verb:               "Update",
		ShortDocumentation: "Update an existing organization.",
		LongDocumentation:  "Update the details of an existing organization.",
		AutoGenerate:       true,
		ParamsType:         types.ParamsUpdateOrganization{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "full_name",
				Description: "The full name of the organization. Appears in the Cloud application as a human-readable name of the organization.",
				Example:     "My Organization to update",
			},
			commands.ParamsSpec{
				Name:        "description",
				Description: "A description for the organization.",
				Example:     "This is my organization description to update",
			},
			commands.ParamsSpec{
				Name:        "customer_mail",
				Description: "The contact email for the organization.",
				Example:     "user@mail.com",
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorEmail(),
				},
			},
			commands.ParamsSpec{
				Name:        "internet_billing_mode",
				Description: "The internet billing mode for the organization. More information about billing modes can be found in the [documentation](https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/internet-access/).",
				Example:     "PAYG",
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorOneOf("PAYG", "TRAFFIC_VOLUME"),
				},
			},
		},
		ModelType: types.ModelGetOrganization{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)

			// Type assert params to types.ParamsUpdateOrganization
			p := params.(types.ParamsUpdateOrganization)

			// Check at least one parameter is provided
			if p.FullName == "" && p.CustomerMail == "" && p.InternetBillingMode == "" && p.Description == nil {
				return nil, fmt.Errorf("no parameters provided for organization update")
			}

			logger := cc.logger.WithGroup("UpdateOrganization")

			// Get Organization to set common values from infraAPI
			orgDefault, err := cc.c.Do(ctx, endpoints.GetOrganization())
			if err != nil {
				return nil, fmt.Errorf("error retrieving current organization information: %w", err)
			}
			data := orgDefault.Result().(*itypes.ApiResponseGetOrg).ToModel()

			// Set request body with provided parameters or keep existing values
			reqBody := &itypes.ApiRequestUpdateOrg{
				FullName: func() string {
					if p.FullName != "" {
						return p.FullName
					}
					return data.DisplayName
				}(),
				Description: func() string {
					if p.Description != nil {
						return *p.Description
					}
					return data.Description
				}(),
				CustomerMail: func() string {
					if p.CustomerMail != "" {
						return p.CustomerMail
					}
					return data.CustomerMail
				}(),
				InternetBillingMode: func() string {
					if p.InternetBillingMode != "" {
						return p.InternetBillingMode
					}
					return data.InternetBillingMode
				}(),
			}

			// Update Organization
			_, err = cc.c.Do(
				ctx,
				endpoints.UpdateOrganization(),
				cav.SetBody(reqBody),
			)
			if err != nil {
				return nil, fmt.Errorf("error updating organization: %w", err)
			}

			logger.DebugContext(ctx, "Successfully initiated organization update")

			// return Get organization after update
			return cc.GetOrganization(ctx)
		},
	})
}
