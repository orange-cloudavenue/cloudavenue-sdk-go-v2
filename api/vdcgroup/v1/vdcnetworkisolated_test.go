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

func TestCreateVdcNetworkIsolated(t *testing.T) {
	createdID := generator.MustGenerate("{urn:network}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	guestVLAN := true
	shared := true

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{
			ID:   vdcGroupID,
			Name: vdcGroupName,
		}},
	}, nil)

	ms.CleanResponse(endpoints.CreateVdcNetwork())
	ms.SetResponse(endpoints.CreateVdcNetwork(), &itypes.ApiResponseVdcNetwork{
		ID:                      createdID,
		Name:                    "isolated-net-1",
		Description:             "desc",
		NetworkType:             types.VdcNetworkTypeIsolated,
		OwnerRef:                &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		GuestVlanTaggingAllowed: &guestVLAN,
		Shared:                  &shared,
		Subnets: itypes.ApiVdcNetworkSubnets{
			Values: []itypes.ApiVdcNetworkSubnetValue{{
				Gateway:      "192.168.10.1",
				PrefixLength: 24,
				DNSServer1:   "1.1.1.1",
				DNSServer2:   "8.8.8.8",
				DNSSuffix:    "example.local",
				IPRanges: itypes.ApiVdcNetworkIPRanges{Values: []itypes.ApiVdcNetworkIPRangeValue{{
					StartAddress: "192.168.10.10",
					EndAddress:   "192.168.10.20",
				}}},
			}},
		},
	}, nil)

	resp, err := client.CreateVdcNetworkIsolated(t.Context(), types.ParamsCreateVdcNetworkIsolated{
		Name:         "isolated-net-1",
		Description:  "desc",
		VdcGroupName: vdcGroupName,
		Subnet: types.ParamsSubnet{
			Gateway:      "192.168.10.1",
			PrefixLength: 24,
			DNSServer1:   "1.1.1.1",
			DNSServer2:   "8.8.8.8",
			DNSSuffix:    "example.local",
			IPRanges: []types.ParamsVdcNetworkIPRange{{
				StartAddress: "192.168.10.10",
				EndAddress:   "192.168.10.20",
			}},
		},
		GuestVlanTaggingAllowed: &guestVLAN,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createdID, resp.ID)
	assert.Equal(t, "isolated-net-1", resp.Name)
	assert.Equal(t, types.VdcNetworkTypeIsolated, resp.NetworkType)
	assert.Equal(t, vdcGroupID, resp.OwnerID)
	assert.Equal(t, vdcGroupName, resp.OwnerName)
	assert.NotNil(t, resp.GuestVlanTaggingAllowed)
	assert.True(t, *resp.GuestVlanTaggingAllowed)
	assert.NotNil(t, resp.Shared)
	assert.True(t, *resp.Shared)
	assert.Equal(t, "192.168.10.1", resp.Subnet.Gateway)
	assert.Equal(t, 24, resp.Subnet.PrefixLength)
	assert.Len(t, resp.Subnet.IPRanges, 1)
	assert.Equal(t, "192.168.10.10", resp.Subnet.IPRanges[0].StartAddress)
	assert.Equal(t, "192.168.10.20", resp.Subnet.IPRanges[0].EndAddress)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.CreateVdcNetwork())
}

func TestListVdcNetwork(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	guestVLAN := true
	shared := true

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{
			ID:   vdcGroupID,
			Name: vdcGroupName,
		}},
	}, nil)

	ms.CleanResponse(endpoints.ListVdcNetwork())
	ms.SetResponse(endpoints.ListVdcNetwork(), &itypes.ApiResponseListVdcNetwork{
		Values: []itypes.ApiResponseVdcNetwork{{
			ID:                      generator.MustGenerate("{urn:network}"),
			Name:                    "isolated-net-1",
			NetworkType:             types.VdcNetworkTypeIsolated,
			OwnerRef:                &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
			GuestVlanTaggingAllowed: &guestVLAN,
			Shared:                  &shared,
		}},
	}, nil)

	resp, err := client.ListVdcNetwork(t.Context(), types.ParamsListVdcNetwork{VdcGroupName: vdcGroupName})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.VdcNetworks, 1)
	assert.Equal(t, types.VdcNetworkTypeIsolated, resp.VdcNetworks[0].NetworkType)
	assert.Equal(t, vdcGroupID, resp.VdcNetworks[0].OwnerID)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.ListVdcNetwork())
}

func TestGetVdcNetworkIsolated(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVdcNetwork())
	ms.SetResponse(endpoints.GetVdcNetwork(), &itypes.ApiResponseVdcNetwork{
		ID:          networkID,
		Name:        "isolated-net-1",
		NetworkType: types.VdcNetworkTypeIsolated,
		OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	resp, err := client.GetVdcNetworkIsolated(t.Context(), types.ParamsGetVdcNetworkIsolated{ID: networkID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, networkID, resp.ID)
	assert.Equal(t, types.VdcNetworkTypeIsolated, resp.NetworkType)
	assert.Equal(t, vdcGroupID, resp.OwnerID)

	ms.CleanResponse(endpoints.GetVdcNetwork())
}
