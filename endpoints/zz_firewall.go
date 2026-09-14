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

// GetDfwPolicies - Get the Distributed Firewall policies (enabled state and default policy) of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/get/
func GetDfwPolicies() *cav.Endpoint {
	return cav.MustGetEndpoint("GetDfwPolicies")
}

// UpdateDfwPolicies - Update the Distributed Firewall policies (enabled state) of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/put/
func UpdateDfwPolicies() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateDfwPolicies")
}

// UpdateDfwDefaultPolicy - Update the default Distributed Firewall policy of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/put/
func UpdateDfwDefaultPolicy() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateDfwDefaultPolicy")
}

// GetDfwRules - Get the Distributed Firewall rules of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/rules/get/
func GetDfwRules() *cav.Endpoint {
	return cav.MustGetEndpoint("GetDfwRules")
}

// UpdateDfwRules - Replace (bulk) the Distributed Firewall rules of a VDC Group
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/dfwPolicies/default/rules/put/
func UpdateDfwRules() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateDfwRules")
}
