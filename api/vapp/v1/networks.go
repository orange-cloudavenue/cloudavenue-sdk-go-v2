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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetVAppNetworkConfig = "VApp.NetworkConfig.Get"
	opCreateVAppNetwork    = "VApp.NetworkConfig.Create"
	opAddOrgNetworkToVApp  = "VApp.NetworkConfig.AddOrgNetwork"
	opUpdateVAppNetwork    = "VApp.NetworkConfig.Update"
	opRemoveVAppNetwork    = "VApp.NetworkConfig.Remove"
)

var (
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
)

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
		section.OVF = ovfEnvelopeXMLNS
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
		ref := itypes.APIRequestVAppReference(*src.ParentNetwork)
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
