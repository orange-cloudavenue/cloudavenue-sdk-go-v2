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

// ListVdcGroup - List Vdc Groups
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/get/
func ListVdcGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVdcGroup")
}

// CreateVdcGroup - Create a Vdc Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/post/
func CreateVdcGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVdcGroup")
}

// UpdateVdcGroup - Update a Vdc Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/put/
func UpdateVdcGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVdcGroup")
}

// DeleteVdcGroup - Delete a Vdc Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/delete/
func DeleteVdcGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVdcGroup")
}
