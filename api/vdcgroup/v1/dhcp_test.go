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
	"net/http"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const dhcpTaskIdentifier = "87ab1934-0146-4fb0-80bc-815fea03214d"

func configureDhcpTask(t *testing.T, ms interface {
	SetResponseFunc(*cav.Endpoint, http.HandlerFunc)
	CleanResponse(*cav.Endpoint)
},
) {
	t.Helper()
	jobEndpoint := cav.MustGetEndpoint("GetJobVmware")
	ms.SetResponseFunc(jobEndpoint, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"` + dhcpTaskIdentifier + `","status":"success"}`))
	})
	t.Cleanup(func() { ms.CleanResponse(jobEndpoint) })
}

func TestGetNetworkDhcpConfig(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	client, ms := newClient(t)
	ms.SetResponse(endpoints.GetNetworkDhcpConfig(), &itypes.DhcpConfig{
		Enabled:    true,
		LeaseTime:  86400,
		Mode:       "NETWORK",
		IPAddress:  "10.0.0.2",
		DNSServers: []string{"8.8.8.8", "1.1.1.1"},
		DhcpPools:  []itypes.DhcpPool{{IPRange: itypes.IPRange{StartAddress: "10.0.0.100", EndAddress: "10.0.0.200"}}},
	}, nil)
	t.Cleanup(func() { ms.CleanResponse(endpoints.GetNetworkDhcpConfig()) })

	resp, err := client.GetNetworkDhcpConfig(t.Context(), types.ParamsGetNetworkDhcpConfig{VDCNetworkID: networkID})

	require.NoError(t, err)
	assert.True(t, resp.Enabled)
	assert.Equal(t, int64(86400), resp.LeaseTime)
	assert.Equal(t, "NETWORK", resp.Mode)
	assert.Equal(t, "10.0.0.2", resp.ListenerIPAddress)
	assert.Equal(t, []string{"8.8.8.8", "1.1.1.1"}, resp.DNSServers)
	require.Len(t, resp.Pools, 1)
	assert.Equal(t, "10.0.0.100", resp.Pools[0].StartAddress)
	assert.Equal(t, "10.0.0.200", resp.Pools[0].EndAddress)
}

func TestUpdateNetworkDhcpConfigWaitsForTask(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	client, ms := newClient(t)
	configureDhcpTask(t, ms)
	ms.SetResponseFunc(endpoints.UpdateNetworkDhcpConfig(), func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", "/api/task/"+dhcpTaskIdentifier)
		w.WriteHeader(http.StatusAccepted)
	})
	ms.SetResponse(endpoints.GetNetworkDhcpConfig(), &itypes.DhcpConfig{Enabled: true, Mode: "NETWORK"}, nil)
	t.Cleanup(func() {
		ms.CleanResponse(endpoints.UpdateNetworkDhcpConfig())
		ms.CleanResponse(endpoints.GetNetworkDhcpConfig())
	})

	resp, err := client.UpdateNetworkDhcpConfig(t.Context(), types.ParamsUpdateNetworkDhcpConfig{
		VDCNetworkID: networkID,
		Config: types.ModelDhcpConfig{
			Enabled:           true,
			Mode:              "NETWORK",
			LeaseTime:         86400,
			ListenerIPAddress: "10.0.0.2",
			Pools:             []types.ModelDhcpPool{{StartAddress: "10.0.0.100", EndAddress: "10.0.0.200"}},
		},
	})

	require.NoError(t, err)
	assert.True(t, resp.Enabled)
	assert.Equal(t, "NETWORK", resp.Mode)
}

func TestDeleteNetworkDhcpConfigWaitsForTask(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	client, ms := newClient(t)
	configureDhcpTask(t, ms)
	ms.SetResponseFunc(endpoints.DeleteNetworkDhcpConfig(), func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", "/api/task/"+dhcpTaskIdentifier)
		w.WriteHeader(http.StatusAccepted)
	})
	t.Cleanup(func() { ms.CleanResponse(endpoints.DeleteNetworkDhcpConfig()) })

	assert.NoError(t, client.DeleteNetworkDhcpConfig(t.Context(), types.ParamsDeleteNetworkDhcpConfig{VDCNetworkID: networkID}))
}

func TestNetworkDhcpConfigValidation(t *testing.T) {
	client, ms := newClient(t)
	t.Cleanup(func() {
		ms.CleanResponse(endpoints.GetNetworkDhcpConfig())
		ms.CleanResponse(endpoints.UpdateNetworkDhcpConfig())
		ms.CleanResponse(endpoints.DeleteNetworkDhcpConfig())
	})

	_, err := client.GetNetworkDhcpConfig(t.Context(), types.ParamsGetNetworkDhcpConfig{})
	assert.Error(t, err)
	_, err = client.UpdateNetworkDhcpConfig(t.Context(), types.ParamsUpdateNetworkDhcpConfig{})
	assert.Error(t, err)
	err = client.DeleteNetworkDhcpConfig(t.Context(), types.ParamsDeleteNetworkDhcpConfig{})
	assert.Error(t, err)
}
