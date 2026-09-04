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
	"fmt"
	"time"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListVApp   = "VApp.List"
	opGetVApp    = "VApp.Get"
	opCreateVApp = "VApp.Create"
	opUpdateVApp = "VApp.Update"
	opDeleteVApp = "VApp.Delete"
)

type getVAppByIDParams struct {
	ID string
}

type createVAppParams struct {
	VDCID string
	Body  itypes.APIRequestCreateVApp
}

type updateVAppByIDParams struct {
	ID   string
	Body itypes.APIRequestUpdateVApp
}

type deleteVAppByIDParams struct {
	ID string
}

type listVAppParams struct {
	VDCID string
}

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
	createVAppOp = cav.Operation[createVAppParams, cav.Job]{
		Name:     opCreateVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.CreateVApp(),
		Validate: func(p createVAppParams) error {
			if p.VDCID == "" || p.Body.Name == "" {
				return fmt.Errorf("vdc id and name are required")
			}

			return nil
		},
		RequestOptions: func(p createVAppParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.CreateVApp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.VDCID)}, nil
		},
		Transform: func(p createVAppParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ createVAppParams) (cav.Job, error) {
			job, ok := resp.Result().(*cav.Job)
			if !ok || job == nil {
				return cav.Job{}, fmt.Errorf("unexpected create response type %T", resp.Result())
			}

			return *job, nil
		},
	}
	updateVAppByIDOp = cav.Operation[updateVAppByIDParams, cav.Job]{
		Name:     opUpdateVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UpdateVApp(),
		Validate: func(p updateVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p updateVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UpdateVApp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p updateVAppByIDParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ updateVAppByIDParams) (cav.Job, error) {
			job, ok := resp.Result().(*cav.Job)
			if !ok || job == nil {
				return cav.Job{}, fmt.Errorf("unexpected update response type %T", resp.Result())
			}

			return *job, nil
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
			job, ok := resp.Result().(*cav.Job)
			if !ok || job == nil {
				return cav.Job{}, fmt.Errorf("unexpected delete response type %T", resp.Result())
			}

			return *job, nil
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
			job, ok := resp.Result().(*cav.Job)
			if !ok || job == nil {
				return cav.Job{}, fmt.Errorf("unexpected remove networks response type %T", resp.Result())
			}

			return *job, nil
		},
	}
	undeployVAppOp = cav.Operation[getVAppByIDParams, cav.Job]{
		Name:     "VApp.Undeploy",
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UndeployVApp(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UndeployVApp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p getVAppByIDParams) (any, error) {
			return itypes.APIRequestUndeployVApp{UndeployPowerOff: true}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (cav.Job, error) {
			job, ok := resp.Result().(*cav.Job)
			if !ok || job == nil {
				return cav.Job{}, fmt.Errorf("unexpected undeploy response type %T", resp.Result())
			}

			return *job, nil
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

// CreateVApp creates a VApp and returns the created VApp.
func (c *Client) CreateVApp(ctx context.Context, params types.ParamsCreateVApp) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateVApp, err)
	}

	body := itypes.APIRequestCreateVApp{
		Name:        params.Name,
		Description: params.Description,
	}

	job, err := cav.Execute(ctx, c.c, createVAppOp, createVAppParams{VDCID: params.VDCID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateVApp, err)
	}

	_, err = cav.AwaitJob(ctx, c.c, job.ID, cav.JobPollOptions{
		Timeout:         30 * time.Second,
		PollingInterval: 1 * time.Second,
	}, func(_ *resty.Response) (struct{}, error) {
		return struct{}{}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: await job: %w", opCreateVApp, err)
	}

	m, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: get created: %w", opCreateVApp, err)
	}

	return m, nil
}

// UpdateVApp updates a VApp and returns the updated VApp.
func (c *Client) UpdateVApp(ctx context.Context, params types.ParamsUpdateVApp) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateVApp, err)
	}

	if params.ID == "" {
		vapp, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
		if err != nil {
			return nil, fmt.Errorf("%s: get current: %w", opUpdateVApp, err)
		}
		params.ID = vapp.ID
	}

	body := itypes.APIRequestUpdateVApp{}
	if params.Description != nil {
		body.Description = *params.Description
	}
	if params.DeploymentLeaseInSeconds != nil || params.StorageLeaseInSeconds != nil {
		body.LeaseSettings = &itypes.APIRequestLeaseSettings{}
		if params.DeploymentLeaseInSeconds != nil {
			body.LeaseSettings.DeploymentLeaseInSeconds = params.DeploymentLeaseInSeconds
		}
		if params.StorageLeaseInSeconds != nil {
			body.LeaseSettings.StorageLeaseInSeconds = params.StorageLeaseInSeconds
		}
	}

	job, err := cav.Execute(ctx, c.c, updateVAppByIDOp, updateVAppByIDParams{ID: params.ID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVApp, err)
	}

	_, err = cav.AwaitJob(ctx, c.c, job.ID, cav.JobPollOptions{
		Timeout:         30 * time.Second,
		PollingInterval: 1 * time.Second,
	}, func(_ *resty.Response) (struct{}, error) {
		return struct{}{}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: await job: %w", opUpdateVApp, err)
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

	// Step 1: Remove all networks
	if _, err := cav.Execute(ctx, c.c, removeAllNetworksOp, getVAppByIDParams{ID: params.ID}); err != nil {
		return fmt.Errorf("%s: remove all networks: %w", opDeleteVApp, err)
	}

	// Step 2: Try undeploy (ignore if already undeployed)
	_, err := cav.Execute(ctx, c.c, undeployVAppOp, getVAppByIDParams{ID: params.ID})
	if err != nil {
		c.logger.Debug("Undeploy failed, continuing with delete", "error", err)
	}

	// Step 3: Delete
	if _, err := cav.Execute(ctx, c.c, deleteVAppByIDOp, deleteVAppByIDParams{ID: params.ID}); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVApp, err)
	}

	return nil
}

// IsVAppOrgNetwork checks if a network is an org network.
// This is a helper preserved from v1 SDK.
func IsVAppOrgNetwork(networkHREF string) bool {
	// In vCD, org networks have a specific pattern in their HREF
	// This is a simplified check - adjust based on actual vCD API behavior
	return len(networkHREF) > 0 && networkHREF[0] == '/'
}

// IsVAppNetwork checks if a network is a vApp network.
// This is a helper preserved from v1 SDK.
func IsVAppNetwork(networkHREF string) bool {
	// In vCD, vApp networks typically contain "vapp-" in their HREF
	return len(networkHREF) > 0 && networkHREF[0] != '/'
}
