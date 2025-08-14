/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcgroup

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

//go:generate command-generator -path vdcgroup_commands.go

func init() {
	// * VdcGroup
	cmds.Register(commands.Command{
		Namespace:         "VdcGroup",
		LongDocumentation: "This command allows you to manage VDC Groups.",
	})

	// * ListVdcGroup
	cmds.Register(commands.Command{
		Namespace:          "VdcGroup",
		Verb:               "List",
		ShortDocumentation: "List Vdc Groups",
		LongDocumentation:  "List all Virtual Data Center Groups (Vdc Groups) available in your organization. If no filters are applied, it returns all Vdc Groups.",

		ParamsType: types.ParamsListVdcGroup{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the Vdc Group to filter by",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdcGroup"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the Vdc Group to filter by",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
				},
			},
		},
		ModelType: types.ModelListVdcGroup{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsListVdcGroup)
			ep := endpoints.ListVdcGroup()

			logger := cc.logger.WithGroup("ListVdcGroup")

			qP := ""
			if p.Name != "" {
				qP = fmt.Sprintf("name==%s", p.Name)
			}
			if p.ID != "" {
				qP = fmt.Sprintf("id==%s", p.ID)
			}

			resp, err := cc.c.Do(
				ctx,
				ep,
				cav.WithQueryParam(ep.QueryParams[0], qP),
			)
			if err != nil {
				logger.Error("Failed to get Vdc", "error", err)
				return nil, err
			}

			return resp.Result().(*itypes.ApiResponseListVdcGroup).ToModel(), nil
		},
		AutoGenerate: true,
	})

	// * GetVdcGroup
	cmds.Register(commands.Command{
		Namespace:          "VdcGroup",
		Verb:               "Get",
		ShortDocumentation: "Get a Vdc Group",
		LongDocumentation:  "Retrieve detailed information about a specific Vdc Group.",

		ParamsType: types.ParamsGetVdcGroup{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the Vdc Group to filter by",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdcGroup"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the Vdc Group to filter by",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		ModelType: types.ModelGetVdcGroup{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsGetVdcGroup)

			vdcgroups, err := cc.ListVdcGroup(ctx, types.ParamsListVdcGroup{
				ID:   p.ID,
				Name: p.Name,
			})
			if err != nil {
				cc.logger.Error("Failed to list Vdc Groups", "error", err)
				return nil, err
			}

			if vdcgroups == nil || len(vdcgroups.VdcGroups) == 0 {
				return nil, fmt.Errorf("Vdc Group not found")
			}

			if len(vdcgroups.VdcGroups) > 1 {
				return nil, fmt.Errorf("multiple Vdc Groups found")
			}

			return &vdcgroups.VdcGroups[0], nil
		},
		AutoGenerate: true,
	})

	// * CreateVdcGroup
	cmds.Register(commands.Command{
		Namespace:          "VdcGroup",
		Verb:               "Create",
		ShortDocumentation: "Create a Vdc Group",
		LongDocumentation:  "Create a new Virtual Data Center Group (Vdc Group) in your organization.",

		ParamsType: types.ParamsCreateVdcGroup{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the Vdc Group",
				Example:     "my-vdc-group",
				Required:    true,
			},
			commands.ParamsSpec{
				Name:        "description",
				Description: "Description of the Vdc Group",
				Required:    false,
			},
			commands.ParamsSpec{
				Name:        "vdcs.{index}.id",
				Description: "ID of the Vdc to associate with the Vdc Group",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdcs.{index}.name"),
					commands.ValidatorOmitempty(),
				},
			},
			commands.ParamsSpec{
				Name:        "vdcs.{index}.name",
				Description: "Name of the Vdc to associate with the Vdc Group",
				Required:    false,
				Example:     "my-vdc",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdcs.{index}.id"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		ModelType: types.ModelGetVdcGroup{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsCreateVdcGroup)

			epList := endpoints.ListVdcGroup()

			// List existing Vdc Groups to fast fail if VdcGroup already exist and get in the response context the orgID and the siteID
			respList, err := cc.c.Do(
				ctx,
				epList,
				cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("name==%s", p.Name)),
			)
			if err != nil {
				cc.logger.Error("Failed to list Vdc Groups", "error", err)
				return nil, err
			}

			rL := respList.Result().(*itypes.ApiResponseListVdcGroup)

			if len(rL.Values) > 0 {
				return nil, fmt.Errorf("Vdc Group already exists")
			}

			cd := cav.GetExtraDataFromContext(respList.Request.Context())

			// Create new Vdc Group
			body := itypes.ApiRequestCreateVdcGroup{
				OrgID:               cd.OrganizationID,
				Name:                p.Name,
				Description:         p.Description,
				NetworkProviderType: "NSX_T",
				Type:                "LOCAL",
				Vdcs:                make([]itypes.ApiResponseVdcGroupParticipatingVdc, 0),
			}

			necessaryRequestVdcID := false
			for _, vdc := range p.Vdcs {
				if vdc.ID == "" {
					necessaryRequestVdcID = true
					break
				}
			}

			listVdc := make(map[string]string)
			if necessaryRequestVdcID {
				respListVdc, err := cc.c.Do(ctx, endpoints.ListVdc())
				if err != nil {
					cc.logger.Error("Failed to list Vdcs", "error", err)
					return nil, err
				}

				// For each vdc
				for _, vdc := range respListVdc.Result().(*itypes.ApiResponseListVDC).Records {
					listVdc[vdc.Name] = vdc.ID
				}
			}

			for _, vdc := range p.Vdcs {
				x := itypes.ApiResponseVdcGroupParticipatingVdc{
					Vdc: itypes.ApiResponseVdcGroupParticipatingVdcRef{
						ID: func() string {
							if vdc.ID != "" {
								return vdc.ID
							}
							return listVdc[vdc.Name]
						}(),
						Name: vdc.Name,
					},
					Site: itypes.ApiResponseVdcGroupParticipatingSiteRef{
						ID: cd.SiteID,
					},
					FaultDomainTag:       "AZ01",
					NetworkProviderScope: "AZ01",
				}

				body.Vdcs = append(body.Vdcs, x)
			}

			// No response returned (Job)
			_, err = cc.c.Do(
				ctx,
				endpoints.CreateVdcGroup(),
				cav.SetBody(body),
			)
			if err != nil {
				cc.logger.Error("Failed to create Vdc Group", "error", err)
				return nil, err
			}

			return cc.GetVdcGroup(ctx, types.ParamsGetVdcGroup{
				Name: p.Name,
			})
		},
		AutoGenerate: true,
	})

	// * UpdateVdcGroup
	cmds.Register(commands.Command{
		Namespace:          "VdcGroup",
		Verb:               "Update",
		ShortDocumentation: "Update a Vdc Group",
		LongDocumentation:  "Update an existing Virtual Data Center Group (Vdc Group) in your organization.",

		ParamsType: types.ParamsUpdateVdcGroup{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the Vdc Group",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdcGroup"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the Vdc Group",
				Example:     "my-vdc-group",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
			commands.ParamsSpec{
				Name:        "description",
				Description: "Description of the Vdc Group",
				Required:    false,
			},
		},
		ModelType: types.ModelGetVdcGroup{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsUpdateVdcGroup)

			epList := endpoints.ListVdcGroup()

			param := cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("id==%s", p.ID))
			if p.ID == "" {
				param = cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("name==%s", p.Name))
			}

			// List existing Vdc Groups to fast fail if VdcGroup already exist and get in the response context the orgID and the siteID
			respList, err := cc.c.Do(
				ctx,
				epList,
				param,
			)
			if err != nil {
				cc.logger.Error("Failed to list Vdc Groups", "error", err)
				return nil, err
			}

			rL := respList.Result().(*itypes.ApiResponseListVdcGroup)

			if len(rL.Values) == 0 || len(rL.Values) > 1 {
				cc.logger.Error("Failed to get unique Vdc Group", "error", err)
				return nil, fmt.Errorf("failed to get unique Vdc Group")
			}

			p.ID = rL.Values[0].ID
			p.Name = rL.Values[0].Name

			body := itypes.ApiRequestUpdateVdcGroup{
				Id:                  p.ID,
				OrgID:               rL.Values[0].OrgID,
				Name:                p.Name,
				Description:         p.Description,
				Vdcs:                rL.Values[0].Vdcs,
				NetworkProviderType: "NSX_T",
				Type:                "LOCAL",
			}

			_, err = cc.c.Do(
				ctx,
				endpoints.UpdateVdcGroup(),
				cav.WithPathParam(endpoints.UpdateVdcGroup().PathParams[0], p.ID),
				cav.SetBody(body),
			)
			if err != nil {
				cc.logger.Error("Failed to update Vdc Group", "error", err)
				return nil, err
			}

			return cc.GetVdcGroup(ctx, types.ParamsGetVdcGroup{
				ID:   p.ID,
				Name: p.Name,
			})
		},
		AutoGenerate: true,
	})

	// * DeleteVdcGroup
	cmds.Register(commands.Command{
		Namespace:          "VdcGroup",
		Verb:               "Delete",
		ShortDocumentation: "Delete a Vdc Group",
		LongDocumentation:  "Delete an existing Virtual Data Center Group (Vdc Group) from your organization.",

		ParamsType: types.ParamsDeleteVdcGroup{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the Vdc Group",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdcGroup"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the Vdc Group",
				Example:     "my-vdc-group",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
			commands.ParamsSpec{
				Name:        "force",
				Description: "Force delete the Vdc Group. Value `true` means to forcefully delete the object that contains other objects even if those objects are in a state that does not allow removal. The default is `false` therefore, objects are not removed if they are not in a state that normally allows removal. Force also implies recursive delete where other contained objects are removed. Errors may be ignored. Invalid value (not true or false) are ignored. The VDC contains in the Vdc Group are not deleted.",
				Required:    false,
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsDeleteVdcGroup)

			if p.ID == "" {
				vdcGroup, err := cc.GetVdcGroup(ctx, types.ParamsGetVdcGroup{
					Name: p.Name,
				})
				if err != nil {
					cc.logger.Error("Failed to get Vdc Group", "error", err)
					return nil, err
				}

				p.ID = vdcGroup.ID
			}

			_, err := cc.c.Do(
				ctx,
				endpoints.DeleteVdcGroup(),
				cav.WithPathParam(endpoints.DeleteVdcGroup().PathParams[0], p.ID),
				cav.WithQueryParam(endpoints.DeleteVdcGroup().QueryParams[0], fmt.Sprintf("%t", p.Force)),
			)
			return nil, err
		},
		AutoGenerate: true,
	})

	// * AddVdcToVdcGroup
	cmds.Register(commands.Command{
		Namespace:                  "VdcGroup",
		Resource:                   "Vdc",
		Verb:                       "Add",
		AutoGenerateCustomFuncName: "AddVdcToVdcGroup",
		ShortDocumentation:         "Add a Vdc to a Vdc Group",
		LongDocumentation:          "Add an existing Virtual Data Center (Vdc) to a Virtual Data Center Group (Vdc Group) in your organization.",

		ParamsType: types.ParamsAddVdcToVdcGroup{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the Vdc Group",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdcGroup"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the Vdc Group",
				Example:     "my-vdc-group",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
			commands.ParamsSpec{
				Name:        "vdcs.{index}.id",
				Description: "ID of the Vdc to add",
				Required:    true,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdcs.{index}.name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "vdcs.{index}.name",
				Description: "Name of the Vdc to add",
				Required:    true,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdcs.{index}.id"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsAddVdcToVdcGroup)

			epList := endpoints.ListVdcGroup()

			param := cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("id==%s", p.ID))
			if p.ID == "" {
				param = cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("name==%s", p.Name))
			}

			// List existing Vdc Groups to fast fail if VdcGroup already exist and get in the response context the orgID and the siteID
			respList, err := cc.c.Do(
				ctx,
				epList,
				param,
			)
			if err != nil {
				cc.logger.Error("Failed to list Vdc Groups", "error", err)
				return nil, err
			}

			cd := cav.GetExtraDataFromContext(respList.Request.Context())
			rL := respList.Result().(*itypes.ApiResponseListVdcGroup)

			if len(rL.Values) == 0 || len(rL.Values) > 1 {
				cc.logger.Error("Failed to get unique Vdc Group", "error", err)
				return nil, fmt.Errorf("failed to get unique Vdc Group")
			}

			p.ID = rL.Values[0].ID
			p.Name = rL.Values[0].Name

			body := itypes.ApiRequestUpdateVdcGroup{
				Id:                  p.ID,
				OrgID:               rL.Values[0].OrgID,
				Name:                p.Name,
				Description:         rL.Values[0].Description,
				Vdcs:                rL.Values[0].Vdcs,
				NetworkProviderType: "NSX_T",
				Type:                "LOCAL",
			}

			necessaryRequestVdcID := false
			for _, vdc := range p.Vdcs {
				if vdc.ID == "" {
					necessaryRequestVdcID = true
					break
				}
			}

			listVdc := make(map[string]string)
			if necessaryRequestVdcID {
				respListVdc, err := cc.c.Do(ctx, endpoints.ListVdc())
				if err != nil {
					cc.logger.Error("Failed to list Vdcs", "error", err)
					return nil, err
				}

				// For each vdc
				for _, vdc := range respListVdc.Result().(*itypes.ApiResponseListVDC).Records {
					listVdc[vdc.Name] = vdc.ID
				}
			}

			for _, vdc := range p.Vdcs {
				x := itypes.ApiResponseVdcGroupParticipatingVdc{
					Vdc: itypes.ApiResponseVdcGroupParticipatingVdcRef{
						ID: func() string {
							if vdc.ID != "" {
								return vdc.ID
							}
							return listVdc[vdc.Name]
						}(),
						Name: vdc.Name,
					},
					Site: itypes.ApiResponseVdcGroupParticipatingSiteRef{
						ID: cd.SiteID,
					},
					FaultDomainTag:       "AZ01",
					NetworkProviderScope: "AZ01",
				}

				body.Vdcs = append(body.Vdcs, x)
			}

			_, err = cc.c.Do(
				ctx,
				endpoints.UpdateVdcGroup(),
				cav.WithPathParam(endpoints.UpdateVdcGroup().PathParams[0], p.ID),
				cav.SetBody(body),
			)
			if err != nil {
				cc.logger.Error("Failed to update Vdc Group", "error", err)
				return nil, err
			}

			return nil, nil
		},
		AutoGenerate: true,
	})

	// * RemoveVdcToVdcGroup
	cmds.Register(commands.Command{
		Namespace:                  "VdcGroup",
		Resource:                   "Vdc",
		Verb:                       "Remove",
		AutoGenerateCustomFuncName: "RemoveVdcFromVdcGroup",
		ShortDocumentation:         "Remove one or more Vdc from a Vdc Group",
		LongDocumentation:          "Remove one or more Vdc from a Vdc Group. This action will disassociate the specified Vdc(s) from the Vdc Group.",

		ParamsType: types.ParamsRemoveVdcFromVdcGroup{},
		ParamsSpecs: []commands.ParamsSpec{
			{
				Name:        "id",
				Description: "ID of the VDC Group",
				Required:    true,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdcGroup"),
				},
			},
			{
				Name:        "name",
				Description: "Name of the VDC Group",
				Required:    true,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
			{
				Name:        "vdcs.{index}.id",
				Description: "ID of the Vdc to remove",
				Required:    true,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdcs.{index}.name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			{
				Name:        "vdcs.{index}.name",
				Description: "Name of the Vdc to remove",
				Required:    true,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("vdcs.{index}.id"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(types.ParamsRemoveVdcFromVdcGroup)

			epList := endpoints.ListVdcGroup()

			param := cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("id==%s", p.ID))
			if p.ID == "" {
				param = cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("name==%s", p.Name))
			}

			// List existing Vdc Groups to fast fail if VdcGroup already exist and get in the response context the orgID and the siteID
			respList, err := cc.c.Do(
				ctx,
				epList,
				param,
			)
			if err != nil {
				cc.logger.Error("Failed to list Vdc Groups", "error", err)
				return nil, err
			}

			rL := respList.Result().(*itypes.ApiResponseListVdcGroup)

			if len(rL.Values) == 0 || len(rL.Values) > 1 {
				cc.logger.Error("Failed to get unique Vdc Group", "error", err)
				return nil, fmt.Errorf("failed to get unique Vdc Group")
			}

			p.ID = rL.Values[0].ID
			p.Name = rL.Values[0].Name

			body := itypes.ApiRequestUpdateVdcGroup{
				Id:                  p.ID,
				OrgID:               rL.Values[0].OrgID,
				Name:                p.Name,
				Description:         rL.Values[0].Description,
				Vdcs:                rL.Values[0].Vdcs,
				NetworkProviderType: "NSX_T",
				Type:                "LOCAL",
			}

			for _, vdc := range p.Vdcs {
				// Remove each vdc from the body
				for i := len(body.Vdcs) - 1; i >= 0; i-- {
					if body.Vdcs[i].Vdc.ID == vdc.ID || body.Vdcs[i].Vdc.Name == vdc.Name {
						body.Vdcs = append(body.Vdcs[:i], body.Vdcs[i+1:]...)
					}
				}
			}

			_, err = cc.c.Do(
				ctx,
				endpoints.UpdateVdcGroup(),
				cav.WithPathParam(endpoints.UpdateVdcGroup().PathParams[0], p.ID),
				cav.SetBody(body),
			)
			if err != nil {
				cc.logger.Error("Failed to update Vdc Group", "error", err)
				return nil, err
			}

			return nil, nil
		},
		AutoGenerate: true,
	})
}
