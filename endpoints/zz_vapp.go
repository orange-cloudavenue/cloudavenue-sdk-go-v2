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
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-ExecuteQuery.html
func ListVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVApp")
}

// DeployVApp - Deploy a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-DeployVApp.html
func DeployVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("DeployVApp")
}

// PowerOffVApp - Power off a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-PowerOffVApp.html
func PowerOffVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("PowerOffVApp")
}

// PowerOnVApp - Power on a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-PowerOnVApp.html
func PowerOnVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("PowerOnVApp")
}

// RebootVApp - Reboot a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-RebootVApp.html
func RebootVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("RebootVApp")
}

// ResetVApp - Reset a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ResetVApp.html
func ResetVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("ResetVApp")
}

// ShutdownVApp - Shut down a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ShutdownVApp.html
func ShutdownVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("ShutdownVApp")
}

// SuspendVApp - Suspend a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-SuspendVApp.html
func SuspendVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("SuspendVApp")
}

// GetVApp - Get VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-VApp.html
func GetVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVApp")
}

// ComposeVApp - Compose a new VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ComposeVApp.html
func ComposeVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("ComposeVApp")
}

// InstantiateVAppTemplate - Instantiate a VApp from a template
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-InstantiateVAppTemplate.html
func InstantiateVAppTemplate() *cav.Endpoint {
	return cav.MustGetEndpoint("InstantiateVAppTemplate")
}

// UpdateVApp - Update an existing VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-Vm.html
func UpdateVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVApp")
}

// DeleteVApp - Delete an existing VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/DELETE-VApp.html
func DeleteVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVApp")
}

// RemoveAllNetworks - Remove all networks from a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-NetworkConfigSection-vApp.html
func RemoveAllNetworks() *cav.Endpoint {
	return cav.MustGetEndpoint("RemoveAllNetworks")
}

// UndeployVApp - Undeploy a VApp
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-UndeployVApp.html
func UndeployVApp() *cav.Endpoint {
	return cav.MustGetEndpoint("UndeployVApp")
}

// GetVAppLeaseSettings - Get org/VDC-level lease settings
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-VAppLeaseSettings.html
func GetVAppLeaseSettings() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVAppLeaseSettings")
}

// UpdateVAppLeaseSettings - Update org/VDC-level lease settings
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-VAppLeaseSettings.html
func UpdateVAppLeaseSettings() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVAppLeaseSettings")
}
