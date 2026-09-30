/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package organization

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestGetCatalogAccessControl(t *testing.T) {
	catalogURN := "urn:vcloud:catalog:12345678-1234-4b8d-89ab-123456789012"
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetCatalogAccessControl())
	ms.SetResponse(endpoints.GetCatalogAccessControl(), &itypes.APIResponseCatalogAccessControlGrants{
		Values: []itypes.APIResponseAccessControlGrant{
			{
				AccessControlID:  "ac-1",
				SubjectName:      "org-1",
				SubjectType:      "org",
				RoleURN:          "urn:vcloud:role:12345678-1234-4b8d-89ab-123456789012",
				RoleName:         "Catalog Author",
				OrganizationID:   "urn:vcloud:org:12345678-1234-4b8d-89ab-123456789012",
				OrganizationName: "Org 1",
			},
			{
				AccessControlID:  "ac-2",
				SubjectName:      "user-1",
				SubjectType:      "user",
				RoleURN:          "urn:vcloud:role:12345678-1234-4b8d-89ab-123456789012",
				RoleName:         "Catalog Author",
				OrganizationID:   "urn:vcloud:org:12345678-1234-4b8d-89ab-123456789012",
				OrganizationName: "Org 1",
			},
		},
	}, nil)

	resp, err := client.GetCatalogAccessControl(t.Context(), types.ParamsGetCatalogAccessControl{CatalogURN: catalogURN})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Grants, 2)
	assert.Equal(t, "ac-1", resp.Grants[0].AccessControlID)
	assert.Equal(t, "org-1", resp.Grants[0].SubjectName)
	assert.Equal(t, "org", resp.Grants[0].SubjectType)
	assert.Equal(t, "urn:vcloud:role:12345678-1234-4b8d-89ab-123456789012", resp.Grants[0].RoleURN)
	assert.Equal(t, "Catalog Author", resp.Grants[0].RoleName)
	assert.Equal(t, "urn:vcloud:org:12345678-1234-4b8d-89ab-123456789012", resp.Grants[0].OrganizationID)
	assert.Equal(t, "Org 1", resp.Grants[0].OrganizationName)
	assert.Equal(t, "user-1", resp.Grants[1].SubjectName)
	assert.Equal(t, "user", resp.Grants[1].SubjectType)

	ms.CleanResponse(endpoints.GetCatalogAccessControl())
}

func TestGetCatalogAccessControl_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetCatalogAccessControl())

	_, err := client.GetCatalogAccessControl(t.Context(), types.ParamsGetCatalogAccessControl{})
	assert.Error(t, err)
}

func TestSetCatalogAccessControl(t *testing.T) {
	catalogURN := "urn:vcloud:catalog:12345678-1234-4b8d-89ab-123456789012"
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.SetCatalogAccessControl())
	ms.SetResponse(endpoints.SetCatalogAccessControl(), &itypes.APIResponseCatalogAccessControlGrants{
		Values: []itypes.APIResponseAccessControlGrant{
			{
				AccessControlID: "ac-1",
				SubjectName:     "org-1",
				SubjectType:     "org",
				RoleURN:         "urn:vcloud:role:12345678-1234-4b8d-89ab-123456789012",
				RoleName:        "Catalog Author",
			},
		},
	}, nil)

	resp, err := client.SetCatalogAccessControl(t.Context(), types.ParamsSetCatalogAccessControl{
		CatalogURN: catalogURN,
		Grants: []types.ModelCatalogAccessControlGrant{
			{
				SubjectName: "org-1",
				SubjectType: "org",
				RoleURN:     "urn:vcloud:role:12345678-1234-4b8d-89ab-123456789012",
			},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Grants, 1)
	assert.Equal(t, "ac-1", resp.Grants[0].AccessControlID)
	assert.Equal(t, "org-1", resp.Grants[0].SubjectName)
	assert.Equal(t, "org", resp.Grants[0].SubjectType)

	ms.CleanResponse(endpoints.SetCatalogAccessControl())
}

func TestSetCatalogAccessControl_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.SetCatalogAccessControl())

	_, err := client.SetCatalogAccessControl(t.Context(), types.ParamsSetCatalogAccessControl{})
	assert.Error(t, err)
}
