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

func TestCreateVdcNetworkRouted(t *testing.T) {
	createdID := generator.MustGenerate("{urn:network}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	edgeGatewayID := generator.MustGenerate("{uuid}")
	edgeGatewayName := generator.MustGenerate("{word}")
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
		Name:                    "routed-net-1",
		Description:             "desc",
		NetworkType:             types.VdcNetworkTypeRouted,
		OwnerRef:                &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		GuestVlanTaggingAllowed: &guestVLAN,
		Shared:                  &shared,
		Connection: &itypes.ApiVdcNetworkConnection{
			RouterRef:           itypes.ApiObjectReference{ID: edgeGatewayID, Name: edgeGatewayName},
			ConnectionTypeValue: "INTERNAL",
		},
		Subnets: itypes.ApiVdcNetworkSubnets{
			Values: []itypes.ApiVdcNetworkSubnetValue{{
				Gateway:      "10.20.30.1",
				PrefixLength: 24,
				DNSServer1:   "9.9.9.9",
				DNSServer2:   "8.8.4.4",
				DNSSuffix:    "corp.local",
				IPRanges: itypes.ApiVdcNetworkIPRanges{Values: []itypes.ApiVdcNetworkIPRangeValue{{
					StartAddress: "10.20.30.10",
					EndAddress:   "10.20.30.20",
				}}},
			}},
		},
	}, nil)

	resp, err := client.CreateVdcNetworkRouted(t.Context(), types.ParamsCreateVdcNetworkRouted{
		Name:            "routed-net-1",
		Description:     "desc",
		VdcGroupName:    vdcGroupName,
		EdgeGatewayID:   edgeGatewayID,
		EdgeGatewayName: edgeGatewayName,
		Subnet: types.ParamsSubnet{
			Gateway:      "10.20.30.1",
			PrefixLength: 24,
			DNSServer1:   "9.9.9.9",
			DNSServer2:   "8.8.4.4",
			DNSSuffix:    "corp.local",
			IPRanges: []types.ParamsVdcNetworkIPRange{{
				StartAddress: "10.20.30.10",
				EndAddress:   "10.20.30.20",
			}},
		},
		GuestVlanTaggingAllowed: &guestVLAN,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createdID, resp.ID)
	assert.Equal(t, "routed-net-1", resp.Name)
	assert.Equal(t, types.VdcNetworkTypeRouted, resp.NetworkType)
	assert.Equal(t, vdcGroupID, resp.OwnerID)
	assert.Equal(t, vdcGroupName, resp.OwnerName)
	assert.Equal(t, edgeGatewayID, resp.EdgeGatewayID)
	assert.Equal(t, edgeGatewayName, resp.EdgeGatewayName)
	assert.NotNil(t, resp.GuestVlanTaggingAllowed)
	assert.True(t, *resp.GuestVlanTaggingAllowed)
	assert.NotNil(t, resp.Shared)
	assert.True(t, *resp.Shared)
	assert.Equal(t, "10.20.30.1", resp.Subnet.Gateway)
	assert.Equal(t, 24, resp.Subnet.PrefixLength)
	assert.Len(t, resp.Subnet.IPRanges, 1)
	assert.Equal(t, "10.20.30.10", resp.Subnet.IPRanges[0].StartAddress)
	assert.Equal(t, "10.20.30.20", resp.Subnet.IPRanges[0].EndAddress)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.CreateVdcNetwork())
}

func TestGetVdcNetworkRouted(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	edgeGatewayID := generator.MustGenerate("{uuid}")
	edgeGatewayName := generator.MustGenerate("{word}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVdcNetwork())
	ms.SetResponse(endpoints.GetVdcNetwork(), &itypes.ApiResponseVdcNetwork{
		ID:          networkID,
		Name:        "routed-net-1",
		NetworkType: types.VdcNetworkTypeRouted,
		OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		Connection: &itypes.ApiVdcNetworkConnection{
			RouterRef: itypes.ApiObjectReference{ID: edgeGatewayID, Name: edgeGatewayName},
		},
	}, nil)

	resp, err := client.GetVdcNetworkRouted(t.Context(), types.ParamsGetVdcNetworkRouted{ID: networkID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, networkID, resp.ID)
	assert.Equal(t, types.VdcNetworkTypeRouted, resp.NetworkType)
	assert.Equal(t, edgeGatewayID, resp.EdgeGatewayID)
	assert.Equal(t, edgeGatewayName, resp.EdgeGatewayName)

	ms.CleanResponse(endpoints.GetVdcNetwork())
}
