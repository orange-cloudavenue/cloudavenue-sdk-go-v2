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
		Name                string `json:"name" fake:"{resource_name:organization}"`
		FullName            string `json:"fullName" fake:"{company}"`
		Description         string `json:"description" fake:"{sentence:10,20}"`
		IsEnabled           bool   `json:"isEnabled" fake:"true"`
		IsSuspended         bool   `json:"isSuspended" fake:"false"`
		CustomerMail        string `json:"customerMail" fake:"{email}"`
		InternetBillingMode string `json:"internetBillingMode" fake:"PAYG"`
	}

	// Organization represents an organization entity from Vmware Cloud Director.
	ApiResponseGetOrgs struct {
		Organizations []ApiResponseGetOrgDetails `json:"values" fakesize:"1"`
	}

	ApiResponseGetOrgDetails struct {
		ID             string `json:"id" fake:"{urn:org}"`
		Name           string `json:"name" fake:"{resource_name:organization}"`
		DisplayName    string `json:"displayName" fake:"{company}"`
		Description    string `json:"description" fake:"{sentence:10,20}"`
		IsEnabled      bool   `json:"isEnabled" fake:"true"`
		OrgVdcCount    int    `json:"orgVdcCount" fake:"{number:1,5}"`
		CatalogCount   int    `json:"catalogCount" fake:"{number:1,5}"`
		VappCount      int    `json:"vappCount" fake:"{number:1,5}"`
		RunningVMCount int    `json:"runningVMCount" fake:"{number:1,5}"`
		UserCount      int    `json:"userCount" fake:"{number:1,5}"`
		DiskCount      int    `json:"diskCount" fake:"{number:1,5}"`
		CanPublish     bool   `json:"canPublish" fake:"false"`
	}

	ApiRequestUpdateOrg struct {
		FullName            string `json:"fullName" validate:"omitempty"`
		Description         string `json:"description" validate:"omitempty"`
		CustomerMail        string `json:"customerMail" validate:"omitempty,email"`
		InternetBillingMode string `json:"internetBillingMode" validate:"omitempty,oneof=PAYG TRAFFIC_VOLUME"`
	}
)

// From infraAPI
func (r *ApiResponseGetOrg) ToModel() *types.ModelGetOrganization {
	return &types.ModelGetOrganization{
		Name:        r.Name,
		DisplayName: r.FullName, // FullName is mapped to DisplayName
		Description: r.Description,
		IsEnabled:   r.IsEnabled,
		// IsSuspended:         r.IsSuspended,
		CustomerMail:        r.CustomerMail,
		InternetBillingMode: r.InternetBillingMode,
	}
}

// From Vmware Cloud Director
func (r *ApiResponseGetOrgs) ToModel() *types.ModelGetOrganization {
	if len(r.Organizations) == 0 {
		return nil
	}

	return &types.ModelGetOrganization{
		ID:          r.Organizations[0].ID,
		Name:        r.Organizations[0].Name,
		DisplayName: r.Organizations[0].DisplayName,
		Description: r.Organizations[0].Description,
		IsEnabled:   r.Organizations[0].IsEnabled,
		Resources: types.ModelGetOrganizationResources{
			Vdc:       r.Organizations[0].OrgVdcCount,
			Catalog:   r.Organizations[0].CatalogCount,
			Vapp:      r.Organizations[0].VappCount,
			RunningVM: r.Organizations[0].RunningVMCount,
			User:      r.Organizations[0].UserCount,
			Disk:      r.Organizations[0].DiskCount,
		},
	}
}

func (r *ApiRequestUpdateOrg) ToModel() *types.ModelUpdateOrganization {
	return &types.ModelUpdateOrganization{
		FullName:            r.FullName,
		Description:         r.Description,
		CustomerMail:        r.CustomerMail,
		InternetBillingMode: r.InternetBillingMode,
	}
}
