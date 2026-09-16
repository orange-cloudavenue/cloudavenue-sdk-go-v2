/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
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
