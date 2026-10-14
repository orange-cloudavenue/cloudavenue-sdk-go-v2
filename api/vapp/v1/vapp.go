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
	"net/url"
	"path"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListVApp                = "VApp.List"
	opGetVApp                 = "VApp.Get"
	opComposeVApp             = "VApp.Compose"
	opInstantiateVAppTemplate = "VApp.InstantiateTemplate"
	opUpdateVApp              = "VApp.Update"
	opDeleteVApp              = "VApp.Delete"
	opDeployVApp              = "VApp.Deploy"
	opPowerOffVApp            = "VApp.PowerOff"
	opPowerOnVApp             = "VApp.PowerOn"
	opRebootVApp              = "VApp.Reboot"
	opResetVApp               = "VApp.Reset"
	opShutdownVApp            = "VApp.Shutdown"
	opSuspendVApp             = "VApp.Suspend"
)

type getVAppByIDParams struct {
	ID string
}

type composeVAppParams struct {
	VDCID string
	Body  itypes.APIRequestComposeVApp
}

type instantiateVAppTemplateParams struct {
	VDCID string
	Body  itypes.APIRequestInstantiateVAppTemplate
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
	composeVAppOp = cav.Operation[composeVAppParams, cav.Job]{
		Name:     opComposeVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.ComposeVApp(),
		Validate: func(p composeVAppParams) error {
			if p.VDCID == "" || p.Body.Name == "" || p.Body.SourcedItem == nil {
				return fmt.Errorf("vdc id, name, and sourced item are required")
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
			return extractVMwareJob(resp, "update")
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
			return extractVMwareJob(resp, "undeploy")
		},
	}
	deployVAppOp   = newVAppActionOperation(opDeployVApp, endpoints.DeployVApp, "deploy")
	powerOffVAppOp = newVAppActionOperation(opPowerOffVApp, endpoints.PowerOffVApp, "power off")
	powerOnVAppOp  = newVAppActionOperation(opPowerOnVApp, endpoints.PowerOnVApp, "power on")
	rebootVAppOp   = newVAppActionOperation(opRebootVApp, endpoints.RebootVApp, "reboot")
	resetVAppOp    = newVAppActionOperation(opResetVApp, endpoints.ResetVApp, "reset")
	shutdownVAppOp = newVAppActionOperation(opShutdownVApp, endpoints.ShutdownVApp, "shutdown")
	suspendVAppOp  = newVAppActionOperation(opSuspendVApp, endpoints.SuspendVApp, "suspend")
)

func newVAppActionOperation(name string, endpoint func() *cav.Endpoint, action string) cav.Operation[getVAppByIDParams, cav.Job] {
	return cav.Operation[getVAppByIDParams, cav.Job]{
		Name:     name,
		Backend:  cav.BackendVMware,
		Endpoint: endpoint(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoint()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (cav.Job, error) {
			return extractVMwareJob(resp, action)
		},
	}
}

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

// ComposeVApp composes a VApp from one or more source resources and returns it.
func (c *Client) ComposeVApp(ctx context.Context, params types.ParamsComposeVApp) (*types.ModelVApp, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opComposeVApp, err)
	}

	body := itypes.APIRequestComposeVApp{
		OVF:              "http://schemas.dmtf.org/ovf/envelope/1",
		XSI:              "http://www.w3.org/2001/XMLSchema-instance",
		XMLNS:            "http://www.vmware.com/vcloud/v1.5",
		Name:             params.Name,
		Description:      params.Description,
		Deploy:           params.Deploy,
		PowerOn:          params.PowerOn,
		LinkedClone:      params.LinkedClone,
		AllEULAsAccepted: params.AllEULAsAccepted,
	}
	if params.SourcedItem == nil {
		return nil, fmt.Errorf("%s: sourced item is required", opComposeVApp)
	}
	body.SourcedItem = &itypes.APIRequestVAppSourcedItem{
		Source:            toVAppReference(params.SourcedItem.Source),
		SourceDelete:      params.SourcedItem.SourceDelete,
		VAppScopedLocalID: params.SourcedItem.VAppScopedLocalID,
	}
	if params.SourcedItem.StorageProfile != nil {
		ref := toVAppReference(*params.SourcedItem.StorageProfile)
		body.SourcedItem.StorageProfile = &ref
	}
	body.InstantiationParams = toInstantiationParams(params.InstantiationParams)

	job, err := cav.Execute(ctx, c.c, composeVAppOp, composeVAppParams{VDCID: params.VDCID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opComposeVApp, err)
	}

	if err := awaitVAppJob(ctx, c.c, opComposeVApp, job.ID); err != nil {
		return nil, err
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
		OVF:                 "http://schemas.dmtf.org/ovf/envelope/1",
		XSI:                 "http://www.w3.org/2001/XMLSchema-instance",
		XMLNS:               "http://www.vmware.com/vcloud/v1.5",
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
	if err := awaitVAppJob(ctx, c.c, opInstantiateVAppTemplate, job.ID); err != nil {
		return nil, err
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

	body := itypes.APIRequestUpdateVApp{
		Description: &current.Description,
	}
	if params.Description != nil {
		body.Description = params.Description
	}
	if params.DeploymentLeaseInSeconds != nil || params.StorageLeaseInSeconds != nil || current.DeploymentLeaseInSeconds != nil || current.StorageLeaseInSeconds != nil {
		body.LeaseSettings = &itypes.APIRequestLeaseSettings{}
		if current.DeploymentLeaseInSeconds != nil {
			body.LeaseSettings.DeploymentLeaseInSeconds = current.DeploymentLeaseInSeconds
		}
		if params.DeploymentLeaseInSeconds != nil {
			body.LeaseSettings.DeploymentLeaseInSeconds = params.DeploymentLeaseInSeconds
		}
		if current.StorageLeaseInSeconds != nil {
			body.LeaseSettings.StorageLeaseInSeconds = current.StorageLeaseInSeconds
		}
		if params.StorageLeaseInSeconds != nil {
			body.LeaseSettings.StorageLeaseInSeconds = params.StorageLeaseInSeconds
		}
	}

	job, err := cav.Execute(ctx, c.c, updateVAppByIDOp, updateVAppByIDParams{ID: params.ID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVApp, err)
	}

	if err := awaitVAppJob(ctx, c.c, opUpdateVApp, job.ID); err != nil {
		return nil, err
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
	job, err := cav.Execute(ctx, c.c, removeAllNetworksOp, getVAppByIDParams{ID: params.ID})
	if err != nil {
		return fmt.Errorf("%s: remove all networks: %w", opDeleteVApp, err)
	}
	if err := awaitVAppJob(ctx, c.c, opDeleteVApp, job.ID); err != nil {
		return fmt.Errorf("%s: remove all networks: %w", opDeleteVApp, err)
	}

	// Step 2: Try undeploy (ignore if already undeployed)
	job, err = cav.Execute(ctx, c.c, undeployVAppOp, getVAppByIDParams{ID: params.ID})
	if err != nil {
		c.logger.Debug("Undeploy failed, continuing with delete", "error", err)
	} else if err := awaitVAppJob(ctx, c.c, opDeleteVApp, job.ID); err != nil {
		c.logger.Debug("Await undeploy failed, continuing with delete", "error", err)
	}

	// Step 3: Delete
	job, err = cav.Execute(ctx, c.c, deleteVAppByIDOp, deleteVAppByIDParams{ID: params.ID})
	if err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVApp, err)
	}
	if err := awaitVAppJob(ctx, c.c, opDeleteVApp, job.ID); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVApp, err)
	}

	return nil
}

