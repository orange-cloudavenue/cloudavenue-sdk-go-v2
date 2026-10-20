/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiResponseVdcNetworkToModelPreservesAllSubnets(t *testing.T) {
	network := APIResponseVDCNetwork{
		Subnets: APIVDCNetworkSubnets{
			Values: []APIVDCNetworkSubnetValue{
				{
					Gateway:      "10.0.0.1",
					PrefixLength: 24,
					IPRanges: APIVDCNetworkIPRanges{Values: []APIVDCNetworkIPRangeValue{{
						StartAddress: "10.0.0.10",
						EndAddress:   "10.0.0.20",
					}}},
				},
				{
					Gateway:      "10.0.1.1",
					PrefixLength: 24,
					IPRanges: APIVDCNetworkIPRanges{Values: []APIVDCNetworkIPRangeValue{{
						StartAddress: "10.0.1.10",
						EndAddress:   "10.0.1.20",
					}}},
				},
			},
		},
	}

	model := network.ToModel()

	assert.Len(t, model.Subnets, 2)
	assert.Equal(t, "10.0.0.1", model.Subnet.Gateway)
	assert.Equal(t, "10.0.0.1", model.Subnets[0].Gateway)
	assert.Equal(t, "10.0.1.1", model.Subnets[1].Gateway)
	assert.Equal(t, "10.0.1.10", model.Subnets[1].IPRanges[0].StartAddress)
}
