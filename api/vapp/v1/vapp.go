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
	stderrors "errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	ovfEnvelopeXMLNS             = "http://schemas.dmtf.org/ovf/envelope/1"
	opListVApp                   = "VApp.List"
	opGetVApp                    = "VApp.Get"
	opGetVAppNetworkConfig       = "VApp.NetworkConfig.Get"
	opGetVAppACL                 = "VApp.ACL.Get"
	opSetVAppACL                 = "VApp.ACL.Set"
	opRemoveVAppACL              = "VApp.ACL.Remove"
	opGetVAppGuestProperties     = "VApp.GuestProperties.Get"
	opUpdateVAppGuestProperties  = "VApp.GuestProperties.Update"
	opCreateVAppNetwork          = "VApp.NetworkConfig.Create"
	opAddOrgNetworkToVApp        = "VApp.NetworkConfig.AddOrgNetwork"
	opUpdateVAppNetwork          = "VApp.NetworkConfig.Update"
	opRemoveVAppNetwork          = "VApp.NetworkConfig.Remove"
	opComposeVApp                = "VApp.Compose"
	opInstantiateVAppTemplate    = "VApp.InstantiateTemplate"
	opUpdateVApp                 = "VApp.Update"
	opDeleteVApp                 = "VApp.Delete"
	opUndeployVApp               = "VApp.Undeploy"
	opDeployVApp                 = "VApp.Deploy"
	opPowerOffVApp               = "VApp.PowerOff"
	opPowerOnVApp                = "VApp.PowerOn"
	opRebootVApp                 = "VApp.Reboot"
	opResetVApp                  = "VApp.Reset"
	opShutdownVApp               = "VApp.Shutdown"
	opSuspendVApp                = "VApp.Suspend"
	opDiscardSuspendedStateVApp  = "VApp.DiscardSuspendedState"
	vcloudXMLNS                  = "http://www.vmware.com/vcloud/v1.5"
	xmlSchemaInstanceNS          = "http://www.w3.org/2001/XMLSchema-instance"
	pathAdminUserByID            = "/api/admin/user/%s"
	pathAdminGroupByID           = "/api/admin/group/%s"
	urnOrgUserPrefix             = "urn:vcloud:orguser:"
	urnOrgGroupPrefix            = "urn:vcloud:orggroup:"
	urnOrgUserType               = "urn:vcloud:orguser"
	urnOrgGroupType              = "urn:vcloud:orggroup"
	vappACLSubjectKindUser       = "user"
	vappACLSubjectKindGroup      = "group"
	vappACLSubjectSourceLocal    = "local"
	vappACLSubjectSourceExternal = "external"
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

type updateVAppNetworkConfigSectionParams struct {
	ID   string
	Body itypes.APIResponseVAppNetworkConfigSection
}

type updateVAppGuestPropertiesParams struct {
	ID   string
	Body itypes.APIResponseGuestProperties
}

type setVAppACLParams struct {
	ID   string
	Body itypes.ControlAccessParams
}

