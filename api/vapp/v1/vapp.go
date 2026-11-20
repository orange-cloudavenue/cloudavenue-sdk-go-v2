/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListVApp                = "VApp.List"
	opGetVApp                 = "VApp.Get"
	opComposeVApp             = "VApp.Compose"
	opInstantiateVAppTemplate = "VApp.InstantiateTemplate"
	opUpdateVApp              = "VApp.Update"
	opDeleteVApp              = "VApp.Delete"
)

var (
	listVAppOp = cav.Operation[listVAppParams, *itypes.APIResponseListVApp]{
		Name:     opListVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.ListVApp(),
		Validate: func(p listVAppParams) error {
			if p.VDCID == "" {
				return fmt.Errorf("vdc id is required")
			}
			return nil
		},
		RequestOptions: func(p listVAppParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.ListVApp()
			return []cav.EndpointRequestOption{
				cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("id==%s", p.VDCID)),
				cav.WithQueryParam(ep.QueryParams[1], "100"),
				cav.WithQueryParam(ep.QueryParams[2], "records"),
				cav.WithQueryParam(ep.QueryParams[3], "vApp"),
			}, nil
		},
		Extract: func(resp *cav.Response, _ listVAppParams) (*itypes.APIResponseListVApp, error) {
			list, ok := resp.Result().(*itypes.APIResponseListVApp)
			if !ok || list == nil {
				return nil, fmt.Errorf("unexpected list response type %T", resp.Result())
			}

			return list, nil
		},
	}
	getVAppByIDOp = cav.Operation[getVAppByIDParams, *itypes.APIResponseGetVApp]{
		Name:     opGetVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.GetVApp(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.GetVApp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (*itypes.APIResponseGetVApp, error) {
			vapp, ok := resp.Result().(*itypes.APIResponseGetVApp)
			if !ok || vapp == nil {
				return nil, fmt.Errorf("unexpected get response type %T", resp.Result())
			}

			return vapp, nil
		},
	}
	composeVAppOp = cav.Operation[composeVAppParams, cav.Job]{
		Name:     opComposeVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.ComposeVApp(),
		Validate: func(p composeVAppParams) error {
			if p.VDCID == "" || p.Body.Name == "" {
				return fmt.Errorf("vdc id and name are required")
			}

			return nil
		},
		RequestOptions: func(p composeVAppParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.ComposeVApp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.VDCID)}, nil
		},
		Transform: func(p composeVAppParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ composeVAppParams) (cav.Job, error) {
			return extractVMwareJob(resp, "compose")
		},
	}
	instantiateVAppTemplateOp = cav.Operation[instantiateVAppTemplateParams, cav.Job]{
		Name:     opInstantiateVAppTemplate,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.InstantiateVAppTemplate(),
		Validate: func(p instantiateVAppTemplateParams) error {
			if p.VDCID == "" || p.Body.Name == "" || p.Body.Source.ID == "" && p.Body.Source.HREF == "" {
				return fmt.Errorf("vdc id, name, and template source are required")
			}
			return nil
		},
		RequestOptions: func(p instantiateVAppTemplateParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.InstantiateVAppTemplate()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.VDCID)}, nil
		},
		Transform: func(p instantiateVAppTemplateParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ instantiateVAppTemplateParams) (cav.Job, error) {
			return extractVMwareJob(resp, "instantiate template")
		},
	}
	recomposeVAppOp = cav.Operation[recomposeVAppParams, cav.Job]{
		Name:     opUpdateVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.RecomposeVApp(),
		Validate: func(p recomposeVAppParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p recomposeVAppParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.RecomposeVApp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p recomposeVAppParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ recomposeVAppParams) (cav.Job, error) {
			return extractVMwareJob(resp, "recompose")
		},
	}
	updateVAppLeaseSettingsSectionOp = cav.Operation[updateVAppLeaseSettingsSectionParams, cav.Job]{
		Name:     opUpdateVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UpdateVAppLeaseSettingsSection(),
		Validate: func(p updateVAppLeaseSettingsSectionParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p updateVAppLeaseSettingsSectionParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UpdateVAppLeaseSettingsSection()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p updateVAppLeaseSettingsSectionParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ updateVAppLeaseSettingsSectionParams) (cav.Job, error) {
			return extractVMwareJob(resp, "update lease settings")
		},
	}
	deleteVAppByIDOp = cav.Operation[deleteVAppByIDParams, cav.Job]{
		Name:     opDeleteVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.DeleteVApp(),
		Validate: func(p deleteVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p deleteVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.DeleteVApp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ deleteVAppByIDParams) (cav.Job, error) {
			return extractVMwareJob(resp, "delete")
		},
	}
	removeAllNetworksOp = cav.Operation[getVAppByIDParams, cav.Job]{
		Name:     "VApp.RemoveAllNetworks",
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.RemoveAllNetworks(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.RemoveAllNetworks()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (cav.Job, error) {
			return extractVMwareJob(resp, "remove networks")
		},
	}
)

