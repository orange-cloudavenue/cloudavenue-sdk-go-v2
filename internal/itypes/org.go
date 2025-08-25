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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

type (
	// ApiResponseGetOrg from infraAPI
	ApiResponseGetOrg struct {
		Name                string `json:"name"`
		FullName            string `json:"fullName"`
		Description         string `json:"description"`
		IsEnabled           bool   `json:"isEnabled"`
		IsSuspended         bool   `json:"isSuspended"`
		CustomerMail        string `json:"customerMail"`
		InternetBillingMode string `json:"internetBillingMode"`
	}

	// Organization represents an organization entity.
	ApiResponseListOrgs struct {
		Organizations []ApiResponseListOrgQuerrier `json:"values" fakesize:"1"`
	}

	ApiResponseListOrgQuerrier struct {
		ID             string `json:"id" fake:"{uuid}"`           // UUID of the organization
		Name           string `json:"name" fake:"{company_name}"` // Name
		DisplayName    string `json:"displayName" fake:"{company_name}"`
		Description    string `json:"description" fake:"{sentence:10,20}"`
		IsEnabled      bool   `json:"isEnabled" fake:"true"`
		OrgVdcCount    int    `json:"orgVdcCount" fake:"{number:1,5}"`
		CatalogCount   int    `json:"catalogCount" fake:"{number:1,5}"`
		VappCount      int    `json:"vappCount" fake:"{number:1,5}"`
		RunningVMCount int    `json:"runningVMCount" fake:"{number:1,5}"`
		UserCount      int    `json:"userCount" fake:"{number:1,5}"`
		DiskCount      int    `json:"diskCount" fake:"{number:1,5}"`
		// ManagedBy      ManagedBy `json:"managedBy"`
		// CanManageOrgs  bool      `json:"canManageOrgs" fake:"false"`
		CanPublish bool `json:"canPublish" fake:"false"`
	}

	ApiResponseUpdateOrg struct {
		Name                string `json:"name"`
		FullName            string `json:"fullName"`
		Description         string `json:"description"`
		IsEnabled           bool   `json:"isEnabled"`
		IsSuspended         bool   `json:"isSuspended"`
		CustomerMail        string `json:"customerMail"`
		InternetBillingMode string `json:"internetBillingMode"`
	}
)

func (r *ApiResponseListOrgs) ToModel() *types.ModelListOrganization {
	return &types.ModelListOrganization{
		ID:   r.Organizations[0].ID,
		Name: r.Organizations[0].Name,
	}
}