type undeployVAppParams struct {
	ID   string
	Body itypes.APIRequestUndeployVApp
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
	getVAppNetworkConfigSectionOp = cav.Operation[getVAppByIDParams, *itypes.APIResponseVAppNetworkConfigSection]{
		Name:     opGetVAppNetworkConfig,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.GetVAppNetworkConfigSection(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.GetVAppNetworkConfigSection()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (*itypes.APIResponseVAppNetworkConfigSection, error) {
			section, ok := resp.Result().(*itypes.APIResponseVAppNetworkConfigSection)
			if !ok || section == nil {
				return nil, fmt.Errorf("unexpected get response type %T", resp.Result())
			}
			return section, nil
		},
	}
	getVAppACLOp = cav.Operation[getVAppByIDParams, *itypes.ControlAccessParams]{
		Name:     opGetVAppACL,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.GetVAppACL(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.GetVAppACL()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (*itypes.ControlAccessParams, error) {
			acl, ok := resp.Result().(*itypes.ControlAccessParams)
			if !ok || acl == nil {
				return nil, fmt.Errorf("unexpected get response type %T", resp.Result())
			}
			return acl, nil
		},
	}
	setVAppACLOp = cav.Operation[setVAppACLParams, *itypes.ControlAccessParams]{
		Name:     opSetVAppACL,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.SetVAppACL(),
		Validate: func(p setVAppACLParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p setVAppACLParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.SetVAppACL()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p setVAppACLParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ setVAppACLParams) (*itypes.ControlAccessParams, error) {
			acl, ok := resp.Result().(*itypes.ControlAccessParams)
			if !ok || acl == nil {
				return nil, fmt.Errorf("unexpected set response type %T", resp.Result())
			}
			return acl, nil
		},
	}
	getVAppProductSectionsOp = cav.Operation[getVAppByIDParams, *itypes.APIResponseGuestProperties]{
		Name:     opGetVAppGuestProperties,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.GetVAppProductSections(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.GetVAppProductSections()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (*itypes.APIResponseGuestProperties, error) {
			section, ok := resp.Result().(*itypes.APIResponseGuestProperties)
			if !ok || section == nil {
				return nil, fmt.Errorf("unexpected get response type %T", resp.Result())
			}
			return section, nil
		},
	}
	updateVAppNetworkConfigSectionOp = cav.Operation[updateVAppNetworkConfigSectionParams, cav.Job]{
		Name:     "VApp.NetworkConfig.Replace",
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UpdateVAppNetworkConfigSection(),
		Validate: func(p updateVAppNetworkConfigSectionParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p updateVAppNetworkConfigSectionParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UpdateVAppNetworkConfigSection()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p updateVAppNetworkConfigSectionParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ updateVAppNetworkConfigSectionParams) (cav.Job, error) {
			return extractVMwareJob(resp, "replace network config section")
		},
	}
	updateVAppGuestPropertiesOp = cav.Operation[updateVAppGuestPropertiesParams, cav.Job]{
		Name:     opUpdateVAppGuestProperties,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UpdateVAppProductSections(),
		Validate: func(p updateVAppGuestPropertiesParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p updateVAppGuestPropertiesParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UpdateVAppProductSections()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p updateVAppGuestPropertiesParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ updateVAppGuestPropertiesParams) (cav.Job, error) {
			return extractVMwareJob(resp, "replace guest properties")
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
	undeployVAppOp = cav.Operation[undeployVAppParams, cav.Job]{
		Name:     opUndeployVApp,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.UndeployVApp(),
		Validate: func(p undeployVAppParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p undeployVAppParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UndeployVApp()
			return []cav.EndpointRequestOption{
				cav.WithPathParam(ep.PathParams[0], p.ID),
				cav.SetCustomRestyOption(func(req *resty.Request) {
					req.SetHeader("Content-Type", "application/xml")
				}),
			}, nil
		},
		Transform: func(p undeployVAppParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ undeployVAppParams) (cav.Job, error) {
			return extractVMwareJob(resp, "undeploy")
		},
	}
	deployVAppOp                = newVAppActionOperation(opDeployVApp, endpoints.DeployVApp, "deploy")
	powerOffVAppOp              = newVAppActionOperation(opPowerOffVApp, endpoints.PowerOffVApp, "power off")
	powerOnVAppOp               = newVAppActionOperation(opPowerOnVApp, endpoints.PowerOnVApp, "power on")
	rebootVAppOp                = newVAppActionOperation(opRebootVApp, endpoints.RebootVApp, "reboot")
	resetVAppOp                 = newVAppActionOperation(opResetVApp, endpoints.ResetVApp, "reset")
	shutdownVAppOp              = newVAppActionOperation(opShutdownVApp, endpoints.ShutdownVApp, "shutdown")
	suspendVAppOp               = newVAppActionOperation(opSuspendVApp, endpoints.SuspendVApp, "suspend")
	discardSuspendedStateVAppOp = newVAppActionOperation(opDiscardSuspendedStateVApp, endpoints.DiscardSuspendedStateVApp, "discard suspended state")
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

// GetVAppNetworkConfig returns the runtime network configuration section for a VApp.
func (c *Client) GetVAppNetworkConfig(ctx context.Context, params types.ParamsGetVAppNetworkConfig) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetVAppNetworkConfig, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppNetworkConfig, err)
	}

	resp, err := cav.Execute(ctx, c.c, getVAppNetworkConfigSectionOp, getVAppByIDParams{ID: vappID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppNetworkConfig, err)
	}

	return resp.ToModel(), nil
}

// GetVAppACL returns the full ACL state for a vApp.
func (c *Client) GetVAppACL(ctx context.Context, params types.ParamsGetVAppACL) (*types.ModelVAppACL, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetVAppACL, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppACL, err)
	}

	resp, err := cav.Execute(ctx, c.c, getVAppACLOp, getVAppByIDParams{ID: vappID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppACL, err)
	}

	return resp.ToVAppACLModel(), nil
}

// SetVAppACL replaces the full ACL state for a vApp.
func (c *Client) SetVAppACL(ctx context.Context, params types.ParamsSetVAppACL) (*types.ModelVAppACL, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opSetVAppACL, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opSetVAppACL, err)
	}

	resp, err := cav.Execute(ctx, c.c, setVAppACLOp, setVAppACLParams{ID: vappID, Body: toVAppControlAccessParams(c.c.GetConsole().GetAPIVCDEndpoint(), params)})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opSetVAppACL, err)
	}

	return resp.ToVAppACLModel(), nil
}

