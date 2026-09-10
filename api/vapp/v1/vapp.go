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
	"path"
	"net/url"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func awaitVAppJob(ctx context.Context, c cav.Client, opName, jobID string) error {
	_, err := cav.AwaitJobOnBackend(ctx, c, cav.BackendVMware, jobID, cav.JobPollOptions{
		Timeout:         30 * time.Second,
		PollingInterval: 1 * time.Second,
	}, func(_ *resty.Response) (struct{}, error) {
		return struct{}{}, nil
	})
	if err != nil {
		return fmt.Errorf("%s: await job: %w", opName, err)
	}
	return nil
}

func networkHrefPath(networkHREF string) string {
	parsed, err := url.Parse(networkHREF)
	if err != nil || parsed.Path == "" {
		return networkHREF
	}
	return parsed.Path
}

func extractVMwareJob(resp *cav.Response, action string) (cav.Job, error) {
	if job, ok := resp.Result().(*cav.Job); ok && job != nil && job.ID != "" {
		return *job, nil
	}

	if resp == nil || resp.Raw == nil {
		return cav.Job{}, fmt.Errorf("unexpected %s response: missing raw response", action)
	}

	location := resp.Raw.Header().Get("Location")
	if location == "" {
		return cav.Job{}, fmt.Errorf("unexpected %s response: missing job location", action)
	}

	parsed, err := url.Parse(location)
	if err != nil {
		return cav.Job{}, fmt.Errorf("unexpected %s response: parse location: %w", action, err)
	}

	jobID := path.Base(parsed.Path)
	if jobID == "" || jobID == "." || jobID == "/" {
		return cav.Job{}, fmt.Errorf("unexpected %s response: missing job id in location %q", action, location)
	}

	return cav.Job{ID: jobID}, nil
}

const (
	opListVapp   = "Vapp.List"
	opGetVapp    = "Vapp.Get"
	opCreateVapp = "Vapp.Create"
	opUpdateVapp = "Vapp.Update"
	opDeleteVapp = "Vapp.Delete"
)

type getVAppByIDParams struct {
	ID string
}

type createVAppParams struct {
	VDCID string
	Body  itypes.APIRequestCreateVApp
}

type updateVAppByIDParams struct {
	ID    string
	Body  itypes.APIRequestUpdateVApp
}

type deleteVAppByIDParams struct {
	ID string
}

type listVAppParams struct {
	VDCID string
}

var (
	listVAppOp = cav.Operation[listVAppParams, *itypes.APIResponseListVApp]{
		Name:     opListVapp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.ListVapp(),
		Validate: func(p listVAppParams) error {
			if p.VDCID == "" {
				return fmt.Errorf("vdc id is required")
			}
			return nil
		},
		RequestOptions: func(p listVAppParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.ListVapp()
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
		Name:     opGetVapp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.GetVapp(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.GetVapp()
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
		Name:     opCreateVapp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.CreateVapp(),
		Validate: func(p createVAppParams) error {
			if p.VDCID == "" || p.Body.Name == "" {
				return fmt.Errorf("vdc id and name are required")
			}

			return nil
		},
		RequestOptions: func(p createVAppParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.CreateVapp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.VDCID)}, nil
		},
		Transform: func(p createVAppParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ createVAppParams) (cav.Job, error) {
			return extractVMwareJob(resp, "create")
		},
	}
	updateVAppByIDOp = cav.Operation[updateVAppByIDParams, cav.Job]{
		Name:     opUpdateVapp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UpdateVapp(),
		Validate: func(p updateVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p updateVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UpdateVapp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p updateVAppByIDParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ updateVAppByIDParams) (cav.Job, error) {
			return extractVMwareJob(resp, "update")
		},
	}
	deleteVAppByIDOp = cav.Operation[deleteVAppByIDParams, cav.Job]{
		Name:     opDeleteVapp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.DeleteVapp(),
		Validate: func(p deleteVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p deleteVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.DeleteVapp()
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
	undeployVAppOp = cav.Operation[getVAppByIDParams, cav.Job]{
		Name:     "VApp.Undeploy",
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UndeployVapp(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UndeployVapp()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p getVAppByIDParams) (any, error) {
			return itypes.APIRequestUndeployVApp{UndeployPowerOff: true}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (cav.Job, error) {
			return extractVMwareJob(resp, "undeploy")
		},
	}
)

// ListVApp lists VApps in a VDC.
func (c *Client) ListVApp(ctx context.Context, vdcID string) ([]*types.ModelVApp, error) {
	if vdcID == "" {
		return nil, fmt.Errorf("%s: vdc id is required", opListVapp)
	}

	resp, err := cav.Execute(ctx, c.c, listVAppOp, listVAppParams{VDCID: vdcID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListVapp, err)
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
		return nil, fmt.Errorf("%s: validate: %w", opGetVapp, err)
	}

	if params.ID == "" {
		ep := endpoints.ListVapp()
		opts := []cav.EndpointRequestOption{
			cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("name==%s", params.Name)),
			cav.WithQueryParam(ep.QueryParams[1], "100"),
			cav.WithQueryParam(ep.QueryParams[2], "records"),
			cav.WithQueryParam(ep.QueryParams[3], "vApp"),
		}

		rawResp, err := c.c.Do(ctx, ep, opts...)
		if err != nil {
			return nil, fmt.Errorf("%s: list: %w", opGetVapp, err)
		}

		list, ok := rawResp.Result().(*itypes.APIResponseListVApp)
		if !ok || list == nil {
			return nil, fmt.Errorf("%s: unexpected list response type %T", opGetVapp, rawResp.Result())
		}

		for _, vapp := range list.Records {
			m := vapp.ToModel()
			if m.Name == params.Name {
				return &m, nil
			}
		}

		return nil, fmt.Errorf("%s: no VApp found with name %q", opGetVapp, params.Name)
	}

	resp, err := cav.Execute(ctx, c.c, getVAppByIDOp, getVAppByIDParams{ID: params.ID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVapp, err)
	}

	m := resp.ToModel()
	return &m, nil
}

// CreateVApp creates a VApp and returns the created VApp.
func (c *Client) CreateVApp(ctx context.Context, params types.ParamsCreateVApp) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateVapp, err)
	}

	body := itypes.APIRequestCreateVApp{
		Name:        params.Name,
		Description: params.Description,
	}

	job, err := cav.Execute(ctx, c.c, createVAppOp, createVAppParams{VDCID: params.VDCID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateVapp, err)
	}

	if err := awaitVAppJob(ctx, c.c, opCreateVapp, job.ID); err != nil {
		return nil, err
	}

	m, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: get created: %w", opCreateVapp, err)
	}

	return m, nil
}

