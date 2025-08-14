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

// ListVdc - List VDCs
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/38.1/doc/types/ReferenceType.html 
func ListVdc() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVdc")
}
// GetVdc - Get VDC
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Vdc.html 
func GetVdc() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVdc")
}
// GetVdcMetadata - Get VDC Metadata
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VdcMetadata.html 
func GetVdcMetadata() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVdcMetadata")
}
// CreateVdc - Create a new Org VDC
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/vDC/createOrgVdc 
func CreateVdc() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVdc")
}
// UpdateVdc - Update an existing Org VDC
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/vDC/updateOrgVdc 
func UpdateVdc() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVdc")
}
// DeleteVdc - Delete an existing Org VDC
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/vDC/deleteOrgVdc 
func DeleteVdc() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVdc")
}
