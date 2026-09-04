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

	"golang.org/x/sync/errgroup"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListVDC   = "VDC.List"
	opGetVDC    = "VDC.Get"
	opCreateVDC = "VDC.Create"
	opUpdateVDC = "VDC.Update"
	opDeleteVDC = "VDC.Delete"
)

// ListVDC lists VDCs visible to current organization.
func (c *Client) ListVDC(ctx context.Context, params types.ParamsListVDC) (*types.ModelListVDC, error) {
	ep := endpoints.ListVDC()

	query := ""
	if params.Name != "" {
		query = fmt.Sprintf("name==%s", params.Name)
	}
	if params.ID != "" {
		query = fmt.Sprintf("id==%s", params.ID)
	}

	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], query),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListVDC, err)
	}

	return resp.Result().(*itypes.APIResponseListVDC).ToModel(), nil
}

// GetVDC returns detailed information for a VDC by ID or name.
func (c *Client) GetVDC(ctx context.Context, params types.ParamsGetVDC) (*types.ModelGetVDC, error) {
	results, err := c.ListVDC(ctx, types.ParamsListVDC{ID: params.ID, Name: params.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opGetVDC, err)
	}

	if len(results.VDCS) == 0 {
		return nil, fmt.Errorf("%s: no VDC found with the provided parameters", opGetVDC)
	}
	vdc := results.VDCS[0]

	var (
		vdcMetadata *itypes.APIResponseGetVDCMetadatas
		model       types.ModelGetVDC
	)

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		epGetVDCMetadata := endpoints.GetVDCMetadata()
		vdcMetadataResp, err := c.c.Do(
			egCtx,
			epGetVDCMetadata,
			cav.WithPathParam(epGetVDCMetadata.PathParams[0], vdc.ID),
		)
		if err != nil {
			return fmt.Errorf("get metadata: %w", err)
		}

		vdcMetadata = vdcMetadataResp.Result().(*itypes.APIResponseGetVDCMetadatas)
		return nil
	})

	eg.Go(func() error {
		epGetVDC := endpoints.GetVDC()
		vdcResp, err := c.c.Do(
			egCtx,
			epGetVDC,
			cav.WithPathParam(epGetVDC.PathParams[0], vdc.ID),
		)
		if err != nil {
			return fmt.Errorf("get details: %w", err)
		}

		vdcDetails := vdcResp.Result().(*itypes.APIResponseGetVDC)
		model = vdcDetails.ToModel()
		model.NumberOfDisks = vdc.NumberOfDisks
		model.NumberOfStorageProfiles = vdc.NumberOfStorageProfiles
		model.NumberOfVMS = vdc.NumberOfVMS
		model.NumberOfRunningVMS = vdc.NumberOfRunningVMS
		model.NumberOfVAPPS = vdc.NumberOfVAPPS

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVDC, err)
	}

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

	return &model, nil
}

// CreateVDC creates a VDC and returns refreshed state.
func (c *Client) CreateVDC(ctx context.Context, params types.ParamsCreateVDC) (*types.ModelGetVDC, error) {
	if err := validateCreateVDCParams(params); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateVDC, err)
	}

	reqBody := itypes.APIRequestCreateVDC{
		VDC: itypes.APIRequestCreateVDCVDC{
			Name:                params.Name,
			Description:         params.Description,
			ServiceClass:        params.ServiceClass,
			DisponibilityClass:  params.DisponibilityClass,
			BillingModel:        params.BillingModel,
			StorageBillingModel: params.StorageBillingModel,
			VCPUInMhz:           serviceClassToCPUInMhz(params.ServiceClass),
			CPUAllocated:        serviceClassToCPUInMhz(params.ServiceClass) * params.Vcpu,
			MemoryAllocated:     params.Memory,
			StorageProfiles:     make([]itypes.APIRequestVDCStorageProfile, len(params.StorageProfiles)),
		},
	}

	for i, sp := range params.StorageProfiles {
		reqBody.VDC.StorageProfiles[i] = itypes.APIRequestVDCStorageProfile{
			Class:   sp.Class,
			Limit:   sp.Limit,
			Default: sp.Default,
		}
	}

	ep := endpoints.CreateVDC()
	if _, err := c.c.Do(ctx, ep, cav.SetBody(reqBody)); err != nil {
		return nil, fmt.Errorf("%s: create: %w", opCreateVDC, err)
	}

	resp, err := c.GetVDC(ctx, types.ParamsGetVDC{Name: params.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opCreateVDC, err)
	}

	return resp, nil
}

