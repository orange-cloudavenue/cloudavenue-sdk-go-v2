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

// GetDFWPolicies - Get the Distributed Firewall policies (enabled state and default policy) of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/get/ 
func GetDFWPolicies() *cav.Endpoint {
	return cav.MustGetEndpoint("GetDFWPolicies")
}
// UpdateDFWPolicies - Update the Distributed Firewall policies (enabled state) of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/put/ 
func UpdateDFWPolicies() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateDFWPolicies")
}
// UpdateDFWDefaultPolicy - Update the default Distributed Firewall policy of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/put/ 
func UpdateDFWDefaultPolicy() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateDFWDefaultPolicy")
}
// GetDFWRules - Get the Distributed Firewall rules of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/rules/get/ 
func GetDFWRules() *cav.Endpoint {
	return cav.MustGetEndpoint("GetDFWRules")
}
// UpdateDFWRules - Replace (bulk) the Distributed Firewall rules of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/rules/put/ 
func UpdateDFWRules() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateDFWRules")
}