// ListVApp lists VApps in a VDC.
func (c *Client) ListVApp(ctx context.Context, vdcID string) ([]*types.ModelVApp, error) {
	if vdcID == "" {
		return nil, fmt.Errorf("%s: vdc id is required", opListVApp)
	}

	resp, err := cav.Execute(ctx, c.c, listVAppOp, listVAppParams{VDCID: vdcID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListVApp, err)
	}

	vapps := make([]*types.ModelVApp, 0, len(resp.Records))
	for i := range resp.Records {
		m := resp.Records[i].ToModel()
		vapps = append(vapps, &m)
	}

	return vapps, nil
}

// GetVApp returns detailed information for a VApp by ID or name.
func (c *Client) GetVApp(ctx context.Context, params types.ParamsGetVApp) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetVApp, err)
	}

	if params.ID == "" {
		ep := endpoints.ListVApp()
		opts := []cav.EndpointRequestOption{
			cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("name==%s", params.Name)),
			cav.WithQueryParam(ep.QueryParams[1], "100"),
			cav.WithQueryParam(ep.QueryParams[2], "records"),
			cav.WithQueryParam(ep.QueryParams[3], "vApp"),
		}

		rawResp, err := c.c.Do(ctx, ep, opts...)
		if err != nil {
			return nil, fmt.Errorf("%s: list: %w", opGetVApp, err)
		}

		list, ok := rawResp.Result().(*itypes.APIResponseListVApp)
		if !ok || list == nil {
			return nil, fmt.Errorf("%s: unexpected list response type %T", opGetVApp, rawResp.Result())
		}

		for _, vapp := range list.Records {
			m := vapp.ToModel()
			if m.Name == params.Name {
				// The list record has no VDC nor lease: read the full vApp by ID.
				if m.ID != "" {
					return c.GetVApp(ctx, types.ParamsGetVApp{ID: m.ID})
				}
				return &m, nil
			}
		}

		return nil, fmt.Errorf("%s: no VApp found with name %q", opGetVApp, params.Name)
	}

	resp, err := cav.Execute(ctx, c.c, getVAppByIDOp, getVAppByIDParams{ID: params.ID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVApp, err)
	}

	m := resp.ToModel()
	return &m, nil
}

// ComposeVApp composes a VApp from one or more source resources and returns it.
func (c *Client) ComposeVApp(ctx context.Context, params types.ParamsComposeVApp) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opComposeVApp, err)
	}

	body := itypes.APIRequestComposeVApp{
		OVF:              ovfEnvelopeXMLNS,
		XSI:              xmlSchemaInstanceNS,
		XMLNS:            vcloudXMLNS,
		Name:             params.Name,
		Description:      params.Description,
		Deploy:           params.Deploy,
		PowerOn:          params.PowerOn,
		LinkedClone:      params.LinkedClone,
		AllEULAsAccepted: params.AllEULAsAccepted,
	}
	if params.SourcedItem != nil {
		body.SourcedItem = &itypes.APIRequestVAppSourcedItem{
			Source:            itypes.APIRequestVAppReference(params.SourcedItem.Source),
			SourceDelete:      params.SourcedItem.SourceDelete,
			VAppScopedLocalID: params.SourcedItem.VAppScopedLocalID,
		}
		if params.SourcedItem.StorageProfile != nil {
			ref := itypes.APIRequestVAppReference(*params.SourcedItem.StorageProfile)
			body.SourcedItem.StorageProfile = &ref
		}
	}
	body.InstantiationParams = toInstantiationParams(params.InstantiationParams)

	job, err := cav.Execute(ctx, c.c, composeVAppOp, composeVAppParams{VDCID: params.VDCID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opComposeVApp, err)
	}

	if job.Status != cav.JobSuccess {
		if err := c.awaitJob(ctx, opComposeVApp, job.ID); err != nil {
			return nil, err
		}
	}

	m, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: get composed: %w", opComposeVApp, err)
	}

	return m, nil
}

