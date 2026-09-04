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

// ListVDCNetwork - List Org VDC Networks (routed and isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/get/ 
func ListVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVDCNetwork")
}
// GetVDCNetwork - Get an Org VDC Network (routed or isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/get/ 
func GetVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVDCNetwork")
}
// CreateVDCNetwork - Create an Org VDC Network (routed or isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/post/ 
func CreateVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVDCNetwork")
}
// UpdateVDCNetwork - Update an Org VDC Network (routed or isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/put/ 
func UpdateVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVDCNetwork")
}
// DeleteVDCNetwork - Delete an Org VDC Network (routed or isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/delete/ 
func DeleteVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVDCNetwork")
}
