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

func TestUpdateVdcNetworkIsolated(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	guestVLAN := true
	shared := true

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.SetResponse(endpoints.GetVDCNetwork(), &itypes.APIResponseVDCNetwork{
		ID:                      networkID,
		Name:                    "isolated-net-1",
		Description:             "old description",
		NetworkType:             types.VDCNetworkTypeIsolated,
		OwnerRef:                &itypes.APIObjectReference{ID: vdcGroupID, Name: "vdcg1"},
		GuestVLANTaggingAllowed: &guestVLAN,
		Shared:                  &shared,
		Subnets: itypes.APIVDCNetworkSubnets{
			Values: []itypes.APIVDCNetworkSubnetValue{{
				Gateway:      "192.168.10.1",
				PrefixLength: 24,
				IPRanges: itypes.APIVDCNetworkIPRanges{Values: []itypes.APIVDCNetworkIPRangeValue{{
					StartAddress: "192.168.10.10",
					EndAddress:   "192.168.10.20",
				}}},
			}},
		},
	}, nil)

	ms.CleanResponse(endpoints.UpdateVDCNetwork())

	resp, err := client.UpdateVDCNetworkIsolated(t.Context(), types.ParamsUpdateVDCNetworkIsolated{
		ID:          networkID,
		Description: "new description",
		Subnet: &types.ParamsSubnet{
			Gateway:      "10.0.0.1",
			PrefixLength: 24,
			IPRanges: []types.ParamsVDCNetworkIPRange{{
				StartAddress: "10.0.0.10",
				EndAddress:   "10.0.0.20",
			}},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, networkID, resp.ID)
	assert.Equal(t, "new description", resp.Description)
	assert.Equal(t, "10.0.0.1", resp.Subnet.Gateway)
	assert.Len(t, resp.Subnet.IPRanges, 1)
	assert.Equal(t, "10.0.0.10", resp.Subnet.IPRanges[0].StartAddress)

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.CleanResponse(endpoints.UpdateVDCNetwork())
}

func TestUpdateVdcNetworkRouted(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	edgeGwID := generator.MustGenerate("{uuid}")
	guestVLAN := true
	shared := true

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.SetResponse(endpoints.GetVDCNetwork(), &itypes.APIResponseVDCNetwork{
		ID:          networkID,
		Name:        "routed-net-1",
		Description: "old description",
		NetworkType: types.VDCNetworkTypeRouted,
		OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: "vdcg1"},
		Connection: &itypes.APIVDCNetworkConnection{
			RouterRef:           itypes.APIObjectReference{ID: "old-edge-gw", Name: "old-gw"},
			ConnectionTypeValue: "INTERNAL",
		},
		GuestVLANTaggingAllowed: &guestVLAN,
		Shared:                  &shared,
		Subnets: itypes.APIVDCNetworkSubnets{
			Values: []itypes.APIVDCNetworkSubnetValue{{
				Gateway:      "192.168.10.1",
				PrefixLength: 24,
				IPRanges: itypes.APIVDCNetworkIPRanges{Values: []itypes.APIVDCNetworkIPRangeValue{{
					StartAddress: "192.168.10.10",
					EndAddress:   "192.168.10.20",
				}}},
			}},
		},
	}, nil)

	ms.CleanResponse(endpoints.UpdateVDCNetwork())

	resp, err := client.UpdateVDCNetworkRouted(t.Context(), types.ParamsUpdateVDCNetworkRouted{
		ID:          networkID,
		Description: "new description",
		Subnet: &types.ParamsSubnet{
			Gateway:      "10.0.0.1",
			PrefixLength: 24,
			IPRanges: []types.ParamsVDCNetworkIPRange{{
				StartAddress: "10.0.0.10",
				EndAddress:   "10.0.0.20",
			}},
		},
		EdgeGatewayID:   edgeGwID,
		EdgeGatewayName: "new-gw",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, networkID, resp.ID)
	assert.Equal(t, "new description", resp.Description)
	assert.Equal(t, "10.0.0.1", resp.Subnet.Gateway)
	assert.Equal(t, edgeGwID, resp.EdgeGatewayID)
	assert.Equal(t, "new-gw", resp.EdgeGatewayName)
	assert.Len(t, resp.Subnet.IPRanges, 1)

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.CleanResponse(endpoints.UpdateVDCNetwork())
}
