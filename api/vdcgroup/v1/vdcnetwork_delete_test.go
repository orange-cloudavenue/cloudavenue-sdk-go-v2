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

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.SetResponse(endpoints.GetVDCNetwork(), &itypes.APIResponseVDCNetwork{
		ID:          networkID,
		Name:        "isolated-net-1",
		NetworkType: types.VDCNetworkTypeIsolated,
	}, nil)

	ms.CleanResponse(endpoints.DeleteVDCNetwork())

	err := client.DeleteVDCNetworkIsolated(t.Context(), types.ParamsDeleteVDCNetworkIsolated{ID: networkID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.CleanResponse(endpoints.DeleteVDCNetwork())
}

func TestDeleteVdcNetworkRouted(t *testing.T) {
	networkID := generator.MustGenerate("{urn:network}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.SetResponse(endpoints.GetVDCNetwork(), &itypes.APIResponseVDCNetwork{
		ID:          networkID,
		Name:        "routed-net-1",
		NetworkType: types.VDCNetworkTypeRouted,
	}, nil)

	ms.CleanResponse(endpoints.DeleteVDCNetwork())

	err := client.DeleteVDCNetworkRouted(t.Context(), types.ParamsDeleteVDCNetworkRouted{ID: networkID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.CleanResponse(endpoints.DeleteVDCNetwork())
}
