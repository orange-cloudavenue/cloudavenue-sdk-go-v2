/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"reflect"
	"testing"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

func debugNewClient(t *testing.T) (*Client, *mock.Server) {
	t.Helper()

	results := reflect.ValueOf(newClient).Call([]reflect.Value{reflect.ValueOf(t)})
	client, _ := results[0].Interface().(*Client)

	if len(results) < 2 || results[1].IsNil() {
		return client, nil
	}

	ms, _ := results[1].Interface().(*mock.Server)
	return client, ms
}

func TestDebugListT0(t *testing.T) {
	client, ms := debugNewClient(t)
	if ms == nil {
		t.Skip("mock server not available")
	}

	ep := endpoints.ListT0()

	t.Logf("ep.Name=%q, ep.PathTemplate=%q\n", ep.Name, ep.PathTemplate)

	// Check what endpoints are registered
	for _, e := range cav.GetEndpointsUncategorized() {
		if e.Name == "ListT0" || e.Name == "GetEdgeGatewayServices" {
			t.Logf("Registered endpoint: Name=%q PathTemplate=%q BodyResponseType=%T\n", e.Name, e.PathTemplate, e.ResponseType)
		}
	}

	mockResponse := &itypes.APIResponseT0s{
		{
			Type: "tier-0-vrf",
			Name: "test-t0",
		},
	}
	statusCode := 404
	ms.SetResponse(ep, mockResponse, &statusCode)

	t.Logf("Set response for %s: data=%v, statusCode=%d\n", ep.Name, mockResponse, statusCode)

	t0s, err := client.ListT0(t.Context())
	t.Logf("Result: t0s=%v, err=%v\n", t0s, err)
}
