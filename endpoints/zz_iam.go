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
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Users.html
func ListUsers() *cav.Endpoint {
	return cav.MustGetEndpoint("ListUsers")
}
// GetUser - Get user by ID or name
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-User.html
func GetUser() *cav.Endpoint {
	return cav.MustGetEndpoint("GetUser")
}
// CreateUser - Create a new user in organization
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-User.html
func CreateUser() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateUser")
}
// UpdateUser - Update an existing user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-User.html
func UpdateUser() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateUser")
}
// DeleteUser - Delete a user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-User.html
func DeleteUser() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteUser")
}
// EnableUser - Enable a user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-UserEnable.html
func EnableUser() *cav.Endpoint {
	return cav.MustGetEndpoint("EnableUser")
}
// DisableUser - Disable a user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-UserDisable.html
func DisableUser() *cav.Endpoint {
	return cav.MustGetEndpoint("DisableUser")
}
// UnlockUser - Unlock a user
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-UserUnlock.html
func UnlockUser() *cav.Endpoint {
	return cav.MustGetEndpoint("UnlockUser")
}
// ChangePassword - Change a user's password
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-UserChangePassword.html
func ChangePassword() *cav.Endpoint {
	return cav.MustGetEndpoint("ChangePassword")
}
// ListTokens - List tokens in organization
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Tokens.html
func ListTokens() *cav.Endpoint {
	return cav.MustGetEndpoint("ListTokens")
}
// GetToken - Get token by ID
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Token.html
func GetToken() *cav.Endpoint {
	return cav.MustGetEndpoint("GetToken")
}
// CreateToken - Create a new token in organization
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-Token.html
func CreateToken() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateToken")
}
// UpdateToken - Update an existing token
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-Token.html
func UpdateToken() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateToken")
}
// DeleteToken - Delete a token
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-Token.html
func DeleteToken() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteToken")
}
// TestLDAP - Test LDAP connection
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-LDAPTest.html
func TestLDAP() *cav.Endpoint {
	return cav.MustGetEndpoint("TestLDAP")
}
// SyncLDAP - Synchronize LDAP directory
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-LDAPSync.html
func SyncLDAP() *cav.Endpoint {
	return cav.MustGetEndpoint("SyncLDAP")
}
// SearchLDAPUsers - Search LDAP users
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-LDAPSearchUser.html
func SearchLDAPUsers() *cav.Endpoint {
	return cav.MustGetEndpoint("SearchLDAPUsers")
}
// SearchLDAPGroups - Search LDAP groups
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-LDAPSearchGroup.html
func SearchLDAPGroups() *cav.Endpoint {
	return cav.MustGetEndpoint("SearchLDAPGroups")
}
// ListGlobalRoles - List global roles
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-GlobalRoles.html
func ListGlobalRoles() *cav.Endpoint {
	return cav.MustGetEndpoint("ListGlobalRoles")
}
// GetGlobalRole - Get global role by ID
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-GlobalRole.html
func GetGlobalRole() *cav.Endpoint {
	return cav.MustGetEndpoint("GetGlobalRole")
}
// CreateGlobalRole - Create a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRole.html
func CreateGlobalRole() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateGlobalRole")
}
// UpdateGlobalRole - Update a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-GlobalRole.html
func UpdateGlobalRole() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateGlobalRole")
}
// DeleteGlobalRole - Delete a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-GlobalRole.html
func DeleteGlobalRole() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteGlobalRole")
}
// ListGlobalRoleRights - List rights of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-GlobalRoleRights.html
func ListGlobalRoleRights() *cav.Endpoint {
	return cav.MustGetEndpoint("ListGlobalRoleRights")
}
// AddGlobalRoleRights - Add rights to a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleRights.html
func AddGlobalRoleRights() *cav.Endpoint {
	return cav.MustGetEndpoint("AddGlobalRoleRights")
}
// ReplaceGlobalRoleRights - Replace rights of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-GlobalRoleRights.html
func ReplaceGlobalRoleRights() *cav.Endpoint {
	return cav.MustGetEndpoint("ReplaceGlobalRoleRights")
}
// ListGlobalRoleTenants - List tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-GlobalRoleTenants.html
func ListGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("ListGlobalRoleTenants")
}
// SetGlobalRoleTenants - Set tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-GlobalRoleTenants.html
func SetGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("SetGlobalRoleTenants")
}
// PublishGlobalRoleTenants - Publish tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleTenantsPublish.html
func PublishGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("PublishGlobalRoleTenants")
}
// UnpublishGlobalRoleTenants - Unpublish tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleTenantsUnpublish.html
func UnpublishGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("UnpublishGlobalRoleTenants")
}
// PublishAllGlobalRoleTenants - Publish all tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleTenantsPublishAll.html
func PublishAllGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("PublishAllGlobalRoleTenants")
}
// UnpublishAllGlobalRoleTenants - Unpublish all tenants of a global role
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-GlobalRoleTenantsUnpublishAll.html
func UnpublishAllGlobalRoleTenants() *cav.Endpoint {
	return cav.MustGetEndpoint("UnpublishAllGlobalRoleTenants")
}