// UpdateVApp updates a VApp and returns the updated VApp.
func (c *Client) UpdateVApp(ctx context.Context, params types.ParamsUpdateVApp) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateVapp, err)
	}

	if params.ID == "" {
		vapp, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
		if err != nil {
			return nil, fmt.Errorf("%s: get current: %w", opUpdateVapp, err)
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
		return nil, fmt.Errorf("%s: %w", opUpdateVapp, err)
	}

	if err := awaitVAppJob(ctx, c.c, opUpdateVapp, job.ID); err != nil {
		return nil, err
	}

	m, err := c.GetVApp(ctx, types.ParamsGetVApp{ID: params.ID})
	if err != nil {
		return nil, fmt.Errorf("%s: get updated: %w", opUpdateVapp, err)
	}

	return m, nil
}

// DeleteVApp deletes a VApp with the proper sequence:
// 1. Remove all networks
// 2. Try undeploy (ignore if already undeployed)
// 3. Delete
func (c *Client) DeleteVApp(ctx context.Context, params types.ParamsDeleteVApp) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opDeleteVapp, err)
	}

	if params.ID == "" {
		vapp, err := c.GetVApp(ctx, types.ParamsGetVApp{Name: params.Name})
		if err != nil {
			return fmt.Errorf("%s: get current: %w", opDeleteVapp, err)
		}
		params.ID = vapp.ID
	}

	// Step 1: Remove all networks
	job, err := cav.Execute(ctx, c.c, removeAllNetworksOp, getVAppByIDParams{ID: params.ID})
	if err != nil {
		return fmt.Errorf("%s: remove all networks: %w", opDeleteVapp, err)
	}
	if err := awaitVAppJob(ctx, c.c, opDeleteVapp, job.ID); err != nil {
		return fmt.Errorf("%s: remove all networks: %w", opDeleteVapp, err)
	}

	// Step 2: Try undeploy (ignore if already undeployed)
	job, err = cav.Execute(ctx, c.c, undeployVAppOp, getVAppByIDParams{ID: params.ID})
	if err != nil {
		c.logger.Debug("Undeploy failed, continuing with delete", "error", err)
	} else if err := awaitVAppJob(ctx, c.c, opDeleteVapp, job.ID); err != nil {
		c.logger.Debug("Await undeploy failed, continuing with delete", "error", err)
	}

	// Step 3: Delete
	job, err = cav.Execute(ctx, c.c, deleteVAppByIDOp, deleteVAppByIDParams{ID: params.ID})
	if err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVapp, err)
	}
	if err := awaitVAppJob(ctx, c.c, opDeleteVapp, job.ID); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVapp, err)
	}

	return nil
}

// IsVAppOrgNetwork checks if a network is an org network.
// This is a helper preserved from v1 SDK.
func IsVAppOrgNetwork(networkHREF string) bool {
	path := networkHrefPath(networkHREF)
	return strings.Contains(path, "/api/network/") || strings.Contains(path, "/network/")
}

// IsVAppNetwork checks if a network is a vApp network.
// This is a helper preserved from v1 SDK.
func IsVAppNetwork(networkHREF string) bool {
	path := strings.ToLower(networkHrefPath(networkHREF))
	return strings.Contains(path, "vappnetwork")
}
