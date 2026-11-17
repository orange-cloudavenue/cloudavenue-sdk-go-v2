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
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertOrder fails when tokens do not appear in the given relative order.
func assertOrder(t *testing.T, out string, tokens ...string) {
	t.Helper()
	prev := -1
	for _, token := range tokens {
		idx := strings.Index(out, token)
		require.GreaterOrEqual(t, idx, 0, "missing %q in %s", token, out)
		assert.Greater(t, idx, prev, "token %q out of order in %s", token, out)
		prev = idx
	}
}

func TestInstantiateVAppTemplateMarshalSequence(t *testing.T) {
	out, err := xml.Marshal(APIRequestInstantiateVAppTemplate{
		OVF:         "http://schemas.dmtf.org/ovf/envelope/1",
		XMLNS:       "http://www.vmware.com/vcloud/v1.5",
		Name:        "order-probe",
		Description: "desc",
		InstantiationParams: &APIRequestVAppInstantiationParams{
			NetworkConfigSection: &APIRequestVAppNetworkConfigSection{
				Info: "NetInfo",
				NetworkConfig: []APIRequestVAppNetworkConfig{{
					NetworkName: "VM Network",
					Configuration: APIRequestVAppNetworkConfiguration{
						ParentNetwork: &APIRequestVAppReference{HREF: "https://example.com/api/network/1", Name: "VM Network", Type: "application/vnd.vmware.vcloud.network+xml"},
						FenceMode:     "bridged",
					},
				}},
			},
		},
		Source: APIRequestVAppReference{
			HREF: "https://example.com/api/vAppTemplate/vappTemplate-1",
			Name: "imported",
			Type: "application/vnd.vmware.vcloud.vAppTemplate+xml",
		},
	})
	require.NoError(t, err)
	assert.NotContains(t, string(out), "AllEULAsAccepted")

	// Live VCD 39.1 XSD rejects InstantiationParams after Source; the wire
	// sequence must be Description, InstantiationParams, Source.
	assertOrder(t, string(out),
		"<Description>desc</Description>",
		"<InstantiationParams>",
		"<NetworkConfigSection>",
		"<ovf:Info>NetInfo</ovf:Info>",
		"<ParentNetwork ",
		"<FenceMode>bridged</FenceMode>",
		"</Configuration>",
		"</NetworkConfigSection>",
		"</InstantiationParams>",
		"<Source ",
	)
	t.Logf("marshaled instantiate: %s", string(out))
}

func TestComposeVAppMarshalSequence(t *testing.T) {
	out, err := xml.Marshal(APIRequestComposeVApp{
		OVF:         "http://schemas.dmtf.org/ovf/envelope/1",
		XSI:         "http://www.w3.org/2001/XMLSchema-instance",
		XMLNS:       "http://www.vmware.com/vcloud/v1.5",
		Name:        "compose-probe",
		Description: "desc",
		InstantiationParams: &APIRequestVAppInstantiationParams{
			NetworkConnectionSection: &APIRequestVAppNetworkConnectionSection{Info: "NetConnInfo", PrimaryNetworkConnectionIndex: 0},
		},
		SourcedItem: &APIRequestVAppSourcedItem{
			Source:            APIRequestVAppReference{HREF: "https://example.com/api/vAppTemplate/vappTemplate-1"},
			VAppScopedLocalID: "vm-0",
		},
		AllEULAsAccepted: true,
	})
	require.NoError(t, err)

	// ComposeVAppParams XSD sequence: Description, InstantiationParams,
	// SourcedItem, AllEULAsAccepted (AllEULAsAccepted is an element, not an
	// attribute — verified live on CloudAvenue VCD 39.1).
	assertOrder(t, string(out),
		"<Description>desc</Description>",
		"<InstantiationParams>",
		"</InstantiationParams>",
		"<SourcedItem>",
		"</SourcedItem>",
		"<AllEULAsAccepted>",
	)
	t.Logf("marshaled compose: %s", string(out))
}
