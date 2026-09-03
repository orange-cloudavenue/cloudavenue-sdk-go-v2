/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

// NetBackup types

type (
	APIResponseNetbackupInventory struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	APIResponseNetbackupMachine struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Hostname    string `json:"hostname"`
		IPAddress   string `json:"ipAddress"`
		Status      string `json:"status"`
	}

	APIResponseNetbackupProtectionLevel struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Retention   int    `json:"retention"`
	}

	APIRequestNetbackupProtectMachine struct {
		ProtectionLevelID string `json:"protectionLevelId"`
	}

	APIResponseNetbackupProtectMachine struct {
		ID          string `json:"id"`
		MachineID   string `json:"machineId"`
		Status      string `json:"status"`
	}
)

type (
	APIResponseListNetbackupInventory []APIResponseNetbackupInventory
	APIResponseListNetbackupMachines  []APIResponseNetbackupMachine
	APIResponseListNetbackupProtectionLevels []APIResponseNetbackupProtectionLevel
)
