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

// ListNetworkContextProfile - List Network Context Profiles
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/get/
func ListNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("ListNetworkContextProfile")
}

// GetNetworkContextProfile - Get a Network Context Profile
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/networkContextProfileId/get/
func GetNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("GetNetworkContextProfile")
}

// CreateNetworkContextProfile - Create a Network Context Profile
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/post/
func CreateNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateNetworkContextProfile")
}

// UpdateNetworkContextProfile - Update a Network Context Profile
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/networkContextProfileId/put/
func UpdateNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateNetworkContextProfile")
}

// DeleteNetworkContextProfile - Delete a Network Context Profile
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/networkContextProfileId/delete/
func DeleteNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteNetworkContextProfile")
}

// GetNetworkContextProfileAttributes - Get the static reference catalog (App IDs, Domain Names) of attributes usable in Network Context Profiles
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/networkContextProfiles/attributes/get/
func GetNetworkContextProfileAttributes() *cav.Endpoint {
	return cav.MustGetEndpoint("GetNetworkContextProfileAttributes")
}
