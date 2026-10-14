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
	"net/url"
	"path"
	"time"

	"resty.dev/v3"

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

// GetNetworkDhcpConfig retrieves the CloudAPI DHCP configuration of an Org VDC Network.
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

// UpdateNetworkDhcpConfig updates DHCP configuration and waits for its VMware task.
func (c *Client) UpdateNetworkDhcpConfig(ctx context.Context, params types.ParamsUpdateNetworkDhcpConfig) (*types.ModelDhcpConfig, error) {
	if params.VDCNetworkID == "" {
		return nil, fmt.Errorf("%s: vdcNetworkId is required", opUpdateNetworkDhcpConfig)
	}

	body := itypes.APIRequestDhcpConfig{
		Enabled:    params.Config.Enabled,
		LeaseTime:  params.Config.LeaseTime,
		Mode:       params.Config.Mode,
		IPAddress:  params.Config.ListenerIPAddress,
		DNSServers: append([]string(nil), params.Config.DNSServers...),
	}
	for _, pool := range params.Config.Pools {
		body.DhcpPools = append(body.DhcpPools, itypes.DhcpPool{
			IPRange: itypes.IPRange{StartAddress: pool.StartAddress, EndAddress: pool.EndAddress},
		})
	}

	ep := endpoints.UpdateNetworkDhcpConfig()
	resp, err := c.c.Do(
		ctx, ep,
		cav.WithPathParam(ep.PathParams[0], params.VDCNetworkID),
		cav.SetBody(body),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateNetworkDhcpConfig, err)
	}
	if err := awaitDhcpTask(ctx, c.c, resp, opUpdateNetworkDhcpConfig); err != nil {
		return nil, err
	}

	model, err := c.GetNetworkDhcpConfig(ctx, types.ParamsGetNetworkDhcpConfig{VDCNetworkID: params.VDCNetworkID})
	if err != nil {
		return nil, fmt.Errorf("%s: get updated config: %w", opUpdateNetworkDhcpConfig, err)
	}
	return model, nil
}

// DeleteNetworkDhcpConfig deletes DHCP configuration and waits for its VMware task.
func (c *Client) DeleteNetworkDhcpConfig(ctx context.Context, params types.ParamsDeleteNetworkDhcpConfig) error {
	if params.VDCNetworkID == "" {
		return fmt.Errorf("%s: vdcNetworkId is required", opDeleteNetworkDhcpConfig)
	}

	ep := endpoints.DeleteNetworkDhcpConfig()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.VDCNetworkID))
	if err != nil {
		return fmt.Errorf("%s: %w", opDeleteNetworkDhcpConfig, err)
	}
	if err := awaitDhcpTask(ctx, c.c, resp, opDeleteNetworkDhcpConfig); err != nil {
		return err
	}
	return nil
}

func awaitDhcpTask(ctx context.Context, client cav.Client, resp *resty.Response, operation string) error {
	jobID, err := dhcpTaskID(resp)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	if _, err := cav.AwaitJobOnBackend(ctx, client, cav.BackendVMware, jobID, cav.JobPollOptions{
		Timeout:         30 * time.Second,
		PollingInterval: 1 * time.Second,
	}, func(_ *resty.Response) (struct{}, error) {
		return struct{}{}, nil
	}); err != nil {
		return fmt.Errorf("%s: await task: %w", operation, err)
	}
	return nil
}

func dhcpTaskID(resp *resty.Response) (string, error) {
	if resp == nil {
		return "", fmt.Errorf("unexpected async response: missing response")
	}
	location := resp.Header().Get("Location")
	if location == "" {
		return "", fmt.Errorf("unexpected async response: missing task location")
	}
	parsed, err := url.Parse(location)
	if err != nil {
		return "", fmt.Errorf("parse task location: %w", err)
	}
	id := path.Base(parsed.Path)
	if id == "" || id == "." || id == "/" {
		return "", fmt.Errorf("missing task id in location %q", location)
	}
	return id, nil
}
