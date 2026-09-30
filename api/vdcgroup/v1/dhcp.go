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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetNetworkDhcpConfig    = "VdcNetworkDhcp.Get"
	opUpdateNetworkDhcpConfig = "VdcNetworkDhcp.Update"
	opDeleteNetworkDhcpConfig = "VdcNetworkDhcp.Delete"
)

// GetNetworkDhcpConfig retrieves the DHCP configuration of an Org VDC Network.
func (c *Client) GetNetworkDhcpConfig(ctx context.Context, params types.ParamsGetNetworkDhcpConfig) (*types.ModelDhcpConfig, error) {
	if params.VDCNetworkID == "" {
		return nil, fmt.Errorf("%s: vdcNetworkId is required", opGetNetworkDhcpConfig)
	}

	ep := endpoints.GetNetworkDhcpConfig()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.VDCNetworkID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetNetworkDhcpConfig, err)
	}

	dhcp, ok := resp.Result().(*itypes.DhcpConfig)
	if !ok || dhcp == nil {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetNetworkDhcpConfig, resp.Result())
	}

	model := dhcp.ToModel()
	return &model, nil
}

// UpdateNetworkDhcpConfig updates the DHCP configuration of an Org VDC Network.
func (c *Client) UpdateNetworkDhcpConfig(ctx context.Context, params types.ParamsUpdateNetworkDhcpConfig) (*types.ModelDhcpConfig, error) {
	if params.VDCNetworkID == "" {
		return nil, fmt.Errorf("%s: vdcNetworkId is required", opUpdateNetworkDhcpConfig)
	}

	body := itypes.DhcpConfig{
		ID:                       params.Config.ID,
		VDCNetworkID:             params.Config.VDCNetworkID,
		DHCPServerIPAddress:      params.Config.DHCPServerIPAddress,
		DHCPServerPort:           params.Config.DHCPServerPort,
		DHCPRelayServerIPAddress: params.Config.DHCPRelayServerIPAddress,
		DHCPEnabled:              params.Config.DHCPEnabled,
		DHCPLeaseTime:            params.Config.DHCPLeaseTime,
		DHCPIPAddress:            params.Config.DHCPIPAddress,
		DHCPPoolIPAddress:        params.Config.DHCPPoolIPAddress,
		DHCPSubnetMask:           params.Config.DHCPSubnetMask,
		DHCPDefaultGateway:       params.Config.DHCPDefaultGateway,
		DHCPDNS1IPAddress:        params.Config.DHCPDNS1IPAddress,
		DHCPDNS2IPAddress:        params.Config.DHCPDNS2IPAddress,
		DHCPSearchDomain:         params.Config.DHCPSearchDomain,
	}

	if params.Config.DHCPServerCredentials != nil {
		body.DHCPServerCredentials = &itypes.DHCPServerCredentials{
			Username: params.Config.DHCPServerCredentials.Username,
			Password: params.Config.DHCPServerCredentials.Password,
		}
	}

	body.DHCPOptions = make([]itypes.DHCPOption, 0, len(params.Config.DHCPOptions))
	for _, opt := range params.Config.DHCPOptions {
		body.DHCPOptions = append(body.DHCPOptions, itypes.DHCPOption{Code: opt.Code, Value: opt.Value})
	}

	ep := endpoints.UpdateNetworkDhcpConfig()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.VDCNetworkID), cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateNetworkDhcpConfig, err)
	}

	dhcp, ok := resp.Result().(*itypes.DhcpConfig)
	if !ok || dhcp == nil {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUpdateNetworkDhcpConfig, resp.Result())
	}

	model := dhcp.ToModel()
	return &model, nil
}

// DeleteNetworkDhcpConfig deletes the DHCP configuration of an Org VDC Network.
func (c *Client) DeleteNetworkDhcpConfig(ctx context.Context, params types.ParamsDeleteNetworkDhcpConfig) error {
	if params.VDCNetworkID == "" {
		return fmt.Errorf("%s: vdcNetworkId is required", opDeleteNetworkDhcpConfig)
	}

	ep := endpoints.DeleteNetworkDhcpConfig()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.VDCNetworkID)); err != nil {
		return fmt.Errorf("%s: %w", opDeleteNetworkDhcpConfig, err)
	}

	return nil
}