// RemoveVAppACL clears all explicit ACL grants on a vApp.
// It replaces the controlAccess state with IsSharedToEveryone=false
// and an empty subject list, matching govcd RemoveAccessControl semantics.
func (c *Client) RemoveVAppACL(ctx context.Context, params types.ParamsRemoveVAppACL) (*types.ModelVAppACL, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opRemoveVAppACL, err)
	}

	return c.SetVAppACL(ctx, types.ParamsSetVAppACL{ID: params.ID, Name: params.Name})
}

// GetVAppGuestProperties returns the ProductSection guest properties for a VApp.
func (c *Client) GetVAppGuestProperties(ctx context.Context, params types.ParamsGetVAppGuestProperties) (*types.ModelVAppGuestProperties, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetVAppGuestProperties, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppGuestProperties, err)
	}

	resp, err := cav.Execute(ctx, c.c, getVAppProductSectionsOp, getVAppByIDParams{ID: vappID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppGuestProperties, err)
	}

	model := resp.ToModel()
	return &model, nil
}

// UpdateVAppGuestProperties replaces the ProductSection guest properties for a VApp.
func (c *Client) UpdateVAppGuestProperties(ctx context.Context, params types.ParamsUpdateVAppGuestProperties) (*types.ModelVAppGuestProperties, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateVAppGuestProperties, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVAppGuestProperties, err)
	}

	body := toAPIResponseGuestProperties(params.GuestProperties)
	job, err := cav.Execute(ctx, c.c, updateVAppGuestPropertiesOp, updateVAppGuestPropertiesParams{ID: vappID, Body: body})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVAppGuestProperties, err)
	}
	if err := awaitVAppJob(ctx, c.c, opUpdateVAppGuestProperties, job.ID); err != nil {
		return nil, err
	}
	return c.GetVAppGuestProperties(ctx, types.ParamsGetVAppGuestProperties{ID: vappID})
}

// CreateVAppNetwork creates a new isolated runtime vApp network and re-reads the section.
func (c *Client) CreateVAppNetwork(ctx context.Context, params types.ParamsCreateVAppNetwork) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateVAppNetwork, err)
	}

	return c.mutateVAppNetworkConfig(ctx, opCreateVAppNetwork, params.VAppID, params.VAppName, func(section *itypes.APIResponseVAppNetworkConfigSection) error {
		if findRawVAppRuntimeNetwork(section, params.Network.NetworkName) != nil {
			return fmt.Errorf("network %q already exists", params.Network.NetworkName)
		}
		section.NetworkConfig = append(section.NetworkConfig, toRawVAppRuntimeNetwork(params.Network))
		return nil
	})
}

// AddOrgNetworkToVApp adds an org-backed runtime vApp network and re-reads the section.
func (c *Client) AddOrgNetworkToVApp(ctx context.Context, params types.ParamsAddOrgNetworkToVApp) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opAddOrgNetworkToVApp, err)
	}

	return c.mutateVAppNetworkConfig(ctx, opAddOrgNetworkToVApp, params.VAppID, params.VAppName, func(section *itypes.APIResponseVAppNetworkConfigSection) error {
		if findRawVAppRuntimeNetwork(section, params.Network.NetworkName) != nil {
			return fmt.Errorf("network %q already exists", params.Network.NetworkName)
		}
		section.NetworkConfig = append(section.NetworkConfig, toRawVAppRuntimeNetwork(params.Network))
		return nil
	})
}

