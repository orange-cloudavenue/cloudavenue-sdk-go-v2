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
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestGetNetworkDhcpConfig(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetNetworkDhcpConfig())
	ms.SetResponse(endpoints.GetNetworkDhcpConfig(), &itypes.DhcpConfig{
		ID:                       "dhcp-1",
		VDCNetworkID:             networkID,
		DHCPServerIPAddress:      "10.0.0.10",
		DHCPServerPort:           5460,
		DHCPRelayServerIPAddress: "10.0.0.1",
		DHCPEnabled:              true,
		DHCPLeaseTime:            3600,
		DHCPIPAddress:            "10.0.0.0/24",
		DHCPPoolIPAddress:        "10.0.0.100",
		DHCPSubnetMask:           "255.255.255.0",
		DHCPDefaultGateway:       "10.0.0.1",
		DHCPDNS1IPAddress:        "8.8.8.8",
		DHCPDNS2IPAddress:        "8.8.4.4",
		DHCPSearchDomain:         "corp.local",
		DHCPServerCredentials: &itypes.DHCPServerCredentials{
			Username: "dhcpuser",
			Password: "dhcppass",
		},
		DHCPOptions: []itypes.DHCPOption{
			{Code: "6", Value: "8.8.8.8"},
		},
	}, nil)

	resp, err := client.GetNetworkDhcpConfig(t.Context(), types.ParamsGetNetworkDhcpConfig{VDCNetworkID: networkID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "dhcp-1", resp.ID)
	assert.Equal(t, networkID, resp.VDCNetworkID)
	assert.Equal(t, "10.0.0.10", resp.DHCPServerIPAddress)
	assert.Equal(t, 5460, resp.DHCPServerPort)
	assert.Equal(t, "10.0.0.1", resp.DHCPRelayServerIPAddress)
	assert.True(t, resp.DHCPEnabled)
	assert.Equal(t, 3600, resp.DHCPLeaseTime)
	assert.Equal(t, "10.0.0.0/24", resp.DHCPIPAddress)
	assert.Equal(t, "10.0.0.100", resp.DHCPPoolIPAddress)
	assert.Equal(t, "255.255.255.0", resp.DHCPSubnetMask)
	assert.Equal(t, "10.0.0.1", resp.DHCPDefaultGateway)
	assert.Equal(t, "8.8.8.8", resp.DHCPDNS1IPAddress)
	assert.Equal(t, "8.8.4.4", resp.DHCPDNS2IPAddress)
	assert.Equal(t, "corp.local", resp.DHCPSearchDomain)
	assert.NotNil(t, resp.DHCPServerCredentials)
	assert.Equal(t, "dhcpuser", resp.DHCPServerCredentials.Username)
	assert.Equal(t, "dhcppass", resp.DHCPServerCredentials.Password)
	assert.Len(t, resp.DHCPOptions, 1)
	assert.Equal(t, "6", resp.DHCPOptions[0].Code)
	assert.Equal(t, "8.8.8.8", resp.DHCPOptions[0].Value)

	ms.CleanResponse(endpoints.GetNetworkDhcpConfig())
}

func TestGetNetworkDhcpConfig_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetNetworkDhcpConfig())

	_, err := client.GetNetworkDhcpConfig(t.Context(), types.ParamsGetNetworkDhcpConfig{})
	assert.Error(t, err)
}

func TestGetNetworkDhcpConfig_Error(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetNetworkDhcpConfig())
	status := 500
	ms.SetResponse(endpoints.GetNetworkDhcpConfig(), nil, &status)

	_, err := client.GetNetworkDhcpConfig(t.Context(), types.ParamsGetNetworkDhcpConfig{VDCNetworkID: networkID})
	assert.Error(t, err)

	ms.CleanResponse(endpoints.GetNetworkDhcpConfig())
}

func TestUpdateNetworkDhcpConfig(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.UpdateNetworkDhcpConfig())
	ms.SetResponse(endpoints.UpdateNetworkDhcpConfig(), &itypes.DhcpConfig{
		ID:           "dhcp-1",
		VDCNetworkID: networkID,
		DHCPEnabled:  true,
	}, nil)

	resp, err := client.UpdateNetworkDhcpConfig(t.Context(), types.ParamsUpdateNetworkDhcpConfig{
		VDCNetworkID: networkID,
		Config: types.ModelDhcpConfig{
			ID:                  "dhcp-1",
			VDCNetworkID:        networkID,
			DHCPServerIPAddress: "10.0.0.10",
			DHCPServerPort:      5460,
			DHCPEnabled:         true,
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "dhcp-1", resp.ID)
	assert.Equal(t, networkID, resp.VDCNetworkID)
	assert.True(t, resp.DHCPEnabled)

	ms.CleanResponse(endpoints.UpdateNetworkDhcpConfig())
}

func TestUpdateNetworkDhcpConfig_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.UpdateNetworkDhcpConfig())

	_, err := client.UpdateNetworkDhcpConfig(t.Context(), types.ParamsUpdateNetworkDhcpConfig{})
	assert.Error(t, err)
}

func TestDeleteNetworkDhcpConfig(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	client, ms := newClient(t)

	err := client.DeleteNetworkDhcpConfig(t.Context(), types.ParamsDeleteNetworkDhcpConfig{VDCNetworkID: networkID})
	assert.NoError(t, err)

	ms.CleanResponse(endpoints.DeleteNetworkDhcpConfig())
}

func TestDeleteNetworkDhcpConfig_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.DeleteNetworkDhcpConfig())

	err := client.DeleteNetworkDhcpConfig(t.Context(), types.ParamsDeleteNetworkDhcpConfig{})
	assert.Error(t, err)
}
