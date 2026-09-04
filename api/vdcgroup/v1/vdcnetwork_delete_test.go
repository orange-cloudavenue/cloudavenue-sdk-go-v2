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

func TestDeleteVdcNetworkIsolated(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVdcNetwork())
	ms.SetResponse(endpoints.GetVdcNetwork(), &itypes.ApiResponseVdcNetwork{
		ID:          networkID,
		Name:        "isolated-net-1",
		NetworkType: types.VdcNetworkTypeIsolated,
	}, nil)

	ms.CleanResponse(endpoints.DeleteVdcNetwork())

	err := client.DeleteVdcNetworkIsolated(t.Context(), types.ParamsDeleteVdcNetworkIsolated{ID: networkID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetVdcNetwork())
	ms.CleanResponse(endpoints.DeleteVdcNetwork())
}

func TestDeleteVdcNetworkRouted(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVdcNetwork())
	ms.SetResponse(endpoints.GetVdcNetwork(), &itypes.ApiResponseVdcNetwork{
		ID:          networkID,
		Name:        "routed-net-1",
		NetworkType: types.VdcNetworkTypeRouted,
	}, nil)

	ms.CleanResponse(endpoints.DeleteVdcNetwork())

	err := client.DeleteVdcNetworkRouted(t.Context(), types.ParamsDeleteVdcNetworkRouted{ID: networkID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetVdcNetwork())
	ms.CleanResponse(endpoints.DeleteVdcNetwork())
}
