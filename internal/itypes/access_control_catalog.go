/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

// APIResponseCatalogAccessControlGrants is the response of GET /catalogs/{catalogUrn}/accessControls/get/.
// It wraps the list of access control grants.
type APIResponseCatalogAccessControlGrants struct {
	Values []APIResponseAccessControlGrant `json:"values"`
}

// APIResponseAccessControlGrant represents a single access control grant for a catalog.
type APIResponseAccessControlGrant struct {
	AccessControlID  string `json:"accessControlId"`
	SubjectName      string `json:"subjectName"`
	SubjectType      string `json:"subjectType"`
	RoleURN          string `json:"roleUrn"`
	RoleName         string `json:"roleName"`
	OrganizationID   string `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
}

// APIRequestCatalogAccessControlGrant is the request payload shape used by PUT /catalogs/{catalogUrn}/accessControls/put/.
type APIRequestCatalogAccessControlGrant struct {
	SubjectName string `json:"subjectName"`
	SubjectType string `json:"subjectType"`
	RoleURN     string `json:"roleUrn"`
}

// APIRequestCatalogAccessControlGrants is the request body for setting catalog access control.
type APIRequestCatalogAccessControlGrants struct {
	Values []APIRequestCatalogAccessControlGrant `json:"values"`
}

// ToModel converts the internal access control grant list to the public model.
func (r *APIResponseCatalogAccessControlGrants) ToModel() *types.ModelListCatalogAccessControlGrant {
	model := &types.ModelListCatalogAccessControlGrant{
		Grants: make([]types.ModelCatalogAccessControlGrant, 0, len(r.Values)),
	}

	for _, grant := range r.Values {
		model.Grants = append(model.Grants, types.ModelCatalogAccessControlGrant{
			AccessControlID:  grant.AccessControlID,
			SubjectName:      grant.SubjectName,
			SubjectType:      grant.SubjectType,
			RoleURN:          grant.RoleURN,
			RoleName:         grant.RoleName,
			OrganizationID:   grant.OrganizationID,
			OrganizationName: grant.OrganizationName,
		})
	}

	return model
}
