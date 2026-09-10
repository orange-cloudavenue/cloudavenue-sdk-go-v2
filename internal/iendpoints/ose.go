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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

func init() {
	// ListAssociatedTenants returns tenants associated with the current user.
	cav.Endpoint{
		DocumentationURL: "https://s3console.cloudavenue.orange-business.com/api/v1/core/associated-tenants",
		Name:             "ListAssociatedTenants",
		Description:      "List tenants associated with the current user",
		Method:           cav.MethodGET,
		Backend:          cav.BackendOSE,
		PathTemplate:     "/api/v1/core/associated-tenants",
		QueryParams: []cav.QueryParam{
			{
				Name:        "accessible-only",
				Description: "Filter to only accessible tenants",
				Value:       "true",
			},
		},
		ResponseType: itypes.APIResponseListAssociatedTenants{},
	}.Register()

	// GetS3Credentials returns S3 access keys for a user in an organization.
	cav.Endpoint{
		DocumentationURL: "https://s3console.cloudavenue.orange-business.com/api/v1/core/tenants/{organizationID}/users/{userName}/credentials",
		Name:             "GetS3Credentials",
		Description:      "Get S3 credentials for a user in an organization",
		Method:           cav.MethodGET,
		Backend:          cav.BackendOSE,
		PathTemplate:     "/api/v1/core/tenants/{organizationID}/users/{userName}/credentials",
		PathParams: []cav.PathParam{
			{
				Name:        "organizationID",
				Description: "Organization identifier",
				Required:    true,
			},
			{
				Name:        "userName",
				Description: "User name",
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseS3Credentials{},
	}.Register()
}
