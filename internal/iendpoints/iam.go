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

//go:generate endpoint-generator -path iam.go -output iam

func init() {
	const pathAdminOrg = "/api/admin/org/{orgId}"

	// ListUsers
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Users.html",
		Name:             "ListUsers",
		Description:      "List users in organization",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/users",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
		},
		ResponseType: itypes.Users{},
	}.Register()

	// GetUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-User.html",
		Name:             "GetUser",
		Description:      "Get user by ID or name",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/user/{userId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamUserID,
				Description: descUserID,
				Required:    true,
			},
		},
		ResponseType: itypes.User{},
	}.Register()

	// CreateUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-User.html",
		Name:             "CreateUser",
		Description:      "Create a new user in organization",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/users",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.User{},
		ResponseType:    itypes.User{},
	}.Register()

	// UpdateUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-User.html",
		Name:             "UpdateUser",
		Description:      "Update an existing user",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/user/{userId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamUserID,
				Description: descUserID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.User{},
		ResponseType:    itypes.User{},
	}.Register()

	// DeleteUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-User.html",
		Name:             "DeleteUser",
		Description:      "Delete a user",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/user/{userId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamUserID,
				Description: descUserID,
				Required:    true,
			},
		},
		QueryParams: []cav.QueryParam{
			{
				Name:        "takeOwnership",
				Description: "Take ownership of user's resources",
				Required:    false,
			},
		},
		ResponseType: struct{}{},
	}.Register()

	// EnableUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-UserEnable.html",
		Name:             "EnableUser",
		Description:      "Enable a user",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/user/{userId}/action/enable",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamUserID,
				Description: descUserID,
				Required:    true,
			},
		},
		ResponseType: itypes.User{},
	}.Register()

	// DisableUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-UserDisable.html",
		Name:             "DisableUser",
		Description:      "Disable a user",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/user/{userId}/action/disable",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamUserID,
				Description: descUserID,
				Required:    true,
			},
		},
		ResponseType: itypes.User{},
	}.Register()

	// UnlockUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-UserUnlock.html",
		Name:             "UnlockUser",
		Description:      "Unlock a user",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/user/{userId}/action/unlock",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamUserID,
				Description: descUserID,
				Required:    true,
			},
		},
		ResponseType: itypes.User{},
	}.Register()

	// ChangePassword
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-UserChangePassword.html",
		Name:             "ChangePassword",
		Description:      "Change a user's password",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg + "/user/{userId}/action/changePassword",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamUserID,
				Description: descUserID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.NewPassword{},
		ResponseType:    struct{}{},
	}.Register()
}
