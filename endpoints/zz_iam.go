/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// ListUsers - List users in organization
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ListUsers() *cav.Endpoint {
	return cav.MustGetEndpoint("ListUsers")
}
// GetUser - Get user by ID or name
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func GetUser() *cav.Endpoint {
	return cav.MustGetEndpoint("GetUser")
}
// CreateUser - Create a new user in organization
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func CreateUser() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateUser")
}
// UpdateUser - Update an existing user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UpdateUser() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateUser")
}
// DeleteUser - Delete a user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func DeleteUser() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteUser")
}
// EnableUser - Enable a user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func EnableUser() *cav.Endpoint {
	return cav.MustGetEndpoint("EnableUser")
}
// DisableUser - Disable a user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func DisableUser() *cav.Endpoint {
	return cav.MustGetEndpoint("DisableUser")
}
// UnlockUser - Unlock a user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UnlockUser() *cav.Endpoint {
	return cav.MustGetEndpoint("UnlockUser")
}
// ChangePassword - Change a user's password
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ChangePassword() *cav.Endpoint {
	return cav.MustGetEndpoint("ChangePassword")
}
// ListTokens - List tokens in organization
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ListTokens() *cav.Endpoint {
	return cav.MustGetEndpoint("ListTokens")
}
// GetToken - Get token by ID
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func GetToken() *cav.Endpoint {
	return cav.MustGetEndpoint("GetToken")
}
// CreateToken - Create a new token in organization
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func CreateToken() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateToken")
}
// UpdateToken - Update an existing token
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UpdateToken() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateToken")
}
// DeleteToken - Delete a token
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func DeleteToken() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteToken")
}
// TestLDAP - Test LDAP connection
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func TestLDAP() *cav.Endpoint {
	return cav.MustGetEndpoint("TestLDAP")
}
// SyncLDAP - Synchronize LDAP directory
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func SyncLDAP() *cav.Endpoint {
	return cav.MustGetEndpoint("SyncLDAP")
}
// SearchLDAPUsers - Search LDAP users
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func SearchLDAPUsers() *cav.Endpoint {
	return cav.MustGetEndpoint("SearchLDAPUsers")
}
// SearchLDAPGroups - Search LDAP groups
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func SearchLDAPGroups() *cav.Endpoint {
	return cav.MustGetEndpoint("SearchLDAPGroups")
}
// ListGlobalRoles - List global roles
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ListGlobalRoles() *cav.Endpoint {
	return cav.MustGetEndpoint("ListGlobalRoles")
}
// GetGlobalRole - Get global role by ID
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func GetGlobalRole() *cav.Endpoint {
	return cav.MustGetEndpoint("GetGlobalRole")
}
// CreateGlobalRole - Create a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func CreateGlobalRole() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateGlobalRole")
}
// UpdateGlobalRole - Update a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UpdateGlobalRole() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateGlobalRole")
}
// DeleteGlobalRole - Delete a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func DeleteGlobalRole() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteGlobalRole")
}
// ListGlobalRoleRights - List rights of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ListGlobalRoleRights() *cav.Endpoint {
	return cav.MustGetEndpoint("ListGlobalRoleRights")
}
// AddGlobalRoleRights - Add rights to a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func AddGlobalRoleRights() *cav.Endpoint {
	return cav.MustGetEndpoint("AddGlobalRoleRights")
}
// ReplaceGlobalRoleRights - Replace rights of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ReplaceGlobalRoleRights() *cav.Endpoint {
	return cav.MustGetEndpoint("ReplaceGlobalRoleRights")
}
// ListGlobalRoleTenants - List tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ListGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("ListGlobalRoleTenants")
}
// SetGlobalRoleTenants - Set tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func SetGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("SetGlobalRoleTenants")
}
// PublishGlobalRoleTenants - Publish tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func PublishGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("PublishGlobalRoleTenants")
}
// UnpublishGlobalRoleTenants - Unpublish tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UnpublishGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("UnpublishGlobalRoleTenants")
}
// PublishAllGlobalRoleTenants - Publish all tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func PublishAllGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("PublishAllGlobalRoleTenants")
}
// UnpublishAllGlobalRoleTenants - Unpublish all tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UnpublishAllGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("UnpublishAllGlobalRoleTenants")
}

