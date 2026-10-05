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
		BodyRequestType: itypes.UserRequest{},
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
		BodyRequestType: itypes.UserRequest{},
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

	const pathTokens = "/cloudapi/1.0.0/tokens"

	// ListTokens
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Tokens.html",
		Name:             "ListTokens",
		Description:      "List tokens in organization",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokens + "/get/",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseListTokens{},
	}.Register()

	// GetToken
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Token.html",
		Name:             "GetToken",
		Description:      "Get token by ID",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokens + "/id/get/",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamTokenID,
				Description: descTokenID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseToken{},
	}.Register()

	// CreateToken
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-Token.html",
		Name:             "CreateToken",
		Description:      "Create a new token in organization",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokens + "/post/",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestToken{},
		ResponseType:    itypes.APIResponseToken{},
	}.Register()

	// UpdateToken
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-Token.html",
		Name:             "UpdateToken",
		Description:      "Update an existing token",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokens + "/id/put/",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamTokenID,
				Description: descTokenID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestToken{},
		ResponseType:    itypes.APIResponseToken{},
	}.Register()

	// DeleteToken
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-Token.html",
		Name:             "DeleteToken",
		Description:      "Delete a token",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokens + "/id/delete/",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamTokenID,
				Description: descTokenID,
				Required:    true,
			},
		},
		ResponseType: struct{}{},
	}.Register()

	const pathLDAP = "/cloudapi/1.0.0/ldap"

	const pathGlobalRoles = "/cloudapi/1.0.0/globalRoles"

	// TestLDAP
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/ldap/test/post/",
		Name:             "TestLDAP",
		Description:      "Test LDAP connection",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathLDAP + "/test",
		QueryParams: []cav.QueryParam{
			{
				Name:        "username",
				Description: "Username to use when testing LDAP search",
				Required:    false,
			},
		},
		BodyRequestType: itypes.APIRequestLDAPTest{},
		ResponseType:    itypes.APIResponseLDAPTestResult{},
	}.Register()

	// SyncLDAP
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-LDAPSync.html",
		Name:             "SyncLDAP",
		Description:      "Synchronize LDAP directory",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathLDAP + "/sync",
		BodyRequestType:  nil, // No request body for this endpoint.
		ResponseType:     struct{}{},
	}.Register()

	// SearchLDAPUsers
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/ldap/search/user/get/",
		Name:             "SearchLDAPUsers",
		Description:      "Search LDAP users",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathLDAP + "/search/user",
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamQ,
				Description: "String to search for via LDAP",
				Required:    false,
			},
		},
		ResponseType: []itypes.APIResponseLDAPUser{},
	}.Register()

	// SearchLDAPGroups
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/latest/cloudapi/1.0.0/ldap/search/group/get/",
		Name:             "SearchLDAPGroups",
		Description:      "Search LDAP groups",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathLDAP + "/search/group",
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamQ,
				Description: "String to search for via LDAP",
				Required:    false,
			},
		},
		ResponseType: []itypes.APIResponseLDAPGroup{},
	}.Register()

	// ListGlobalRoles
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-GlobalRoles.html",
		Name:             "ListGlobalRoles",
		Description:      "List global roles",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles,
		ResponseType:     itypes.APIResponseListGlobalRoles{},
	}.Register()

	// GetGlobalRole
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-GlobalRole.html",
		Name:             "GetGlobalRole",
		Description:      "Get global role by ID",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseGlobalRole{},
	}.Register()

	// CreateGlobalRole
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRole.html",
		Name:             "CreateGlobalRole",
		Description:      "Create a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles,
		BodyRequestType:  itypes.APIRequestGlobalRole{},
		ResponseType:     itypes.APIResponseGlobalRole{},
	}.Register()

	// UpdateGlobalRole
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-GlobalRole.html",
		Name:             "UpdateGlobalRole",
		Description:      "Update a global role",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRole{},
		ResponseType:    itypes.APIResponseGlobalRole{},
	}.Register()

	// DeleteGlobalRole
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-GlobalRole.html",
		Name:             "DeleteGlobalRole",
		Description:      "Delete a global role",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: struct{}{},
	}.Register()

	// ListGlobalRoleRights
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-GlobalRoleRights.html",
		Name:             "ListGlobalRoleRights",
		Description:      "List rights of a global role",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/rights",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseGlobalRoleRights{},
	}.Register()

	// AddGlobalRoleRights
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleRights.html",
		Name:             "AddGlobalRoleRights",
		Description:      "Add rights to a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/rights",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleRights{},
		ResponseType:    itypes.APIResponseGlobalRoleRights{},
	}.Register()

	// ReplaceGlobalRoleRights
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-GlobalRoleRights.html",
		Name:             "ReplaceGlobalRoleRights",
		Description:      "Replace rights of a global role",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/rights",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleRights{},
		ResponseType:    itypes.APIResponseGlobalRoleRights{},
	}.Register()

	// ListGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-GlobalRoleTenants.html",
		Name:             "ListGlobalRoleTenants",
		Description:      "List tenants of a global role",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/tenants",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseGlobalRoleTenants{},
	}.Register()

	// SetGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-GlobalRoleTenants.html",
		Name:             "SetGlobalRoleTenants",
		Description:      "Set tenants of a global role",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/tenants",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleTenants{},
		ResponseType:    itypes.APIResponseGlobalRoleTenants{},
	}.Register()

	// PublishGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleTenantsPublish.html",
		Name:             "PublishGlobalRoleTenants",
		Description:      "Publish tenants of a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/tenants/publish",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleTenants{},
		ResponseType:    struct{}{},
	}.Register()

	// UnpublishGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleTenantsUnpublish.html",
		Name:             "UnpublishGlobalRoleTenants",
		Description:      "Unpublish tenants of a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/tenants/unpublish",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleTenants{},
		ResponseType:    struct{}{},
	}.Register()

	// PublishAllGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleTenantsPublishAll.html",
		Name:             "PublishAllGlobalRoleTenants",
		Description:      "Publish all tenants of a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/tenants/publishAll",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: struct{}{},
	}.Register()

	// UnpublishAllGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleTenantsUnpublishAll.html",
		Name:             "UnpublishAllGlobalRoleTenants",
		Description:      "Unpublish all tenants of a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles + "/{id}/tenants/unpublishAll",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamGlobalRoleID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: struct{}{},
	}.Register()
}
