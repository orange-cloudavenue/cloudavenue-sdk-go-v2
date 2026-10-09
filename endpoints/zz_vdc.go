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

// ListVDC - List VDCs
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func ListVDC() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVDC")
}
// GetVDC - Get VDC
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func GetVDC() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVDC")
}
// GetVDCMetadata - Get VDC Metadata
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-VdcMetadata.html
func GetVDCMetadata() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVDCMetadata")
}
// CreateVDC - Create a new Org VDC
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func CreateVDC() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVDC")
}
// UpdateVDC - Update an existing Org VDC
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func UpdateVDC() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVDC")
}
// DeleteVDC - Delete an existing Org VDC
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func DeleteVDC() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVDC")
}

