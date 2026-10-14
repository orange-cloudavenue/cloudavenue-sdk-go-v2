/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

func TestCreateOAuthClientEndpoint(t *testing.T) {
	ep := endpoints.CreateOAuthClient()
	require.NotNil(t, ep)
	assert.Equal(t, cav.BackendIHM, ep.Backend)
	assert.Equal(t, http.MethodPost, ep.Method.String())
	assert.Equal(t, "/oauth/tenant/{tenantId}/register", ep.PathTemplate)
	require.Len(t, ep.PathParams, 1)
	assert.Equal(t, "tenantId", ep.PathParams[0].Name)
}
