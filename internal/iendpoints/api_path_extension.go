/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

// S3 endpoints
const (
	// pathOSEAssociatedTenants lists OSE tenants associated with current user.
	pathOSEAssociatedTenants = "/api/v1/core/associated-tenants"
	// pathOSEUserCredentials gets OSE credentials for one user in one tenant.
	pathOSEUserCredentials = "/api/v1/core/tenants/{organizationID}/users/{userName}/credentials"
)

// Netbackup endpoints
const (

	// pathNetBackupBase is the base path for NetBackup self-service API.
	pathNetBackupBase = "/NetBackupSelfService/Api"
	// pathNetBackupAuthToken issues or refreshes a NetBackup auth token.
	pathNetBackupAuthToken = "/NetBackupSelfService/Api/auth/token"
	// pathNetBackupInventory lists NetBackup inventory.
	pathNetBackupInventory = "/NetBackupSelfService/Api/inventory"
	// pathNetBackupMachines lists NetBackup machines.
	pathNetBackupMachines = "/NetBackupSelfService/Api/machines"
	// pathNetBackupProtectionLevels lists NetBackup protection levels.
	pathNetBackupProtectionLevels = "/NetBackupSelfService/Api/protection-levels"
	// pathNetBackupProtectionLevelByID gets one NetBackup protection level by ID.
	pathNetBackupProtectionLevelByID = "/NetBackupSelfService/Api/protection-levels/{id}"
	// pathNetBackupMachineProtect protects one NetBackup machine by ID.
	pathNetBackupMachineProtect = "/NetBackupSelfService/Api/machines/{id}/protect"
)

const (
	pathParamOrganizationID = "organizationID"
	pathParamUserName       = "userName"
)
