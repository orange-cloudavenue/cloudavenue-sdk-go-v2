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

// ListDraasOnPremiseIp - List of on premise IP addresses allowed for this organization's draas offer
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/VCDA/getVcdaIPs 
func ListDraasOnPremiseIp() *cav.Endpoint {
	return cav.MustGetEndpoint("ListDraasOnPremiseIp")
}
// AddDraasOnPremiseIp - Allow a new on premise IP address for this organization's draas offer
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/VCDA/postVcdaIPs 
func AddDraasOnPremiseIp() *cav.Endpoint {
	return cav.MustGetEndpoint("AddDraasOnPremiseIp")
}
// RemoveDraasOnPremiseIp - Remove an on premise IP address from this organization's draas offer
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/VCDA/deleteVcdaIPs 
func RemoveDraasOnPremiseIp() *cav.Endpoint {
	return cav.MustGetEndpoint("RemoveDraasOnPremiseIp")
}
