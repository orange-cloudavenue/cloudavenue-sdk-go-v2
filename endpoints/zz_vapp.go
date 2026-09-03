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

// ListVapp - List VApps
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VApp.html
func ListVapp() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVapp")
}

// GetVapp - Get VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VApp.html
func GetVapp() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVapp")
}

// CreateVapp - Create a new VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp.html
func CreateVapp() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVapp")
}

// UpdateVapp - Update an existing VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-VApp.html
func UpdateVapp() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVapp")
}

// DeleteVapp - Delete an existing VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-VApp.html
func DeleteVapp() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVapp")
}

// RemoveAllNetworks - Remove all networks from a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp-RemoveAllNetworks.html
func RemoveAllNetworks() *cav.Endpoint {
	return cav.MustGetEndpoint("RemoveAllNetworks")
}

// UndeployVapp - Undeploy a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp-Undeploy.html
func UndeployVapp() *cav.Endpoint {
	return cav.MustGetEndpoint("UndeployVapp")
}
