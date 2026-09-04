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

// ListAppPortProfile - List Application Port Profiles
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/get/ 
func ListAppPortProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAppPortProfile")
}
// GetAppPortProfile - Get an Application Port Profile
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/appPortProfileId/get/ 
func GetAppPortProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAppPortProfile")
}
// CreateAppPortProfile - Create an Application Port Profile
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/post/ 
func CreateAppPortProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateAppPortProfile")
}
// UpdateAppPortProfile - Update an Application Port Profile
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/appPortProfileId/put/ 
func UpdateAppPortProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateAppPortProfile")
}
// DeleteAppPortProfile - Delete an Application Port Profile
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/applicationPortProfiles/appPortProfileId/delete/ 
func DeleteAppPortProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteAppPortProfile")
}
