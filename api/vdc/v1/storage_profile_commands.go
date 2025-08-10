package vdc

import (
	"context"
	"slices"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
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
		Namespace: "VDC",
		Resource:  "StorageProfile",
		Verb:      "List",

		ShortDocumentation: "List VDC Storage Profiles",
		LongDocumentation:  "List of storage profiles available in a specific VDC.",

		ParamsType: ParamsListStorageProfiles{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the VDC to get",
				Required:    true,
				Validators: []commands.Validator{
					commands.ValidatorURN("vdc"),
				},
			},
		},
		ModelType: ModelListStorageProfiles{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsListStorageProfiles)

			logger := cc.logger.WithGroup("ListStorageProfiles")

			ep := endpoints.ListStorageProfiles()

			resp, err := cc.c.Do(
				ctx,
				ep,
				cav.WithQueryParam(ep.QueryParams[0], p.ID),
			)
			if err != nil {
				logger.Error("Failed to list VDC Storage Profiles", "error", err)
				return nil, err
			}

			return resp.Result().(*apiResponseListStorageProfiles).toModel(), nil
		},
		AutoGenerate: true,
	})

	// * AddStorageProfile
	cmds.Register(commands.Command{
		Namespace: "VDC",
		Resource:  "StorageProfile",
		Verb:      "Add",

		ShortDocumentation: "Add a new VDC Storage Profile",
		LongDocumentation:  "Add one or more storage profiles to a specific VDC.",
		ModelType:          cav.Job{},
		ParamsType:         ParamsAddStorageProfile{},
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
					commands.ValidatorRequiredIfParamIsNull("vdc_id"),
					commands.ValidatorOmitempty(),
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
			p := paramsIn.(ParamsAddStorageProfile)

			vdc, err := cc.GetVDC(ctx, ParamsGetVDC{
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

				StorageProfiles    []ParamsCreateVDCStorageProfile
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
			p := params.(ParamsAddStorageProfile)

			logger := cc.logger.WithGroup("AddStorageProfile")

			ep := endpoints.UpdateVdc()

			apiR := apiRequestUpdateVDC{
				VDC: apiRequestUpdateVDCVDC{
					Name: p.VdcName,
				},
			}

			logger.Debug("Adding storage profiles to VDC", "vdc_name", p.VdcName, "storage_profiles", p.StorageProfiles)

			for _, sp := range p.StorageProfiles {
				apiR.VDC.StorageProfiles = append(apiR.VDC.StorageProfiles, apiRequestVDCStorageProfile{
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
				logger.Error("Failed to update VDC", "error", err)
				return nil, err
			}

			return nil, nil
		},
		AutoGenerate: true,
	})
}
