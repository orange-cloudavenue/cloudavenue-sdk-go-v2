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

// ListVApp - List VApps
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VApp.html 
func ListVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVApp")
}
// GetVApp - Get VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VApp.html 
func GetVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVApp")
}
// CreateVApp - Create a new VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp.html 
func CreateVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVApp")
}
// UpdateVApp - Update an existing VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-VApp.html 
func UpdateVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVApp")
}
// DeleteVApp - Delete an existing VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-VApp.html 
func DeleteVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVApp")
}
// RemoveAllNetworks - Remove all networks from a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp-RemoveAllNetworks.html 
func RemoveAllNetworks() *cav.Endpoint {
	return cav.MustGetEndpoint("RemoveAllNetworks")
}
// UndeployVApp - Undeploy a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-VApp-Undeploy.html 
func UndeployVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("UndeployVApp")
}
