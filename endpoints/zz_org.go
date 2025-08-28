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

// GetOrganizationDetails - Get organizations details from VMware Cloud Director
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/orgs/get/ 
func GetOrganizationDetails() *cav.Endpoint {
	return cav.MustGetEndpoint("GetOrganizationDetails")
}
// GetOrganization - Get your organization information
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Organizations/get_api_customers_v2_0_configurations 
func GetOrganization() *cav.Endpoint {
	return cav.MustGetEndpoint("GetOrganization")
}
// UpdateOrganization - Update an existing organization
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Organizations/put_api_customers_v2_0_configurations 
func UpdateOrganization() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateOrganization")
}
