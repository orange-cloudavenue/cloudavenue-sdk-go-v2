/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

const (
	// cerberusPathPrefix is the live Cerberus customer-proxy path prefix. The
	// pre-2026 Infrapi-era "/api/customers/" prefix is dead and answers 404.
	cerberusPathPrefix = "/infrapicustomerproxy/"

	// staleInfrapiPathPrefix is the retired prefix that must never reappear.
	staleInfrapiPathPrefix = "/api/customers/"
)

// Pinned upstream paths. These are test-local so they stay an independent
// anchor: the mock chi router is built from PathTemplate itself, so no
// behavioural test can tell a stale prefix from a live one.
const (
	expectedPathDraasIPList    = "/infrapicustomerproxy/v2.0/vcda/ips"
	expectedPathDraasIP        = "/infrapicustomerproxy/v2.0/vcda/ips/{ip}"
	expectedPathEdgeCreate     = "/infrapicustomerproxy/v2.0/{vdcType}/{vdcName}/edges"
	expectedPathEdgeList       = "/infrapicustomerproxy/v2.0/edges"
	expectedPathEdgeByID       = "/infrapicustomerproxy/v2.0/edges/{edgeId}"
	expectedPathNetwork        = "/infrapicustomerproxy/v2.0/network"
	expectedPathServices       = "/infrapicustomerproxy/v2.0/services"
	expectedPathServiceByID    = "/infrapicustomerproxy/v2.0/services/{serviceId}"
	expectedPathConfigurations = "/infrapicustomerproxy/v2.0/configurations"
	expectedPathVDCList        = "/infrapicustomerproxy/v2.0/vdcs"
	expectedPathVDCByName      = "/infrapicustomerproxy/v2.0/vdcs/{vdcName}"
	expectedPathT0List         = "/infrapicustomerproxy/v2.0/tier-0-vrfs"
	expectedPathT0ByName       = "/infrapicustomerproxy/v2.0/tier-0-vrfs/{tier0Name}"
	expectedPathJobByTaskID    = "/infrapicustomerproxy/v1.0/jobs/{taskId}"
)

// cerberusEndpoints pins the literal PathTemplate of every BackendInfrapi
// endpoint in the registry.
var cerberusEndpoints = []struct {
	name         string
	pathTemplate string
}{
	{"ListDraasOnPremiseIP", expectedPathDraasIPList},
	{"AddDraasOnPremiseIP", expectedPathDraasIP},
	{"RemoveDraasOnPremiseIP", expectedPathDraasIP},
	{"CreateEdgeGateway", expectedPathEdgeCreate},
	{"GetEdgeGateway", expectedPathEdgeByID},
	{"QueryEdgeGateway", expectedPathEdgeList},
	{"ListEdgeGateway", expectedPathEdgeList},
	{"DeleteEdgeGateway", expectedPathEdgeByID},
	{"UpdateEdgeGatewayBandwidth", expectedPathEdgeByID},
	{"GetEdgeGatewayServices", expectedPathNetwork},
	{"CreatePublicIP", expectedPathServices},
	{"EnableCloudavenueServices", expectedPathServices},
	{"DisableCloudavenueServices", expectedPathServiceByID},
	{"GetOrganization", expectedPathConfigurations},
	{"UpdateOrganization", expectedPathConfigurations},
	{"CreateVDC", expectedPathVDCList},
	{"ListVDC", expectedPathVDCList},
	{"GetVDC", expectedPathVDCByName},
	{"UpdateVDC", expectedPathVDCByName},
	{"DeleteVDC", expectedPathVDCByName},
	{"GetJobCerberus", expectedPathJobByTaskID},
	{"ListT0", expectedPathT0List},
	{"GetT0", expectedPathT0ByName},
}

// Test_Cerberus_PathTemplate_PinnedLiterals asserts each registered
// BackendInfrapi endpoint carries the exact upstream Cerberus path.
func Test_Cerberus_PathTemplate_PinnedLiterals(t *testing.T) {
	for _, tt := range cerberusEndpoints {
		t.Run(tt.name, func(t *testing.T) {
			ep, err := cav.GetEndpoint(tt.name)
			require.NoError(t, err)
			assert.Equal(t, cav.BackendInfrapi, ep.Backend, "endpoint must target the Cerberus backend")
			assert.Equal(t, tt.pathTemplate, ep.PathTemplate)
		})
	}
}

// Test_Cerberus_PathTemplate_NoStaleInfrapiPrefix is the invariant guard: no
// registered BackendInfrapi endpoint may still use the retired
// "/api/customers/" prefix. It fails loudly for any new or reverted endpoint,
// including ones absent from cerberusEndpoints.
func Test_Cerberus_PathTemplate_NoStaleInfrapiPrefix(t *testing.T) {
	var inspected int

	for _, ep := range cav.GetEndpointsUncategorized() {
		if ep.Backend != cav.BackendInfrapi {
			continue
		}

		inspected++
		assert.NotContains(t, ep.PathTemplate, staleInfrapiPathPrefix,
			"endpoint %q uses the retired Infrapi prefix %q", ep.Name, staleInfrapiPathPrefix)
		assert.True(t, strings.HasPrefix(ep.PathTemplate, cerberusPathPrefix),
			"endpoint %q must start with %q, got %q", ep.Name, cerberusPathPrefix, ep.PathTemplate)
	}

	assert.Equal(t, len(cerberusEndpoints), inspected,
		"every registered BackendInfrapi endpoint must be pinned in cerberusEndpoints")
}

// Test_Cerberus_PathTemplate_TableCoversRegistry keeps the pinned table and the
// live registry in sync, so neither can drift into agreement.
func Test_Cerberus_PathTemplate_TableCoversRegistry(t *testing.T) {
	registered := make([]string, 0, len(cerberusEndpoints))
	for _, ep := range cav.GetEndpointsUncategorized() {
		if ep.Backend == cav.BackendInfrapi {
			registered = append(registered, ep.Name)
		}
	}

	pinned := make([]string, 0, len(cerberusEndpoints))
	for _, tt := range cerberusEndpoints {
		pinned = append(pinned, tt.name)
	}

	slices.Sort(registered)
	slices.Sort(pinned)
	assert.Equal(t, pinned, registered)
}
