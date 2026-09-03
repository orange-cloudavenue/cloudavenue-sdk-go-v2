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
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path appportprofile.go -output appportprofile

// AppPortProfile is a distinct NSX-T resource from FirewallGroup: it is not
// owned by a single OwnerRef, but instead carries both an OrgRef (Org ID) and
// a ContextEntityId (VDC, VdcGroup, or NSX-T Manager URN depending on Scope).
// Unlike FirewallGroup, List/Create use the bare resource path (no
// "/summaries" suffix). Create/Update/Delete are synchronous (no async job).
func init() {
	// ListAppPortProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/get/",
		Name:             "ListAppPortProfile",
		Description:      "List Application Port Profiles",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/applicationPortProfiles",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "Filter to apply to the list of Application Port Profiles. Format: key==value;key==value. Allowed keys: scope, _context, name, id.",
			},
			{
				Name:        "pageSize",
				Description: "The number of items per page.",
				Value:       "100",
			},
		},
		ResponseType: itypes.ApiResponseListAppPortProfile{},
	}.Register()

	// GetAppPortProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/appPortProfileId/get/",
		Name:             "GetAppPortProfile",
		Description:      "Get an Application Port Profile",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/applicationPortProfiles/{appPortProfileId}",
		PathParams: []cav.PathParam{
			{
				Name:        "appPortProfileId",
				Description: "ID of the Application Port Profile to get",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=applicationPortProfile")
				},
			},
		},
		ResponseType: itypes.ApiResponseAppPortProfile{},
	}.Register()

	// CreateAppPortProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/post/",
		Name:             "CreateAppPortProfile",
		Description:      "Create an Application Port Profile",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/applicationPortProfiles",
		BodyRequestType:  itypes.ApiRequestAppPortProfile{},
		ResponseType:     itypes.ApiResponseAppPortProfile{},
	}.Register()

	// UpdateAppPortProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/appPortProfileId/put/",
		Name:             "UpdateAppPortProfile",
		Description:      "Update an Application Port Profile",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/applicationPortProfiles/{appPortProfileId}",
		PathParams: []cav.PathParam{
			{
				Name:        "appPortProfileId",
				Description: "ID of the Application Port Profile to update",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=applicationPortProfile")
				},
			},
		},
		BodyRequestType: itypes.ApiRequestAppPortProfile{},
		ResponseType:    itypes.ApiResponseAppPortProfile{},
	}.Register()

	// DeleteAppPortProfile
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/appPortProfileId/delete/",
		Name:             "DeleteAppPortProfile",
		Description:      "Delete an Application Port Profile",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/applicationPortProfiles/{appPortProfileId}",
		PathParams: []cav.PathParam{
			{
				Name:        "appPortProfileId",
				Description: "ID of the Application Port Profile to delete",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=applicationPortProfile")
				},
			},
		},
	}.Register()
}