// DeployVApp deploys a VApp by ID.
func (c *Client) DeployVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, deployVAppOp, opDeployVApp)
}

// PowerOffVApp powers off a VApp by ID.
func (c *Client) PowerOffVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, powerOffVAppOp, opPowerOffVApp)
}

// PowerOnVApp powers on a VApp by ID.
func (c *Client) PowerOnVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, powerOnVAppOp, opPowerOnVApp)
}

// RebootVApp reboots a VApp by ID.
func (c *Client) RebootVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, rebootVAppOp, opRebootVApp)
}

// ResetVApp resets a VApp by ID.
func (c *Client) ResetVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, resetVAppOp, opResetVApp)
}

// ShutdownVApp shuts down a VApp by ID.
func (c *Client) ShutdownVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, shutdownVAppOp, opShutdownVApp)
}

// SuspendVApp suspends a VApp by ID.
func (c *Client) SuspendVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, suspendVAppOp, opSuspendVApp)
}

func (c *Client) executeVAppAction(ctx context.Context, id string, op cav.Operation[getVAppByIDParams, cav.Job], name string) error {
	job, err := cav.Execute(ctx, c.c, op, getVAppByIDParams{ID: id})
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if err := awaitVAppJob(ctx, c.c, name, job.ID); err != nil {
		return err
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

func toVAppReference(ref types.VAppReference) itypes.APIRequestVAppReference {
	return itypes.APIRequestVAppReference{HREF: ref.HREF, ID: ref.ID, Type: ref.Type, Name: ref.Name}
}

func toInstantiationParams(params *types.VAppInstantiationParams) *itypes.APIRequestVAppInstantiationParams {
	if params == nil {
		return nil
	}

	result := &itypes.APIRequestVAppInstantiationParams{}
	if params.NetworkConfigSection != nil {
		section := &itypes.APIRequestVAppNetworkConfigSection{Info: params.NetworkConfigSection.Info}
		for _, network := range params.NetworkConfigSection.NetworkConfig {
			section.NetworkConfig = append(section.NetworkConfig, itypes.APIRequestVAppNetworkConfig{
				NetworkName: network.NetworkName,
				Network:     network.Network,
			})
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