// UpdateVDC updates selected VDC fields and returns refreshed state.
func (c *Client) UpdateVDC(ctx context.Context, params types.ParamsUpdateVDC) (*types.ModelGetVDC, error) {
	if err := validateUpdateVDCParams(params); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateVDC, err)
	}

	ep := endpoints.UpdateVDC()
	apiR := itypes.APIRequestUpdateVDC{
		VDC: itypes.APIRequestUpdateVDCVDC{Name: params.Name},
	}

	if params.Vcpu != nil || params.Name == "" {
		vdc, err := c.GetVDC(ctx, types.ParamsGetVDC{ID: params.ID, Name: params.Name})
		if err != nil {
			return nil, fmt.Errorf("%s: get current: %w", opUpdateVDC, err)
		}

		apiR.VDC.Name = vdc.Name
		if params.Vcpu != nil {
			apiR.VDC.CPUAllocated = serviceClassToCPUInMhz(vdc.Properties.ServiceClass) * *params.Vcpu
		}
	}

	if params.Description != nil {
		apiR.VDC.Description = *params.Description
	}
	if params.Memory != nil {
		apiR.VDC.MemoryAllocated = *params.Memory
	}

	if _, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], apiR.VDC.Name),
		cav.SetBody(apiR),
	); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateVDC, err)
	}

	return c.GetVDC(ctx, types.ParamsGetVDC{ID: params.ID, Name: params.Name})
}

// DeleteVDC deletes a VDC by ID or name.
func (c *Client) DeleteVDC(ctx context.Context, params types.ParamsDeleteVDC) error {
	if params.ID == "" && params.Name == "" {
		return fmt.Errorf("%s: id or name is required", opDeleteVDC)
	}

	name := params.Name
	if name == "" {
		vdc, err := c.GetVDC(ctx, types.ParamsGetVDC{ID: params.ID})
		if err != nil {
			return fmt.Errorf("%s: get current: %w", opDeleteVDC, err)
		}
		name = vdc.Name
	}

	ep := endpoints.DeleteVDC()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], name)); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVDC, err)
	}

	return nil
}

func validateCreateVDCParams(params types.ParamsCreateVDC) error {
	if params.Name == "" || params.ServiceClass == "" || params.DisponibilityClass == "" || params.BillingModel == "" || params.StorageBillingModel == "" {
		return fmt.Errorf("missing required parameters")
	}
	if len(params.StorageProfiles) == 0 {
		return fmt.Errorf("missing required parameters")
	}

	haveOneDefaultStorageProfile := false
	for _, sp := range params.StorageProfiles {
		if sp.Default {
			haveOneDefaultStorageProfile = true
			break
		}
	}
	if !haveOneDefaultStorageProfile {
		return fmt.Errorf("at least one storage profile must be marked as default")
	}

	return nil
}

func validateUpdateVDCParams(params types.ParamsUpdateVDC) error {
	if params.ID == "" && params.Name == "" {
		return fmt.Errorf("missing required parameters")
	}
	if params.Description == nil && params.Vcpu == nil && params.Memory == nil {
		return fmt.Errorf("missing required parameters")
	}
	return nil
}

func serviceClassToCPUInMhz(serviceClass string) int {
	switch serviceClass {
	case "VOIP":
		return 3000
	default:
		return 2200
	}
}
