/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package subclient

// vmware struct is the VMware subclient that implements the Client interface.
// It provides methods to interact with VMware Cloud Director API.
type vmware struct {
	client
}

// VmwareError represents the error structure returned by VMware Cloud Director API.
// It contains a message and a minor error code.
// This structure is used to parse error responses from the VMware API.
type VmwareError struct {
	// DOC API : https://developer.broadcom.com/xapis/vmware-cloud-director-api/38.1/doc/types/ErrorType.html
	Message       string `json:"message"`
	StatusCode    int    `json:"majorErrorCode"`
	StatusMessage string `json:"minorErrorCode"`
}
