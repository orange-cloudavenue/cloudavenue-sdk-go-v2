package vdc

import (
	"context"
	"fmt"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

//go:generate command-generator -path vdc_commands.go

func init() { //nolint:gocyclo
	// * VDC
	// This command is a high-level command that allows you to manage documentation for the VDC resource.
	cmds.Register(commands.Command{
		Namespace: "VDC",
	})

	// * ListVDC
	cmds.Register(commands.Command{
		Namespace: "VDC",
		Verb:      "List",

		ShortDocumentation: "List VDCs",
		LongDocumentation:  "List all Virtual Data Centers (VDCs) available in your organization. If no filters are applied, it returns all VDCs.",

		ParamsType: ParamsListVDC{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the VDC to filter by",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the VDC to filter by",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorOmitempty(),
				},
			},
		},
		ModelType: ModelListVDC{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsListVDC)
			ep := endpoints.ListVdc()

			logger := cc.logger.WithGroup("ListVDC")

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
				logger.Error("Failed to get VDC", "error", err)
				return nil, err
			}

			return resp.Result().(*apiResponseListVDC).toModel(), nil
		},
		AutoGenerate: true,
	})

	// * GetVDC
	cmds.Register(commands.Command{
		Namespace: "VDC",
		Verb:      "Get",

		ShortDocumentation: "Get VDC details",
		LongDocumentation:  "Retrieve detailed information about a specific Virtual Data Center (VDC) by its name.",

		ParamsType: ParamsGetVDC{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the VDC to get",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the VDC to get",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		ModelType: ModelGetVDC{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsGetVDC)

			logger := cc.logger.WithGroup("GetVDC")

			// First list all VDCs with filter to fast fail if no VDCs are found and retrieve the VDC ID
			epListVDC := endpoints.ListVdc()
			qP := ""
			if p.Name != "" {
				qP = fmt.Sprintf("name==%s", p.Name)
			}
			if p.ID != "" {
				qP = fmt.Sprintf("id==%s", p.ID)
			}

			resp, err := cc.c.Do(
				ctx,
				epListVDC,
				cav.WithQueryParam(epListVDC.QueryParams[0], qP),
			)
			if err != nil {
				logger.Error("Failed to list VDCs", "error", err)
				return nil, err
			}

			results := resp.Result().(*apiResponseListVDC)
			if len(results.Records) == 0 {
				logger.Warn("No VDCs found")
				return nil, fmt.Errorf("The VDC %s does not exist in your organization", p.Name)
			}

			vdc := results.Records[0]

			var (
				vdcMetadata *apiResponseGetVDCMetadatas
				model       ModelGetVDC
			)

			// GET VDC Details and Metadata in Parallel
			eg, egCtx := errgroup.WithContext(ctx)

			eg.Go(func() error {
				// Get VDC Metadata
				epGetVDCMetadata := endpoints.GetVdcMetadata()
				logger.Debug("Fetching VDC metadata", "vdcName", results.Records[0].Name, "vdcID", vdc.ID)
				vdcMetadataResp, err := cc.c.Do(
					egCtx,
					epGetVDCMetadata,
					cav.WithPathParam(epGetVDCMetadata.PathParams[0], results.Records[0].ID),
				)
				if err != nil {
					logger.Error("Failed to get VDC metadata", "error", err, "vdcName", results.Records[0].Name)
					return fmt.Errorf("failed to get VDC metadata for %s: %w", results.Records[0].Name, err)
				}

				vdcMetadata = vdcMetadataResp.Result().(*apiResponseGetVDCMetadatas)
				return nil
			})

			eg.Go(func() error {
				epGetVDC := endpoints.GetVdc()
				logger.Debug("Fetching VDC details", "vdcName", vdc.Name, "vdcID", vdc.ID)
				vdcResp, err := cc.c.Do(
					ctx,
					epGetVDC,
					cav.WithPathParam(epGetVDC.PathParams[0], vdc.ID),
				)
				if err != nil {
					logger.Error("Failed to get VDC details", "error", err, "vdcName", vdc.Name)
					return err
				}

				vdcDetails := vdcResp.Result().(*apiResponseGetVDC)
				model = vdcDetails.toModel()
				model.NumberOfDisks = vdc.NumberOfDisks
				model.NumberOfStorageProfiles = vdc.NumberOfStorageProfiles
				model.NumberOfVMS = vdc.NumberOfVMS
				model.NumberOfRunningVMS = vdc.NumberOfRunningVMS
				model.NumberOfVAPPS = vdc.NumberOfVAPPS

				return nil
			})

			if err := eg.Wait(); err != nil {
				logger.Error("Error while fetching VDC details or metadata", "error", err)
				return nil, fmt.Errorf("failed to get VDC details or metadata: %w", err)
			}

			// Populate model with metadata
			for _, metadata := range vdcMetadata.Metadatas {
				switch metadata.Name {
				case "vdcBillingModel":
					model.Properties.BillingModel = metadata.Value.Value
				case "vdcStorageBillingModel":
					model.Properties.StorageBillingModel = metadata.Value.Value
				case "vdcDisponibilityClass":
					model.Properties.DisponibilityClass = metadata.Value.Value
				case "vdcServiceClass":
					model.Properties.ServiceClass = metadata.Value.Value
				}
			}

			logger.Debug("Successfully retrieved VDC details", "vdcName", model.Name, "vdcID", model.ID)

			return &model, nil
		},
		AutoGenerate: true,
	})

	// * CreateVDC
	cmds.Register(commands.Command{
		Namespace: "VDC",
		Verb:      "Create",

		ShortDocumentation: "Create a new VDC",
		LongDocumentation:  "Create a new Virtual Data Center (VDC) with the specified parameters.",

		ParamsType: ParamsCreateVDC{},
		ParamsSpecs: commands.ParamsSpecs{
			{
				Name:        "name",
				Description: "Name of the VDC to create",
				Required:    true,
				Example:     "my-vdc",
			},
			{
				Name:        "description",
				Description: "Description of the VDC to create",
				Required:    false,
			},
			{
				Name:        "service_class",
				Description: "Service class of the VDC to create",
				Required:    true,
				Example:     "STD",
				Validators: []commands.Validator{
					commands.ValidatorOneOf("ECO", "STD", "HP", "VOIP"),
				},
			},
			{
				Name:        "disponibility_class",
				Description: "Disponibility class of the VDC to create",
				Required:    true,
				Example:     "ONE-ROOM",
				Validators: []commands.Validator{
					commands.ValidatorOneOf("ONE-ROOM", "DUAL-ROOM", "HA-DUAL-ROOM"),
				},
			},
			{
				Name:        "billing_model",
				Description: "Billing model of the VDC to create",
				Required:    true,
				Example:     "PAYG",
				Validators: []commands.Validator{
					commands.ValidatorOneOf("PAYG", "DRAAS", "RESERVED"),
				},
			},
			{
				Name:        "storage_billing_model",
				Description: "Storage billing model of the VDC to create",
				Required:    true,
				Example:     "PAYG",
				Validators: []commands.Validator{
					commands.ValidatorOneOf("PAYG", "RESERVED"),
				},
			},
			{
				Name:        "vcpu",
				Description: "Number of vCPUs to allocate to the VDC.",
				Required:    true,
				Example:     "50",
			},
			{
				Name:        "memory",
				Description: "Amount of memory (in GB) to allocate to the VDC.",
				Required:    true,
				Example:     "500",
			},
			{
				Name:        "storage_profiles.{index}.class",
				Description: "Class of the storage profile to create. Predefined classes or dedicated storage classes can be used. For predefined classes you have different properties (`_r1`, `_r2` and `_hm`) that can be used to define the storage profile.",
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
				Example:     "true",
			},
		},
		ParamsRules: vdcRules,
		ModelType:   ModelGetVDC{},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsCreateVDC)

			logger := cc.logger.WithGroup("CreateVDC")

			reqBody := apiRequestCreateVDC{
				VDC: apiRequestCreateVDCVDC{
					Name:                p.Name,
					Description:         p.Description,
					ServiceClass:        p.ServiceClass,
					DisponibilityClass:  p.DisponibilityClass,
					BillingModel:        p.BillingModel,
					StorageBillingModel: p.StorageBillingModel,
					VCPUInMhz:           serviceClassToCPUInMhz(p.ServiceClass),
					CPUAllocated:        serviceClassToCPUInMhz(p.ServiceClass) * p.Vcpu,
					MemoryAllocated:     p.Memory,
					StorageProfiles:     make([]apiRequestVDCStorageProfile, len(p.StorageProfiles)),
				},
			}

			for i, sp := range p.StorageProfiles {
				reqBody.VDC.StorageProfiles[i] = apiRequestVDCStorageProfile{
					Class:   sp.Class,
					Limit:   sp.Limit,
					Default: sp.Default,
				}
			}

			haveOneDefaultStorageProfile := false
			for _, sp := range reqBody.VDC.StorageProfiles {
				if sp.Default {
					haveOneDefaultStorageProfile = true
					break
				}
			}
			if !haveOneDefaultStorageProfile {
				return nil, fmt.Errorf("at least one storage profile must be marked as default")
			}

			ep := endpoints.CreateVdc()
			_, err := cc.c.Do(
				ctx,
				ep,
				cav.SetBody(reqBody),
			)
			if err != nil {
				logger.Error("Failed to create VDC", "error", err)
				return nil, err
			}

			resp, err := cc.GetVDC(ctx, ParamsGetVDC{
				Name: p.Name,
			})
			if err != nil {
				logger.Error("Failed to get VDC", "error", err)
				return nil, err
			}

			return resp, nil
		},
		AutoGenerate: true,
	})

	// * UpdateVDC
	cmds.Register(commands.Command{
		Namespace:          "VDC",
		Verb:               "Update",
		ShortDocumentation: "UpdateVDC updates an existing VDC",
		LongDocumentation:  "Update VDC performs a PUT request to update an existing VDC. Enter only the fields you want to update.",
		ParamsType:         ParamsUpdateVDC{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the VDC to get",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the VDC to get",
				Required:    false,
				Example:     "my-vdc",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
			commands.ParamsSpec{
				Name:        "description",
				Description: "The description of the VDC.",
				Required:    false,
			},
			commands.ParamsSpec{
				Name:        "vcpu",
				Description: "The number of virtual CPUs for the VDC.",
				Required:    false,
				Example:     "50",
			},
			commands.ParamsSpec{
				Name:        "memory",
				Description: "The amount of memory for the VDC.",
				Required:    false,
				Example:     "500",
			},
		},
		ParamsRules: func() commands.ParamsRules {
			pR := make(commands.ParamsRules, 0)

			searchField := []string{"vcpu", "memory"}

			for _, spec := range vdcRules {
				if slices.Contains(searchField, spec.Target) {
					pR = append(pR, spec)
				}
			}
			return pR
		}(),

		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsUpdateVDC)
			ep := endpoints.UpdateVdc()

			logger := cc.logger.WithGroup("UpdateVDC")

			apiR := apiRequestUpdateVDC{
				VDC: apiRequestUpdateVDCVDC{
					Name: p.Name,
				},
			}

			if p.Vcpu != nil || p.Name == "" {
				vdc, err := cc.GetVDC(ctx, ParamsGetVDC{
					ID:   p.ID,
					Name: p.Name,
				})
				if err != nil {
					logger.Error("Failed to get VDC", "error", err)
					return nil, err
				}

				apiR.VDC.Name = vdc.Name

				if p.Vcpu != nil {
					apiR.VDC.CPUAllocated = serviceClassToCPUInMhz(vdc.Properties.ServiceClass) * *p.Vcpu
				}
			}

			if p.Description != nil {
				apiR.VDC.Description = *p.Description
			}

			if p.Memory != nil {
				apiR.VDC.MemoryAllocated = *p.Memory
			}

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.Name),
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

	// * DeleteVDC
	cmds.Register(commands.Command{
		Namespace:          "VDC",
		Verb:               "Delete",
		ShortDocumentation: "DeleteVDC deletes an existing VDC",
		LongDocumentation:  "Delete VDC performs a DELETE request to delete an existing VDC.",
		AutoGenerate:       true,
		ParamsType:         ParamsDeleteVDC{},
		ParamsSpecs: commands.ParamsSpecs{
			commands.ParamsSpec{
				Name:        "id",
				Description: "ID of the VDC to delete",
				Required:    false,
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("name"),
					commands.ValidatorOmitempty(),
					commands.ValidatorURN("vdc"),
				},
			},
			commands.ParamsSpec{
				Name:        "name",
				Description: "Name of the VDC to delete",
				Required:    false,
				Example:     "my-vdc",
				Validators: []commands.Validator{
					commands.ValidatorRequiredIfParamIsNull("id"),
					commands.ValidatorOmitempty(),
				},
			},
		},
		RunnerFunc: func(ctx context.Context, cmd *commands.Command, client, params any) (any, error) {
			cc := client.(*Client)
			p := params.(ParamsDeleteVDC)
			ep := endpoints.DeleteVdc()

			logger := cc.logger.WithGroup("DeleteVDC")

			if p.Name == "" {
				// Delete require vdc name
				vdc, err := cc.GetVDC(ctx, ParamsGetVDC{
					ID: p.ID,
				})
				if err != nil {
					logger.Error("Failed to get VDC", "error", err)
					return nil, err
				}

				p.Name = vdc.Name
			}

			_, err := cc.c.Do(
				ctx,
				ep,
				cav.WithPathParam(ep.PathParams[0], p.Name),
			)
			if err != nil {
				logger.Error("Failed to delete VDC", "error", err)
				return nil, err
			}

			return nil, nil
		},
	})
}