// UpdateVAppNetwork updates an existing runtime vApp network and re-reads the section.
func (c *Client) UpdateVAppNetwork(ctx context.Context, params types.ParamsUpdateVAppNetwork) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateVAppNetwork, err)
	}

	return c.mutateVAppNetworkConfig(ctx, opUpdateVAppNetwork, params.VAppID, params.VAppName, func(section *itypes.APIResponseVAppNetworkConfigSection) error {
		network := findRawVAppRuntimeNetwork(section, params.NetworkName)
		if network == nil {
			return fmt.Errorf("network %q does not exist", params.NetworkName)
		}
		updated := params.Network
		updated.NetworkName = params.NetworkName
		patchRawVAppRuntimeNetwork(network, updated)
		return nil
	})
}

// RemoveVAppNetwork removes an existing runtime vApp network and re-reads the section.
func (c *Client) RemoveVAppNetwork(ctx context.Context, params types.ParamsRemoveVAppNetwork) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opRemoveVAppNetwork, err)
	}

	return c.mutateVAppNetworkConfig(ctx, opRemoveVAppNetwork, params.VAppID, params.VAppName, func(section *itypes.APIResponseVAppNetworkConfigSection) error {
		index := findRawVAppRuntimeNetworkIndex(section, params.NetworkName)
		if index < 0 {
			return fmt.Errorf("network %q does not exist", params.NetworkName)
		}
		section.NetworkConfig = append(section.NetworkConfig[:index], section.NetworkConfig[index+1:]...)
		return nil
	})
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

	// Step 2: Try undeploy (ignore only the explicit already-undeployed case)
	job, err = cav.Execute(ctx, c.c, undeployVAppOp, undeployVAppParams{ID: params.ID, Body: itypes.ToAPIRequestUndeployVApp(types.VAppUndeployModePowerOff)})
	if err != nil {
		if !isIgnorableUndeployError(err) {
			return fmt.Errorf("%s: undeploy: %w", opDeleteVApp, err)
		}
		c.logger.Debug("Undeploy skipped during delete because VApp is already undeployed", "error", err)
	} else if err := awaitVAppJob(ctx, c.c, opDeleteVApp, job.ID); err != nil {
		if !isIgnorableUndeployError(err) {
			return fmt.Errorf("%s: undeploy: %w", opDeleteVApp, err)
		}
		c.logger.Debug("Undeploy completion skipped during delete because VApp is already undeployed", "error", err)
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

// UndeployVApp undeploys a VApp using the requested VMware undeploy mode.
func (c *Client) UndeployVApp(ctx context.Context, params types.ParamsUndeployVApp) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opUndeployVApp, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return fmt.Errorf("%s: %w", opUndeployVApp, err)
	}

	job, err := cav.Execute(ctx, c.c, undeployVAppOp, undeployVAppParams{ID: vappID, Body: itypes.ToAPIRequestUndeployVApp(params.Mode)})
	if err != nil {
		return fmt.Errorf("%s: %w", opUndeployVApp, err)
	}
	if err := awaitVAppJob(ctx, c.c, opUndeployVApp, job.ID); err != nil {
		return err
	}
	return nil
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

