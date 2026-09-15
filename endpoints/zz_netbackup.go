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

// GetNetbackupToken - Get NetBackup OAuth2 token using password grant
//
// DocumentationURL: https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/auth/token 
func GetNetbackupToken() *cav.Endpoint {
	return cav.MustGetEndpoint("GetNetbackupToken")
}
// RefreshNetbackupToken - Refresh NetBackup OAuth2 token using refresh token grant
//
// DocumentationURL: https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/auth/token 
func RefreshNetbackupToken() *cav.Endpoint {
	return cav.MustGetEndpoint("RefreshNetbackupToken")
}
// ListNetbackupInventory - List NetBackup inventory
//
// DocumentationURL: https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/inventory 
func ListNetbackupInventory() *cav.Endpoint {
	return cav.MustGetEndpoint("ListNetbackupInventory")
}
// ListNetbackupMachines - List NetBackup machines
//
// DocumentationURL: https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/machines 
func ListNetbackupMachines() *cav.Endpoint {
	return cav.MustGetEndpoint("ListNetbackupMachines")
}
// GetNetbackupProtectionLevel - Get NetBackup protection level by ID
//
// DocumentationURL: https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/protection-levels/{id} 
func GetNetbackupProtectionLevel() *cav.Endpoint {
	return cav.MustGetEndpoint("GetNetbackupProtectionLevel")
}
// ListNetbackupProtectionLevels - List NetBackup protection levels
//
// DocumentationURL: https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/protection-levels 
func ListNetbackupProtectionLevels() *cav.Endpoint {
	return cav.MustGetEndpoint("ListNetbackupProtectionLevels")
}
// ProtectNetbackupMachine - Protect a NetBackup machine
//
// DocumentationURL: https://backup.cloudavenue.orange-business.com/NetBackupSelfService/Api/machines/{id}/protect 
func ProtectNetbackupMachine() *cav.Endpoint {
	return cav.MustGetEndpoint("ProtectNetbackupMachine")
}
