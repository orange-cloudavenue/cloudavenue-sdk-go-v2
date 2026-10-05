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

// ListVApp - List VApps
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ListVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVApp")
}
// GetVApp - Get VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func GetVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVApp")
}
// CreateVApp - Create a new VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func CreateVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVApp")
}
// UpdateVApp - Update an existing VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UpdateVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVApp")
}
// DeleteVApp - Delete an existing VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func DeleteVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVApp")
}
// RemoveAllNetworks - Remove all networks from a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func RemoveAllNetworks() *cav.Endpoint {
	return cav.MustGetEndpoint("RemoveAllNetworks")
}
// UndeployVApp - Undeploy a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UndeployVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("UndeployVApp")
}
// GetVAppLeaseSettings - Get org/VDC-level lease settings
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func GetVAppLeaseSettings() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVAppLeaseSettings")
}
// UpdateVAppLeaseSettings - Update org/VDC-level lease settings
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UpdateVAppLeaseSettings() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVAppLeaseSettings")
}

