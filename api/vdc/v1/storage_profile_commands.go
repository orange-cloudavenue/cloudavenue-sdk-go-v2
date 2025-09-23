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
	"fmt"
	"slices"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

//go:generate command-generator -path storage_profile_commands.go

func init() { //nolint:gocyclo
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
		ShortDocumentation: "Retrieve VDC storage profiles",
		LongDocumentation:  "Retrieves a comprehensive list of storage profiles. When no filters are specified, all storage profiles across all VDCs are returned. Filtering options include storage profile ID/name and VDC ID/name. Filters can be combined (e.g., profile filter + VDC filter). When both ID and name are provided for the same resource, they must reference the same object to return results.",
		AutoGenerate:       true,
		ParamsType:         types.ParamsListStorageProfile{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "Unique identifier of the storage profile to retrieve",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdcstorageProfile"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_id",
				Description: "Unique identifier of the VDC containing the storage profiles",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_name",
				Description: "Name of the VDC containing the storage profiles",
				Required:    false,
				Example:     "my-vdc",
			},
			commands.ParamsSpec{
				Name:        "class",
				Description: "Storage class name of the profile to retrieve",
				Required:    false,
				Example:     "gold",
			},
		},
		ModelType: types.ModelListStorageProfiles{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsListStorageProfile)

			ep := endpoints.ListStorageProfile()

			// Set the filter query parameter based on the provided parameters
			// When Required, Validators and TransformFunc are used, use a temp variable to build the filter (ref: QueryParam)
			var value string

			var filters []string
			if p.ID != "" {
				filters = append(filters, "id=="+p.ID)
			}
			if p.Class != "" {
				filters = append(filters, "name=="+p.Class) // Corrected from 'name' to 'class'
			}
			if p.VdcID != "" {
				filters = append(filters, "vdc=="+p.VdcID)
			}
			if p.VdcName != "" {
				filters = append(filters, "vdcName=="+p.VdcName)
			}
			if len(filters) > 0 {
				// Join all filters with ';' to allow multiple filters in the query param
				value = filters[0]
				if len(filters) > 1 {
					value = ""
					for i, f := range filters {
						if i > 0 {
							value += ";"
						}
						value += f
					}
				}
			}

			// Execute the request with the query parameters
			resp, err := cc.c.Do(
				ctx,
				ep,
				cav.WithQueryParam(ep.QueryParams[0], value),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to list VDC Storage Profiles: %w", err)
			}

			return resp.Result().(*itypes.ApiResponseListStorageProfiles).ToModel(), nil
		},
	})

	// * AddStorageProfile
	cmds.Register(commands.Command{
		Namespace: "VDC",
		Resource:  "StorageProfile",
		Verb:      "Add",

		ShortDocumentation: "Create VDC storage profiles",
		LongDocumentation:  "Creates one or more storage profiles within a specified VDC. Each profile requires a storage class and capacity limit, with an optional default designation.",
		AutoGenerate:       true,
		ParamsType:         types.ParamsAddStorageProfile{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "vdc_id",
				Description: "Unique identifier of the target VDC for the new storage profile",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_name",
				Description: "Name of the target VDC for the new storage profile",
				Required:    false,
				Example:     "my-vdc",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_id"),
					commands.ValidatorOmitempty(),
				},
			},
			{
				Name:        "storage_profiles.{index}.class",
				Description: "Storage class for the profile. Supports predefined and dedicated storage classes (see rules for available options)",
				Required:    true,
				Example:     "gold",
			},
			{
				Name:        "storage_profiles.{index}.limit",
				Description: "Storage capacity limit in GiB (maximum amount of storage available to the VDC)",
				Required:    true,
				Example:     "500",
				Validators: []commands.Validator{
					commands.ValidatorBetween(100, 81920),
				},
			},
			{
				Name:        "storage_profiles.{index}.default",
				Description: "Designates this storage profile as the default for the VDC when no specific profile is specified",
				Required:    false,
				Example:     "false",
			},
		},
		ParamsRules: func() commands.ParamsRules {
			pR := make(commands.ParamsRules, 0)

			searchField := []string{"storage_profiles.{index}.class", "storage_profiles.{index}.limit"}

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

			vdc, err := cc.GetVDC(ctx, types.ParamsGetVDC{
				ID:   p.VdcID,
				Name: p.VdcName,
			})
			if err != nil {
				return nil, err
			}

			return struct {
				// VdcId is the unique identifier of the VDC to add the storage profile to.
				VdcId string
				// VdcName is the name of the VDC to add the storage profile to.
				VdcName string

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

			vdc, err := cc.GetVDC(ctx, types.ParamsGetVDC{
				ID:   p.VdcID,
				Name: p.VdcName,
			})
			if err != nil {
				return nil, err
			}

			apiR := itypes.ApiRequestUpdateVDC{
				VDC: itypes.ApiRequestUpdateVDCVDC{
					Name: vdc.Name,
				},
			}

			logger.DebugContext(ctx, "Adding storage profiles to VDC", "vdc_name", vdc.Name, "storage_profiles", p.StorageProfiles)

			for _, sp := range p.StorageProfiles {
				apiR.VDC.StorageProfiles = append(apiR.VDC.StorageProfiles, itypes.ApiRequestVDCStorageProfile{
					Class:   sp.Class,
					Limit:   sp.Limit,
					Default: sp.Default,
				})
			}

			_, err = cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], vdc.Name),
				cav.SetBody(apiR),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to add storage profiles to VDC %s: %w", vdc.Name, err)
			}

			return nil, nil
		},
	})

	// * DeleteStorageProfile
	cmds.Register(commands.Command{
		Namespace:          "VDC",
		Resource:           "StorageProfile",
		Verb:               "Delete",
		ShortDocumentation: "Remove VDC storage profiles",
		LongDocumentation:  "Removes a storage profile from the specified VDC. Deletion is restricted for default profiles, the last remaining profile, or profiles currently in use.",
		AutoGenerate:       true,
		ParamsType:         types.ParamsDeleteStorageProfile{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "vdc_id",
				Description: "Unique identifier of the VDC containing the storage profile to remove",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_name",
				Description: "Name of the VDC containing the storage profile to remove",
				Required:    false,
				Example:     "my-vdc",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_id"),
					commands.ValidatorOmitempty(),
				},
			},
			{
				Name:        "storage_profiles.{index}.class",
				Description: "Storage class identifier of the profile to remove",
				Required:    true,
				Example:     "gold",
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsDeleteStorageProfile)

			logger := cc.logger.WithGroup("DeleteStorageProfile")

			ep := endpoints.UpdateVdc()

			// Use FastFailure to ensure VDC ID or Name are provided and valid
			listSP, err := cc.ListStorageProfile(ctx, types.ParamsListStorageProfile{
				VdcName: p.VdcName,
				VdcID:   p.VdcID,
			})
			if err != nil {
				return nil, err
			}
			if len(listSP.VDCS) == 0 {
				return nil, errors.New("no VDC found with the provided ID or Name")
			}
			if len(listSP.VDCS) > 1 {
				return nil, errors.New("multiple VDCs found with the provided ID or Name, please specify a unique VDC")
			}

			// Set the VDC ID and Name from the list response
			vdc := listSP.VDCS[0]
			p.VdcID = vdc.ID
			p.VdcName = vdc.Name

			// Check if the storage profile to delete is the only one left
			if len(vdc.StorageProfiles) == 1 {
				return nil, errors.New("cannot delete storage profile, at least one storage profile must exist for a VDC")
			}

			// Create api VDC request to delete all storage profiles
			apiR := itypes.ApiRequestUpdateVDC{
				VDC: itypes.ApiRequestUpdateVDCVDC{
					Name: p.VdcName,
				},
			}

			// Initialize the storage profiles slice in the API request
			// This will be used to delete the storage profile
			// Note: Only the storage profiles to delete will be processed
			// We will set the Limit to 0 to delete the storage profile
			// and Default to false.
			// As we cannot delete a default storage profile
			// And we need to ensure that the storage profile is not the only one left in the VDC.
			apiR.VDC.StorageProfiles = make([]itypes.ApiRequestVDCStorageProfile, 0, len(p.StorageProfiles))
			for _, pSP := range p.StorageProfiles {
				found := false
				for _, sp := range vdc.StorageProfiles {
					if sp.Class == pSP.Class {
						found = true
						if sp.Default {
							return nil, errors.Newf("cannot delete the default storage profile %s from VDC %s", sp.Class, p.VdcName)
						}
						if sp.Used > 0 {
							return nil, errors.Newf("cannot delete a non-empty storage profile %s from VDC %s", sp.Class, p.VdcName)
						}
						logger.DebugContext(ctx, "Deleting storage profile from VDC", "vdc_name", p.VdcName, "storage_profile_class", sp.Class)
						apiR.VDC.StorageProfiles = append(apiR.VDC.StorageProfiles, itypes.ApiRequestVDCStorageProfile{
							Class:   pSP.Class,
							Limit:   0,
							Default: false,
						})
					}
				}
				if !found {
					return nil, errors.Newf("storage profile class %s not found in VDC %s", pSP.Class, p.VdcName)
				}
			}

			_, err = cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.VdcName),
				cav.SetBody(apiR),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to remove storage profile(s) from VDC %s: %w", p.VdcName, err)
			}
			logger.DebugContext(ctx, "Storage profile deleted successfully", "vdc_name", p.VdcName, "storage_profile_class", p.StorageProfiles[0].Class)
			return nil, nil
		},
	})

	// * UpdateStorageProfile
	cmds.Register(commands.Command{
		Namespace:          "VDC",
		Resource:           "StorageProfile",
		Verb:               "Update",
		ShortDocumentation: "Modify VDC storage profiles",
		LongDocumentation:  "Modifies one or more storage profiles within a VDC. Supported updates include capacity limits and default profile designation. Storage class names cannot be modified.",
		AutoGenerate:       true,
		ParamsType:         types.ParamsUpdateStorageProfile{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "vdc_id",
				Description: "Unique identifier of the VDC containing the storage profile to modify",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdc_name",
				Description: "Name of the VDC containing the storage profile to modify",
				Required:    false,
				Example:     "my-vdc",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdc_id"),
					commands.ValidatorOmitempty(),
				},
			},
			{
				Name:        "storage_profiles.{index}.class",
				Description: "Storage class identifier of the profile to modify (see Rules Table for valid classes)",
				Required:    true,
				Example:     "gold",
			},
			{
				Name:        "storage_profiles.{index}.default",
				Description: "Designates this storage profile as the default for the VDC (only one default profile allowed per VDC)",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
				},
				Example: "true",
			},
			{
				Name:        "storage_profiles.{index}.limit",
				Description: "Updated storage capacity limit in GiB (cannot be lower than current usage; decreasing limits may cause service interruption)",
				Required:    false,
				Example:     "500",
				Validators: []commands.Validator{
					commands.ValidatorBetween(100, 81920),
				},
			},
		},

		// No additional rules are applied here.
		// Explanation:
		// - During an update, storage classes are already known.
		// - Re-checking this information would only add complexity
		//   without any benefit (the limit is the same for all classes).
		// We decided to remove the rules to keep the code simpler and more readable.

		ParamsRules: nil,

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsUpdateStorageProfile)

			logger := cc.logger.WithGroup("UpdateStorageProfile")

			ep := endpoints.UpdateVdc()

			// Use FastFailure to ensure VDC ID or Name are provided and valid
			listSP, err := cc.ListStorageProfile(ctx, types.ParamsListStorageProfile{
				VdcName: p.VdcName,
				VdcID:   p.VdcID,
			})
			if err != nil {
				return nil, err
			}
			if len(listSP.VDCS) == 0 {
				return nil, errors.New("no VDC found with the provided ID or Name")
			}
			if len(listSP.VDCS) > 1 {
				return nil, errors.New("multiple VDCs found with the provided ID or Name, please specify a unique VDC")
			}

			// Set the VDC ID and Name from the list response
			vdc := listSP.VDCS[0]
			p.VdcID = vdc.ID
			p.VdcName = vdc.Name

			// Create api VDC request to update all storage profiles
			apiR := itypes.ApiRequestUpdateVDC{
				VDC: itypes.ApiRequestUpdateVDCVDC{
					Name: p.VdcName,
				},
			}

			// Initialize the storage profiles slice in the API request
			// This will be used to update the storage profile
			// Note: Only the storage profiles to update will be processed
			apiR.VDC.StorageProfiles = make([]itypes.ApiRequestVDCStorageProfile, 0, len(p.StorageProfiles))
			// Create a map of current storage profiles for easy lookup
			currentStorageProfiles := make(map[string]types.ModelListStorageProfile, len(vdc.StorageProfiles))
			for _, sp := range vdc.StorageProfiles {
				currentStorageProfiles[sp.Class] = sp
			}

			// Track if a default storage profile is set in the update request
			var defaultCount int
			for _, sp := range p.StorageProfiles {
				// Check if the storage profile to update exists in the current storage profiles
				if _, ok := currentStorageProfiles[sp.Class]; !ok {
					return nil, fmt.Errorf("storage profile class %s not found in VDC %s", sp.Class, p.VdcName)
				}

				// If limit is set, check if the new limit is not less than the current used storage
				if sp.Limit < currentStorageProfiles[sp.Class].Used && sp.Limit > 0 {
					return nil, fmt.Errorf("new limit for storage profile %s cannot be less than the current used (%d GiB)", sp.Class, currentStorageProfiles[sp.Class].Used)
				}

				// Set the storage profile in the API request
				apiR.VDC.StorageProfiles = append(apiR.VDC.StorageProfiles, itypes.ApiRequestVDCStorageProfile{
					Class: sp.Class,
					Limit: func() int {
						if sp.Limit > 0 {
							return sp.Limit
						}
						return currentStorageProfiles[sp.Class].Limit
					}(),
					Default: func() bool {
						// If sp.Default is set, use its value and increment the counter if it is true.
						if sp.Default != nil {
							if *sp.Default {
								defaultCount++
							}
							return *sp.Default
						}
						// If Default is not specified, use the current value
						if currentStorageProfiles[sp.Class].Default {
							return true
						}
						return false
					}(),
				})
			}

			if defaultCount > 1 {
				return nil, fmt.Errorf("multiple storage profiles have default=true, only one is allowed")
			}

			_, err = cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.VdcName),
				cav.SetBody(apiR),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to update storage profile(s) in VDC %s: %w", p.VdcName, err)
			}
			logger.DebugContext(ctx, "Storage profile updated successfully", "vdc_name", p.VdcName)

			x, err := cc.ListStorageProfile(ctx, types.ParamsListStorageProfile{
				VdcID:   p.VdcID,
				VdcName: p.VdcName,
			})

			// Return the updated storage profiles for the VDC,
			// Only one VDC as we used VDC ID or Name in the request.
			return &x.VDCS[0], err
		},
		ModelType: types.ModelListStorageProfilesVDC{},
	})
}
