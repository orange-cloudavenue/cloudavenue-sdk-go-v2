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

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{
			ID:   vdcGroupID,
			Name: vdcGroupName,
		}},
	}, nil)

	ms.CleanResponse(endpoints.CreateVDCNetwork())
	ms.SetResponse(endpoints.CreateVDCNetwork(), &itypes.APIResponseVDCNetwork{
		ID:                      createdID,
		Name:                    "isolated-net-1",
		Description:             "desc",
		NetworkType:             types.VDCNetworkTypeIsolated,
		OwnerRef:                &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		GuestVLANTaggingAllowed: &guestVLAN,
		Shared:                  &shared,
		Subnets: itypes.APIVDCNetworkSubnets{
			Values: []itypes.APIVDCNetworkSubnetValue{{
				Gateway:      "192.168.10.1",
				PrefixLength: 24,
				DNSServer1:   "1.1.1.1",
				DNSServer2:   "8.8.8.8",
				DNSSuffix:    "example.local",
				IPRanges: itypes.APIVDCNetworkIPRanges{Values: []itypes.APIVDCNetworkIPRangeValue{{
					StartAddress: "192.168.10.10",
					EndAddress:   "192.168.10.20",
				}}},
			}},
		},
	}, nil)

	resp, err := client.CreateVDCNetworkIsolated(t.Context(), types.ParamsCreateVDCNetworkIsolated{
		Name:         "isolated-net-1",
		Description:  "desc",
		VDCGroupName: vdcGroupName,
		Subnet: types.ParamsSubnet{
			Gateway:      "192.168.10.1",
			PrefixLength: 24,
			DNSServer1:   "1.1.1.1",
			DNSServer2:   "8.8.8.8",
			DNSSuffix:    "example.local",
			IPRanges: []types.ParamsVDCNetworkIPRange{{
				StartAddress: "192.168.10.10",
				EndAddress:   "192.168.10.20",
			}},
		},
		GuestVLANTaggingAllowed: &guestVLAN,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createdID, resp.ID)
	assert.Equal(t, "isolated-net-1", resp.Name)
	assert.Equal(t, types.VDCNetworkTypeIsolated, resp.NetworkType)
	assert.Equal(t, vdcGroupID, resp.OwnerID)
	assert.Equal(t, vdcGroupName, resp.OwnerName)
	assert.NotNil(t, resp.GuestVLANTaggingAllowed)
	assert.True(t, *resp.GuestVLANTaggingAllowed)
	assert.NotNil(t, resp.Shared)
	assert.True(t, *resp.Shared)
	assert.Equal(t, "192.168.10.1", resp.Subnet.Gateway)
	assert.Equal(t, 24, resp.Subnet.PrefixLength)
	assert.Len(t, resp.Subnet.IPRanges, 1)
	assert.Equal(t, "192.168.10.10", resp.Subnet.IPRanges[0].StartAddress)
	assert.Equal(t, "192.168.10.20", resp.Subnet.IPRanges[0].EndAddress)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.CreateVDCNetwork())
}

func TestListVdcNetwork(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	guestVLAN := true
	shared := true

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{
			ID:   vdcGroupID,
			Name: vdcGroupName,
		}},
	}, nil)

	ms.CleanResponse(endpoints.ListVDCNetwork())
	ms.SetResponse(endpoints.ListVDCNetwork(), &itypes.APIResponseListVDCNetwork{
		Values: []itypes.APIResponseVDCNetwork{{
			ID:                      generator.MustGenerate("{urn:network}"),
			Name:                    "isolated-net-1",
			NetworkType:             types.VDCNetworkTypeIsolated,
			OwnerRef:                &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
			GuestVLANTaggingAllowed: &guestVLAN,
			Shared:                  &shared,
		}},
	}, nil)

	resp, err := client.ListVDCNetwork(t.Context(), types.ParamsListVDCNetwork{VDCGroupName: vdcGroupName})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.VDCNetworks, 1)
	assert.Equal(t, types.VDCNetworkTypeIsolated, resp.VDCNetworks[0].NetworkType)
	assert.Equal(t, vdcGroupID, resp.VDCNetworks[0].OwnerID)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.ListVDCNetwork())
}

func TestGetVdcNetworkIsolated(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.SetResponse(endpoints.GetVDCNetwork(), &itypes.APIResponseVDCNetwork{
		ID:          networkID,
		Name:        "isolated-net-1",
		NetworkType: types.VDCNetworkTypeIsolated,
		OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	resp, err := client.GetVDCNetworkIsolated(t.Context(), types.ParamsGetVDCNetworkIsolated{ID: networkID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, networkID, resp.ID)
	assert.Equal(t, types.VDCNetworkTypeIsolated, resp.NetworkType)
	assert.Equal(t, vdcGroupID, resp.OwnerID)

	ms.CleanResponse(endpoints.GetVDCNetwork())
}