// InstantiateVAppTemplate creates a VApp from a VApp template and returns it.
func (c *Client) InstantiateVAppTemplate(ctx context.Context, params types.ParamsInstantiateVAppTemplate) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opInstantiateVAppTemplate, err)
	}

	body := itypes.APIRequestInstantiateVAppTemplate{
		OVF:                 ovfEnvelopeXMLNS,
		XSI:                 xmlSchemaInstanceNS,
		XMLNS:               vcloudXMLNS,
		Name:                params.Name,
		Description:         params.Description,
		Deploy:              params.Deploy,
		PowerOn:             params.PowerOn,
		LinkedClone:         params.LinkedClone,
		AllEULAsAccepted:    params.AllEULAsAccepted,
		InstantiationParams: toInstantiationParams(params.InstantiationParams),
		Source:              itypes.APIRequestVAppReference{ID: params.TemplateID},
	}

	job, err := cav.Execute(ctx, c.c, instantiateVAppTemplateOp, instantiateVAppTemplateParams{VDCID: params.VDCID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opInstantiateVAppTemplate, err)
	}
	if job.Status != cav.JobSuccess {
		if err := c.awaitJob(ctx, opInstantiateVAppTemplate, job.ID); err != nil {
			return nil, err
		}
	}
	return c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
}

// UpdateVApp updates a VApp and returns the updated VApp.
func (c *Client) UpdateVApp(ctx context.Context, params types.ParamsUpdateVApp) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateVApp, err)
	}

	var current *types.ModelVApp
	if params.ID == "" {
		vapp, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
		if err != nil {
			return nil, fmt.Errorf("%s: get current: %w", opUpdateVApp, err)
		}
		params.ID = vapp.ID
		current = vapp
	}

	if current == nil {
		vapp, err := c.GetVApp(ctx, types.ParamsGetVApp{ID: params.ID})
		if err != nil {
			return nil, fmt.Errorf("%s: get current: %w", opUpdateVApp, err)
		}
		current = vapp
	}

	// The description is changed with the recompose action.
	if params.Description != nil && *params.Description != current.Description {
		job, err := cav.Execute(ctx, c.c, recomposeVAppOp, recomposeVAppParams{ID: params.ID, Body: itypes.APIRequestRecomposeVApp{
			OVF:         ovfEnvelopeXMLNS,
			XSI:         xmlSchemaInstanceNS,
			XMLNS:       vcloudXMLNS,
			Name:        current.Name,
			Deploy:      current.Deployed,
			Description: *params.Description,
		}})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", opUpdateVApp, err)
		}
		if job.Status != cav.JobSuccess {
			if err := c.awaitJob(ctx, opUpdateVApp, job.ID); err != nil {
				return nil, err
			}
		}
	}

	// The lease is changed on its own section: both values are sent, the one
	// that is not requested keeps its current value.
	if leaseChanged(params, current) {
		lease := itypes.APIRequestLeaseSettingsSection{
			OVF:                      ovfEnvelopeXMLNS,
			XMLNS:                    vcloudXMLNS,
			Info:                     "Lease section settings",
			DeploymentLeaseInSeconds: mergedLease(params.DeploymentLeaseInSeconds, current.DeploymentLeaseInSeconds),
			StorageLeaseInSeconds:    mergedLease(params.StorageLeaseInSeconds, current.StorageLeaseInSeconds),
		}

		job, err := cav.Execute(ctx, c.c, updateVAppLeaseSettingsSectionOp, updateVAppLeaseSettingsSectionParams{ID: params.ID, Body: lease})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", opUpdateVApp, err)
		}
		if job.Status != cav.JobSuccess {
			if err := c.awaitJob(ctx, opUpdateVApp, job.ID); err != nil {
				return nil, err
			}
		}
	}

	m, err := c.GetVApp(ctx, types.ParamsGetVApp{ID: params.ID})
	if err != nil {
		return nil, fmt.Errorf("%s: get updated: %w", opUpdateVApp, err)
	}

	return m, nil
}

