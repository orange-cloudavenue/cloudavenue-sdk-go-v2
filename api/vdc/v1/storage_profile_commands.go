/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"context"
	"slices"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

//go:generate command-generator -path storage_profile_commands.go

func init() {
	// * StorageProfiles
	cmds.Register(commands.Command{
		Namespace: "VDC",
		Resource:  "StorageProfile",
		Verb:      "",
	})

	// * ListStorageProfiles
	cmds.Register(commands.Command{
		Namespace:          "VDC",
		Resource:           "StorageProfile",
		Verb:               "List",
		ShortDocumentation: "List VDC Storage Profiles",
		LongDocumentation:  "List of storage profiles in All VDC.",
		AutoGenerate:       true,
		ParamsType:         ParamsListStorageProfile{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the storage profile to list",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdcstorageProfile"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_id",
				Description: "ID of the VDC to get the storage profile from",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_name",
				Description: "Name of the VDC to get the storage profile from",
				Required:    false,
				Example:     "my-vdc",
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the storage profile to list",
				Required:    false,
				Example:     "gold",
			},
		},
		ModelType: types.ModelListStorageProfiles{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsListStorageProfile)

			logger := cc.logger.WithGroup("ListStorageProfile")

			ep := endpoints.ListStorageProfile()

			// Set the filter query parameter based on the provided parameters
			// When Required, Validators and TransformFunc are used, use a temp variable to build the filter (ref: QueryParam)
			var value string
			if p.ID != "" {
				// If ID is provided, we filter by ID
				value = "id==" + p.ID
			}
			if p.Name != "" {
				// If Name is provided, we filter by Name
				value = "name==" + p.Name
			}
			if p.VdcID != "" {
				// If VdcId is provided, we filter by VDC ID
				value = "vdc==" + p.VdcID
			}
			if p.VdcName != "" {
				// If VdcName is provided, we filter by VDC Name
				value = "vdcName==" + p.VdcName
			}
			logger.DebugContext(ctx, "Listing storage profiles", "params", p)

			// Execute the request with the query parameters
			resp, err := cc.c.Do(
				ctx,
				ep,
				cav.WithQueryParam(ep.QueryParams[0], value),
			)
			if err != nil {
				logger.ErrorContext(ctx, "Failed to list VDC Storage Profiles", "error", err)
				return nil, err
			}

			return resp.Result().(*itypes.ApiResponseListStorageProfiles).ToModel(), nil
		},
	})

	// * AddStorageProfile
	cmds.Register(commands.Command{
		Namespace: "VDC",
		Resource:  "StorageProfile",
		Verb:      "Add",

		ShortDocumentation: "Add a new VDC Storage Profile",
		LongDocumentation:  "Add one or more storage profiles to a specific VDC.",
		ModelType:          cav.Job{},
		ParamsType:         types.ParamsAddStorageProfile{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "vdc_id",
				Description: "ID of the VDC to add the storage profile to",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_name",
				Description: "Name of the VDC to add the storage profile to",
				Required:    false,
				Example:     "my-vdc",
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorRequiredIfParamIsNull("vdc_id"),
				},
			},
			{
				Name:        "storage_profiles.{index}.class",
				Description: "Class of the storage profile to create. Predefined classes or dedicated storage classes can be used. See rules for more information.",
				Required:    true,
				Example:     "gold",
			},
			{
				Name:        "storage_profiles.{index}.limit",
				Description: "Limit of the storage profile to create. This is the maximum amount of storage that can be used by the VDC. This is in GiB.",
				Required:    true,
				Example:     "500",
				Validators: []commands.Validator{
					commands.ValidatorBetween(100, 81920),
				},
			},
			{
				Name:        "storage_profiles.{index}.default",
				Description: "Default storage profile to create. This will be used if no specific profile is provided.",
				Required:    false,
				Example:     "false",
			},
		},
		ParamsRules: func() commands.ParamsRules {
			pR := make(commands.ParamsRules, 0)

			searchField := []string{"storage_profiles.{index}.class"}

			for _, spec := range vdcRules {
				if slices.Contains(searchField, spec.Target) {
					pR = append(pR, spec)
				}
			}
			return pR
		}(),

		// PreRulesRunnerFunc is called before the main command for inject DisponibilityClass in a rules validation
		PreRulesRunnerFunc: func(ctx context.Context, cmd *commands.Command, client, paramsIn any) (paramsOut any, err error) {
			cc := client.(*Client)
			p := paramsIn.(types.ParamsAddStorageProfile)

			vdc, err := cc.GetVDC(ctx, ParamsGetVDC{
				ID:   p.VdcId,
				Name: p.VdcName,
			})
			if err != nil {
				return nil, err
			}

			return struct {
				// VdcId is the unique identifier of the VDC to add the storage profile to.
				VdcId string //nolint:revive
				// VdcName is the name of the VDC to add the storage profile to.
				VdcName string //nolint:revive

				StorageProfiles    []types.ParamsCreateVDCStorageProfile
				DisponibilityClass string
			}{
				VdcId:              vdc.ID,
				VdcName:            vdc.Name,
				StorageProfiles:    p.StorageProfiles,
				DisponibilityClass: vdc.Properties.DisponibilityClass,
			}, nil
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsAddStorageProfile)

			logger := cc.logger.WithGroup("AddStorageProfile")

			ep := endpoints.UpdateVdc()

			apiR := itypes.ApiRequestUpdateVDC{
				VDC: itypes.ApiRequestUpdateVDCVDC{
					Name: p.VdcName,
				},
			}

			logger.DebugContext(ctx, "Adding storage profiles to VDC", "vdc_name", p.VdcName, "storage_profiles", p.StorageProfiles)

			for _, sp := range p.StorageProfiles {
				apiR.VDC.StorageProfiles = append(apiR.VDC.StorageProfiles, itypes.ApiRequestVDCStorageProfile{
					Class:   sp.Class,
					Limit:   sp.Limit,
					Default: sp.Default,
				})
			}

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.VdcName),
				cav.SetBody(apiR),
			)
			if err != nil {
				logger.ErrorContext(ctx, "Failed to update VDC", "error", err)
				return nil, err
			}

			return nil, nil
		},
		AutoGenerate: true,
	})

	// * DeleteStorageProfile
	cmds.Register(commands.Command{
		Namespace:          "VDC",
		Resource:           "StorageProfile",
		Verb:               "Delete",
		ShortDocumentation: "Delete a VDC Storage Profile",
		LongDocumentation:  "Delete a storage profile from a specific VDC. This will remove the storage profile from the VDC and all associated resources.",
		AutoGenerate:       true,
		ModelType:          cav.Job{},
		ParamsType:         ParamsDeleteStorageProfile{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "vdc_id",
				Description: "ID of the VDC to delete the storage profile from",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_name",
				Description: "Name of the VDC to delete the storage profile from",
				Required:    false,
				Example:     "my-vdc",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_id"),
					commands.ValidatorOmitempty(),
				},
			},
			{
				Name:        "storage_profile.class",
				Description: "Class of the storage profile to delete. This is the unique identifier of the storage profile to delete.",
				Required:    true,
				Example:     "gold",
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsDeleteStorageProfile)

			logger := cc.logger.WithGroup("DeleteStorageProfile")

			ep := endpoints.UpdateVdc()

			apiR := apiRequestUpdateVDC{
				VDC: apiRequestUpdateVDCVDC{
					Name: p.VdcName,
				},
			}

			logger.DebugContext(ctx, "Deleting storage profile from VDC", "vdc_name", p.VdcName, "storage_profile_class", p.StorageProfile[0].Class)

			apiR.VDC.StorageProfiles = []apiRequestVDCStorageProfile{
				{
					Class: p.StorageProfile[0].Class,
					Limit: 0, // Setting limit to 0 to delete the storage profile
				},
			}

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.VdcName),
				cav.SetBody(apiR),
			)
			if err != nil {
				logger.ErrorContext(ctx, "Failed to update VDC", "error", err)
				return nil, err
			}
			logger.InfoContext(ctx, "Storage profile deleted successfully", "vdc_name", p.VdcName, "storage_profile_class", p.StorageProfile[0].Class)
			return nil, nil
		},
	})
}
