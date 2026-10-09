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
	// IAM user operations use VMware CloudAPI. They are distinct from both the
	// Infrapi customer API and the legacy AdminOrg XML API.
	// ListUsers
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/users/get/",
		Name:             "ListUsers",
		Description:      "List users in organization",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUsers,
		QueryParams: []cav.QueryParam{
			{Name: "page", Description: "Page to fetch", Required: true, Value: "1"},
			{Name: "pageSize", Description: "Results per page to fetch", Required: true, Value: "128"},
		},
		ResponseType: itypes.APIResponseListUsers{},
	}.Register()

	// GetUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/users/userUrn/get/",
		Name:             "GetUser",
		Description:      "Get user by ID or name",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUser,
		PathParams:       []cav.PathParam{{Name: pathParamUserUrn, Description: descUserURN, Required: true}},
		ResponseType:     itypes.APIUser{},
	}.Register()

	// CreateUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/users/post/",
		Name:             "CreateUser",
		Description:      "Create a new user in organization",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUsers,
		PathParams:       []cav.PathParam{},
		BodyRequestType:  itypes.APIUser{},
		ResponseType:     itypes.APIUser{},
	}.Register()

	// UpdateUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/users/userUrn/put/",
		Name:             "UpdateUser",
		Description:      "Update an existing user",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUser,
		PathParams:       []cav.PathParam{{Name: pathParamUserUrn, Description: descUserURN, Required: true}},
		BodyRequestType:  itypes.APIUser{},
		ResponseType:     itypes.APIUser{},
	}.Register()

	// TakeOwnership transfers entities owned by a user to the caller.
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/users/userUrn/takeOwnership/post/",
		Name:             "TakeOwnership",
		Description:      "Take ownership of a user",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUserTakeOwnership,
		PathParams:       []cav.PathParam{{Name: pathParamUserUrn, Description: descUserURN, Required: true}},
		ResponseType:     struct{}{},
	}.Register()

	// DeleteUser
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/users/userUrn/delete/",
		Name:             "DeleteUser",
		Description:      "Delete a user",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUser,
		PathParams:       []cav.PathParam{{Name: pathParamUserUrn, Description: descUserURN, Required: true}},
		ResponseType:     struct{}{},
	}.Register()

	// EnableUser
	cav.Endpoint{
		DocumentationURL: docURLVMware,
		Name:             "EnableUser",
		Description:      "Enable a user",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUser,
		PathParams:       []cav.PathParam{{Name: pathParamUserUrn, Description: descUserURN, Required: true}},
		BodyRequestType:  itypes.APIUser{},
		ResponseType:     itypes.APIUser{},
	}.Register()

	// DisableUser
	cav.Endpoint{
		DocumentationURL: docURLVMware,
		Name:             "DisableUser",
		Description:      "Disable a user",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUser,
		PathParams:       []cav.PathParam{{Name: pathParamUserUrn, Description: descUserURN, Required: true}},
		BodyRequestType:  itypes.APIUser{},
		ResponseType:     itypes.APIUser{},
	}.Register()

	// UnlockUser
	cav.Endpoint{
		DocumentationURL: docURLVMware,
		Name:             "UnlockUser",
		Description:      "Unlock a user",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUser,
		PathParams:       []cav.PathParam{{Name: pathParamUserUrn, Description: descUserURN, Required: true}},
		BodyRequestType:  itypes.APIUser{},
		ResponseType:     itypes.APIUser{},
	}.Register()

	// ChangePassword
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/users/userUrn/changePassword/post/",
		Name:             "ChangePassword",
		Description:      "Change a user's password",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIUserChangePassword,
		PathParams:       []cav.PathParam{{Name: pathParamUserUrn, Description: descUserURN, Required: true}},
		BodyRequestType:  itypes.APIRequestPasswordChange{},
		ResponseType:     struct{}{},
	}.Register()

	// ListTokens
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/tokens/get/",
		Name:             "ListTokens",
		Description:      "List tokens in organization",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokenGet,
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
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/sessions/id/token/get/",
		Name:             "GetToken",
		Description:      "Get token by ID",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokenGetByID,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamID,
				Description: descTokenID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseToken{},
	}.Register()

	// CreateToken
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/tokens/post/",
		Name:             "CreateToken",
		Description:      "Create a new token in organization",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokenCreate,
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
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/tokens/id/put/",
		Name:             "UpdateToken",
		Description:      "Update an existing token",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokenUpdateByID,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamID,
				Description: descTokenID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestToken{},
		ResponseType:    itypes.APIResponseToken{},
	}.Register()

	// DeleteToken
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/tokens/id/delete/",
		Name:             "DeleteToken",
		Description:      "Delete a token",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTokenDeleteByID,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
			{
				Name:        pathParamID,
				Description: descTokenID,
				Required:    true,
			},
		},
		ResponseType: struct{}{},
	}.Register()

	// TestLDAP
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/ldap/test/post/",
		Name:             "TestLDAP",
		Description:      "Test LDAP connection",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathLDAPTest,
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
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/ldap/sync/post/",
		Name:             "SyncLDAP",
		Description:      "Synchronize LDAP directory",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathLDAPSync,
		BodyRequestType:  nil, // No request body for this endpoint.
		ResponseType:     struct{}{},
	}.Register()

	// SearchLDAPUsers
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/ldap/search/user/get/",
		Name:             "SearchLDAPUsers",
		Description:      "Search LDAP users",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathLDAPSearchUsers,
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
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/ldap/search/group/get/",
		Name:             "SearchLDAPGroups",
		Description:      "Search LDAP groups",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathLDAPSearchGroups,
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
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/get/",
		Name:             "ListGlobalRoles",
		Description:      "List global roles",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoles,
		ResponseType:     itypes.APIResponseListGlobalRoles{},
	}.Register()

	// GetGlobalRole
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/get/",
		Name:             "GetGlobalRole",
		Description:      "Get global role by ID",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleByID,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseGlobalRole{},
	}.Register()

	// CreateGlobalRole
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/post/",
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
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/put/",
		Name:             "UpdateGlobalRole",
		Description:      "Update a global role",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleByID,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRole{},
		ResponseType:    itypes.APIResponseGlobalRole{},
	}.Register()

	// DeleteGlobalRole
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/delete/",
		Name:             "DeleteGlobalRole",
		Description:      "Delete a global role",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleByID,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: struct{}{},
	}.Register()

	// ListGlobalRoleRights
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/rights/get/",
		Name:             "ListGlobalRoleRights",
		Description:      "List rights of a global role",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleRights,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseGlobalRoleRights{},
	}.Register()

	// AddGlobalRoleRights
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/rights/post/",
		Name:             "AddGlobalRoleRights",
		Description:      "Add rights to a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleRights,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleRights{},
		ResponseType:    itypes.APIResponseGlobalRoleRights{},
	}.Register()

	// ReplaceGlobalRoleRights
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/rights/put/",
		Name:             "ReplaceGlobalRoleRights",
		Description:      "Replace rights of a global role",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleRights,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleRights{},
		ResponseType:    itypes.APIResponseGlobalRoleRights{},
	}.Register()

	// ListGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/tenants/get/",
		Name:             "ListGlobalRoleTenants",
		Description:      "List tenants of a global role",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleTenants,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: itypes.APIResponseGlobalRoleTenants{},
	}.Register()

	// SetGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/tenants/put/",
		Name:             "SetGlobalRoleTenants",
		Description:      "Set tenants of a global role",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleTenants,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleTenants{},
		ResponseType:    itypes.APIResponseGlobalRoleTenants{},
	}.Register()

	// PublishGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/tenants/publish/post/",
		Name:             "PublishGlobalRoleTenants",
		Description:      "Publish tenants of a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleTenantsPublish,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleTenants{},
		ResponseType:    struct{}{},
	}.Register()

	// UnpublishGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/tenants/unpublish/post/",
		Name:             "UnpublishGlobalRoleTenants",
		Description:      "Unpublish tenants of a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleTenantsUnpublish,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestGlobalRoleTenants{},
		ResponseType:    struct{}{},
	}.Register()

	// PublishAllGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/tenants/publishAll/post/",
		Name:             "PublishAllGlobalRoleTenants",
		Description:      "Publish all tenants of a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleTenantsPublishAll,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: struct{}{},
	}.Register()

	// UnpublishAllGlobalRoleTenants
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/globalRoles/id/tenants/unpublishAll/post/",
		Name:             "UnpublishAllGlobalRoleTenants",
		Description:      "Unpublish all tenants of a global role",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathGlobalRoleTenantsUnpublishAll,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamID,
				Description: descGlobalRoleID,
				Required:    true,
			},
		},
		ResponseType: struct{}{},
	}.Register()
}