// DeleteVApp deletes a VApp with the proper sequence:
// 1. Remove all networks
// 2. Try undeploy (ignore if already undeployed)
// 3. Delete
func (c *Client) DeleteVApp(ctx context.Context, params types.ParamsDeleteVApp) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opDeleteVApp, err)
	}

	if params.ID == "" {
		vapp, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
		if err != nil {
			return fmt.Errorf("%s: get current: %w", opDeleteVApp, err)
		}
		params.ID = vapp.ID
	}

	job, err := cav.Execute(ctx, c.c, removeAllNetworksOp, getVAppByIDParams{ID: params.ID})
	if err != nil {
		if !isIgnorableRemoveAllNetworksError(err) {
			return fmt.Errorf("%s: remove all networks: %w", opDeleteVApp, err)
		}
		// A vApp without networks answers 404 to removeAllNetworks; there is
		// nothing to remove, so skip the step instead of aborting the delete.
		c.logger.Debug("Remove all networks skipped during delete because VApp has no networks", "error", err)
	} else if err := c.awaitJob(ctx, opDeleteVApp, job.ID); err != nil {
		return fmt.Errorf("%s: remove all networks: %w", opDeleteVApp, err)
	}

	job, err = cav.Execute(ctx, c.c, undeployVAppOp, undeployVAppParams{ID: params.ID, Body: itypes.ToAPIRequestUndeployVApp(types.VAppUndeployModePowerOff)})
	if err != nil {
		if !isIgnorableUndeployError(err) {
			return fmt.Errorf("%s: undeploy: %w", opDeleteVApp, err)
		}
		c.logger.Debug("Undeploy skipped during delete because VApp is already undeployed", "error", err)
	} else if err := c.awaitJob(ctx, opDeleteVApp, job.ID); err != nil {
		if !isIgnorableUndeployError(err) {
			return fmt.Errorf("%s: undeploy: %w", opDeleteVApp, err)
		}
		c.logger.Debug("Undeploy completion skipped during delete because VApp is already undeployed", "error", err)
	}

	job, err = cav.Execute(ctx, c.c, deleteVAppByIDOp, deleteVAppByIDParams{ID: params.ID})
	if err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVApp, err)
	}
	if err := c.awaitJob(ctx, opDeleteVApp, job.ID); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVApp, err)
	}

	return nil
}

// isIgnorableRemoveAllNetworksError reports whether err is the 404 that VCD
// returns when removeAllNetworks is invoked on a vApp that has no networks.
func isIgnorableRemoveAllNetworksError(err error) bool {
	var apiErr *pkgerrors.APIError
	if !errors.As(err, &apiErr) || apiErr == nil {
		return false
	}

	return apiErr.StatusCode == http.StatusNotFound
}

func (c *Client) resolveVAppID(ctx context.Context, id, name string) (string, error) {
	if id != "" {
		return id, nil
	}
	vapp, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: name})
	if err != nil {
		return "", fmt.Errorf("get current: %w", err)
	}
	return vapp.ID, nil
}

func toInstantiationParams(params *types.VAppInstantiationParams) *itypes.APIRequestVAppInstantiationParams {
	if params == nil {
		return nil
	}

	result := &itypes.APIRequestVAppInstantiationParams{}
	if params.NetworkConfigSection != nil {
		section := &itypes.APIRequestVAppNetworkConfigSection{Info: params.NetworkConfigSection.Info}
		for _, network := range params.NetworkConfigSection.NetworkConfig {
			wire := itypes.APIRequestVAppNetworkConfig{
				NetworkName:   network.NetworkName,
				Configuration: itypes.APIRequestVAppNetworkConfiguration{FenceMode: network.FenceMode},
			}
			if network.ParentNetwork != nil {
				wire.Configuration.ParentNetwork = &itypes.APIRequestVAppReference{
					HREF: network.ParentNetwork.HREF,
					ID:   network.ParentNetwork.ID,
					Type: network.ParentNetwork.Type,
					Name: network.ParentNetwork.Name,
				}
			}
			section.NetworkConfig = append(section.NetworkConfig, wire)
		}
		result.NetworkConfigSection = section
	}
	if params.NetworkConnectionSection != nil {
		result.NetworkConnectionSection = &itypes.APIRequestVAppNetworkConnectionSection{
			Info:                          params.NetworkConnectionSection.Info,
			PrimaryNetworkConnectionIndex: params.NetworkConnectionSection.PrimaryNetworkConnectionIndex,
		}
	}
	return result
}

// leaseChanged reports whether the update changes the lease of the vApp.
func leaseChanged(params types.ParamsUpdateVApp, current *types.ModelVApp) bool {
	return differs(params.DeploymentLeaseInSeconds, current.DeploymentLeaseInSeconds) ||
		differs(params.StorageLeaseInSeconds, current.StorageLeaseInSeconds)
}

// differs reports whether a requested lease (nil: not requested) differs from the current one.
// An unknown current lease (nil) is considered different.
func differs(requested, current *int) bool {
	if requested == nil {
		return false
	}
	return current == nil || *requested != *current
}

// mergedLease returns the requested lease, or the current one when not requested (0 when unknown).
func mergedLease(requested, current *int) *int {
	if requested != nil {
		v := *requested
		return &v
	}
	v := 0
	if current != nil {
		v = *current
	}
	return &v
}
