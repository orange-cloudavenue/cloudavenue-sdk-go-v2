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

//go:generate endpoint-generator -path admin.go -output admin

func init() {
	const pathAdminOrgs = "/api/admin/orgs"
	const pathAdminOrg = "/api/admin/org/{orgId}"

	// ListAdminOrgs
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminOrgs.html",
		Name:             "ListAdminOrgs",
		Description:      "List organizations (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrgs,
		ResponseType:     itypes.AdminOrgs{},
	}.Register()

	// GetAdminOrg
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminOrg.html",
		Name:             "GetAdminOrg",
		Description:      "Get an organization by ID (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
		},
		ResponseType: itypes.AdminOrg{},
	}.Register()
}
