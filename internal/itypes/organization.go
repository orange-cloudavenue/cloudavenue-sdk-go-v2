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
	// APIResponseGetOrg is infrapi organization response.
	APIResponseGetOrg struct {
		Name                string `json:"name" fake:"{resource_name:organization}"`
		FullName            string `json:"fullName" fake:"{company}"`
		Description         string `json:"description" fake:"{sentence:10,20}"`
		IsEnabled           bool   `json:"isEnabled" fake:"true"`
		IsSuspended         bool   `json:"isSuspended" fake:"false"`
		CustomerMail        string `json:"customerMail" fake:"{email}"`
		InternetBillingMode string `json:"internetBillingMode" fake:"PAYG"`
	}

	// APIResponseGetOrgs is VMware Cloud Director organization list response.
	APIResponseGetOrgs struct {
		Organizations []APIResponseGetOrgDetails `json:"values" fakesize:"1"`
	}

	APIResponseGetOrgDetails struct {
		ID             string `json:"id" fake:"{urn:org}"`
		Name           string `json:"name" fake:"{resource_name:organization}"`
		DisplayName    string `json:"displayName" fake:"{company}"`
		Description    string `json:"description" fake:"{sentence:10,20}"`
		IsEnabled      bool   `json:"isEnabled" fake:"true"`
		OrgVDCCount    int    `json:"orgVdcCount" fake:"{number:1,5}"`
		CatalogCount   int    `json:"catalogCount" fake:"{number:1,5}"`
		VAppCount      int    `json:"vappCount" fake:"{number:1,5}"`
		RunningVMCount int    `json:"runningVMCount" fake:"{number:1,5}"`
		UserCount      int    `json:"userCount" fake:"{number:1,5}"`
		DiskCount      int    `json:"diskCount" fake:"{number:1,5}"`
		CanPublish     bool   `json:"canPublish" fake:"false"`
	}

	APIRequestUpdateOrg struct {
		FullName            string `json:"fullName" validate:"omitempty"`
		Description         string `json:"description" validate:"omitempty"`
		CustomerMail        string `json:"customerMail" validate:"omitempty,email"`
		InternetBillingMode string `json:"internetBillingMode" validate:"omitempty,oneof=PAYG TRAFFIC_VOLUME"`
	}
)

// ToModel converts infrapi organization response.
func (r *APIResponseGetOrg) ToModel() *types.ModelGetOrganization {
	return &types.ModelGetOrganization{
		Name:                r.Name,
		FullName:            r.FullName,
		Description:         r.Description,
		Email:               r.CustomerMail,
		InternetBillingMode: r.InternetBillingMode,
	}
}

// ToModel converts VMware Cloud Director organization response.
func (r *APIResponseGetOrgs) ToModel() *types.ModelGetOrganization {
	if len(r.Organizations) == 0 {
		return nil
	}

	return &types.ModelGetOrganization{
		ID:          r.Organizations[0].ID,
		Name:        r.Organizations[0].Name,
		FullName:    r.Organizations[0].DisplayName,
		Description: r.Organizations[0].Description,
		Enabled:     r.Organizations[0].IsEnabled,
		Resources: types.ModelGetOrganizationResources{
			VDC:       r.Organizations[0].OrgVDCCount,
			Catalog:   r.Organizations[0].CatalogCount,
			VApp:      r.Organizations[0].VAppCount,
			VMRunning: r.Organizations[0].RunningVMCount,
			User:      r.Organizations[0].UserCount,
			Disk:      r.Organizations[0].DiskCount,
		},
	}
}
