/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package endpoints_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

const (
	// cerberusPathPrefix is the live Cerberus customer-proxy path prefix.
	cerberusPathPrefix = "/infrapicustomerproxy/"

	// staleInfrapiPathPrefix is the retired prefix that must never reappear.
	staleInfrapiPathPrefix = "/api/customers/"
)

// Test_GeneratedWrappers_ResolveCerberusPaths checks the public generated
// wrapper surface. The wrappers are name-only indirections and embed no path,
// so the paths they hand out come straight from the registry: every one must
// already be the migrated Cerberus prefix.
func Test_GeneratedWrappers_ResolveCerberusPaths(t *testing.T) {
	wrappers := map[string]*cav.Endpoint{
		"ListDraasOnPremiseIP":       endpoints.ListDraasOnPremiseIP(),
		"AddDraasOnPremiseIP":        endpoints.AddDraasOnPremiseIP(),
		"RemoveDraasOnPremiseIP":     endpoints.RemoveDraasOnPremiseIP(),
		"CreateEdgeGateway":          endpoints.CreateEdgeGateway(),
		"DeleteEdgeGateway":          endpoints.DeleteEdgeGateway(),
		"UpdateEdgeGatewayBandwidth": endpoints.UpdateEdgeGatewayBandwidth(),
		"GetEdgeGatewayServices":     endpoints.GetEdgeGatewayServices(),
		"CreatePublicIP":             endpoints.CreatePublicIP(),
		"EnableCloudavenueServices":  endpoints.EnableCloudavenueServices(),
		"DisableCloudavenueServices": endpoints.DisableCloudavenueServices(),
		"GetOrganization":            endpoints.GetOrganization(),
		"UpdateOrganization":         endpoints.UpdateOrganization(),
		"CreateVDC":                  endpoints.CreateVDC(),
		"UpdateVDC":                  endpoints.UpdateVDC(),
		"DeleteVDC":                  endpoints.DeleteVDC(),
		"GetJobCerberus":             endpoints.GetJobCerberus(),
		"ListT0":                     endpoints.ListT0(),
		"GetT0":                      endpoints.GetT0(),
	}

	for name, ep := range wrappers {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, name, ep.Name)
			assert.Equal(t, cav.BackendInfrapi, ep.Backend)
			assert.NotContains(t, ep.PathTemplate, staleInfrapiPathPrefix,
				"endpoint %q uses the retired Infrapi prefix", name)
			assert.True(t, strings.HasPrefix(ep.PathTemplate, cerberusPathPrefix),
				"endpoint %q must start with %q, got %q", name, cerberusPathPrefix, ep.PathTemplate)
		})
	}
}