// DiscardSuspendedStateVApp discards the suspended state of a VApp by ID.
func (c *Client) DiscardSuspendedStateVApp(ctx context.Context, id string) error {
	return c.executeVAppAction(ctx, id, discardSuspendedStateVAppOp, opDiscardSuspendedStateVApp)
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

func toVAppReference(ref types.VAppReference) itypes.APIRequestVAppReference {
	return itypes.APIRequestVAppReference{HREF: ref.HREF, ID: ref.ID, Type: ref.Type, Name: ref.Name}
}

func toAPIResponseGuestProperties(model types.ModelVAppGuestProperties) itypes.APIResponseGuestProperties {
	result := itypes.APIResponseGuestProperties{
		XMLNS:                        vcloudXMLNS,
		XMLNSXsi:                     xmlSchemaInstanceNS,
		XsiNoNamespaceSchemaLocation: "http://www.vmware.com/vcloud/v1.5/vmguest.xsd",
		ProductSectionList:           make([]itypes.APIResponseProductSection, 0, len(model.Sections)),
	}
	sections := append([]types.ModelVAppGuestPropertySection(nil), model.Sections...)
	sort.SliceStable(sections, func(i, j int) bool {
		if sections[i].Class != sections[j].Class {
			return sections[i].Class < sections[j].Class
		}
		return sections[i].Instance < sections[j].Instance
	})

	for _, section := range sections {
		keys := make([]string, 0, len(section.Properties))
		for key := range section.Properties {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		properties := make([]itypes.APIResponseProductProperty, 0, len(section.Properties))
		for _, key := range keys {
			properties = append(properties, itypes.APIResponseProductProperty{Key: key, Value: section.Properties[key]})
		}
		result.ProductSectionList = append(result.ProductSectionList, itypes.APIResponseProductSection{
			Class:    section.Class,
			Instance: section.Instance,
			Property: properties,
		})
	}
	return result
}

func isIgnorableUndeployError(err error) bool {
	var apiErr *pkgerrors.APIError
	if !stderrors.As(err, &apiErr) || apiErr == nil {
		return false
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		return false
	}

	message := strings.ToLower(strings.TrimSpace(apiErr.Message))
	return strings.Contains(message, "requested operation could not be executed") &&
		strings.Contains(message, "vapp") &&
		strings.Contains(message, "not running")
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

func toVAppControlAccessParams(apiEndpoint string, params types.ParamsSetVAppACL) itypes.ControlAccessParams {
	result := itypes.ControlAccessParams{
		Xmlns:               vcloudXMLNS,
		IsSharedToEveryone:  params.IsSharedToEveryone,
		EveryoneAccessLevel: params.EveryoneAccessLevel,
	}

	if len(params.SharedWith) == 0 {
		return result
	}

	settings := make([]*itypes.AccessSetting, 0, len(params.SharedWith))
	for _, subject := range params.SharedWith {
		setting := &itypes.AccessSetting{AccessLevel: subject.AccessLevel}
		switch subject.SubjectSource {
		case vappACLSubjectSourceLocal:
			setting.Subject = &itypes.LocalSubject{
				HREF: toVAppACLSubjectHREF(apiEndpoint, subject),
				Type: toVAppACLSubjectType(subject.SubjectKind),
			}
		case vappACLSubjectSourceExternal:
			setting.ExternalSubject = &itypes.ExternalSubject{
				IsUser:    subject.SubjectKind == vappACLSubjectKindUser,
				SubjectID: subject.SubjectID,
			}
		}
		settings = append(settings, setting)
	}

	result.AccessSettings = &itypes.AccessSettingList{AccessSetting: settings}
	return result
}

func toVAppACLSubjectHREF(apiEndpoint string, subject types.ParamsSetVAppACLSubject) string {
	if subject.SubjectSource != vappACLSubjectSourceLocal {
		return ""
	}
	switch subject.SubjectKind {
	case vappACLSubjectKindUser:
		return localSubjectHREFFromURN(apiEndpoint, subject.SubjectID, urnOrgUserPrefix, pathAdminUserByID)
	case vappACLSubjectKindGroup:
		return localSubjectHREFFromURN(apiEndpoint, subject.SubjectID, urnOrgGroupPrefix, pathAdminGroupByID)
	default:
		return subject.SubjectID
	}
}

func toVAppACLSubjectType(subjectKind string) string {
	switch subjectKind {
	case vappACLSubjectKindUser:
		return urnOrgUserType
	case vappACLSubjectKindGroup:
		return urnOrgGroupType
	default:
		return ""
	}
}

func localSubjectHREFFromURN(apiEndpoint, subjectID, urnPrefix, pathTemplate string) string {
	if subjectID == "" {
		return ""
	}
	if after, ok := strings.CutPrefix(subjectID, urnPrefix); ok {
		return apiEndpoint + fmt.Sprintf(pathTemplate, after)
	}
	return subjectID
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

func (c *Client) mutateVAppNetworkConfig(ctx context.Context, opName, vappID, vappName string, mutate func(section *itypes.APIResponseVAppNetworkConfigSection) error) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
	resolvedVAppID, err := c.resolveVAppID(ctx, vappID, vappName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opName, err)
	}

	section, err := cav.Execute(ctx, c.c, getVAppNetworkConfigSectionOp, getVAppByIDParams{ID: resolvedVAppID})
	if err != nil {
		return nil, fmt.Errorf("%s: get current: %w", opName, err)
	}
	if section == nil {
		section = &itypes.APIResponseVAppNetworkConfigSection{}
	}
	if section.Info == "" {
		section.Info = "Runtime network config"
	}
	if section.XMLNS == "" {
		section.XMLNS = vcloudXMLNS
	}
	if section.OVF == "" {
		section.OVF = "http://schemas.dmtf.org/ovf/envelope/1"
	}
	if err := mutate(section); err != nil {
		return nil, fmt.Errorf("%s: %w", opName, err)
	}
	job, err := cav.Execute(ctx, c.c, updateVAppNetworkConfigSectionOp, updateVAppNetworkConfigSectionParams{ID: resolvedVAppID, Body: *section})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opName, err)
	}
	if err := awaitVAppJob(ctx, c.c, opName, job.ID); err != nil {
		return nil, err
	}
	updated, err := c.GetVAppNetworkConfig(ctx, types.ParamsGetVAppNetworkConfig{ID: resolvedVAppID})
	if err != nil {
		return nil, fmt.Errorf("%s: get updated: %w", opName, err)
	}
	return updated, nil
}

