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

// ListVdcNetwork - List Org VDC Networks (routed and isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/get/
func ListVdcNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVdcNetwork")
}

// GetVdcNetwork - Get an Org VDC Network (routed or isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/get/
func GetVdcNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVdcNetwork")
}

// CreateVdcNetwork - Create an Org VDC Network (routed or isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/post/
func CreateVdcNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVdcNetwork")
}

// UpdateVdcNetwork - Update an Org VDC Network (routed or isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/put/
func UpdateVdcNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVdcNetwork")
}

// DeleteVdcNetwork - Delete an Org VDC Network (routed or isolated)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/orgVdcNetworks/vdcNetworkId/delete/
func DeleteVdcNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVdcNetwork")
}