func findRawVAppRuntimeNetwork(section *itypes.APIResponseVAppNetworkConfigSection, name string) *itypes.APIResponseVAppNetworkConfig {
	index := findRawVAppRuntimeNetworkIndex(section, name)
	if index < 0 {
		return nil
	}
	return &section.NetworkConfig[index]
}

func findRawVAppRuntimeNetworkIndex(section *itypes.APIResponseVAppNetworkConfigSection, name string) int {
	for i := range section.NetworkConfig {
		if section.NetworkConfig[i].NetworkName == name {
			return i
		}
	}
	return -1
}

func patchRawVAppRuntimeNetwork(dst *itypes.APIResponseVAppNetworkConfig, src types.ModelVAppRuntimeNetworkConfig) {
	dst.NetworkName = src.NetworkName
	dst.Description = src.Description
	patchRawVAppRuntimeConfiguration(&dst.Configuration, src.Configuration)
}

func patchRawVAppRuntimeConfiguration(dst *itypes.APIResponseVAppRuntimeConfiguration, src types.ModelVAppRuntimeNetworkConfiguration) {
	dst.FenceMode = src.FenceMode
	dst.RetainNetInfoAcrossDeployments = src.RetainNetInfoAcrossDeployments
	dst.IsDeployed = src.IsDeployed
	dst.IsShared = src.IsShared
	dst.IsConnected = &src.IsConnected
	if src.ParentNetwork == nil {
		dst.ParentNetwork = nil
	} else {
		ref := toVAppReference(*src.ParentNetwork)
		dst.ParentNetwork = &ref
	}
	dst.IPScopes = toRawVAppRuntimeIPScopes(src.IPScopes)
}

func toRawVAppRuntimeIPScopes(scopes []types.ModelVAppRuntimeIPScope) *itypes.APIResponseVAppRuntimeIPScopes {
	if len(scopes) == 0 {
		return nil
	}
	result := &itypes.APIResponseVAppRuntimeIPScopes{IPScope: make([]itypes.APIResponseVAppRuntimeIPScope, 0, len(scopes))}
	for _, scope := range scopes {
		isInherited := scope.IsInherited
		result.IPScope = append(result.IPScope, itypes.APIResponseVAppRuntimeIPScope{
			IsInherited:  &isInherited,
			Gateway:      scope.Gateway,
			Netmask:      scope.Netmask,
			DNS1:         scope.DNS1,
			DNS2:         scope.DNS2,
			DNSSuffix:    scope.DNSSuffix,
			PrefixLength: scope.PrefixLength,
			IPRanges:     toRawVAppRuntimeIPRanges(scope.IPRanges),
		})
	}
	return result
}

func toRawVAppRuntimeIPRanges(ranges []types.ModelVAppRuntimeIPRange) itypes.APIResponseVAppRuntimeIPRanges {
	result := itypes.APIResponseVAppRuntimeIPRanges{IPRange: make([]itypes.APIResponseVAppRuntimeIPRange, 0, len(ranges))}
	for _, ipRange := range ranges {
		result.IPRange = append(result.IPRange, itypes.APIResponseVAppRuntimeIPRange{StartAddress: ipRange.StartAddress, EndAddress: ipRange.EndAddress})
	}
	return result
}

func toRawVAppRuntimeNetwork(model types.ModelVAppRuntimeNetworkConfig) itypes.APIResponseVAppNetworkConfig {
	result := itypes.APIResponseVAppNetworkConfig{
		Description: model.Description,
		NetworkName: model.NetworkName,
	}
	patchRawVAppRuntimeConfiguration(&result.Configuration, model.Configuration)
	return result
}
