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

// ListFirewallGroup - List Firewall Groups (Security Groups, IP Sets, Dynamic Security Groups)
//
// DocumentationURL: 
func ListFirewallGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("ListFirewallGroup")
}
// GetFirewallGroup - Get a Firewall Group (Security Group, IP Set, or Dynamic Security Group)
//
// DocumentationURL: 
func GetFirewallGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("GetFirewallGroup")
}
// CreateFirewallGroup - Create a Firewall Group (Security Group, IP Set, or Dynamic Security Group)
//
// DocumentationURL: 
func CreateFirewallGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateFirewallGroup")
}
// UpdateFirewallGroup - Update a Firewall Group (Security Group, IP Set, or Dynamic Security Group)
//
// DocumentationURL: 
func UpdateFirewallGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateFirewallGroup")
}
// DeleteFirewallGroup - Delete a Firewall Group (Security Group, IP Set, or Dynamic Security Group)
//
// DocumentationURL: 
func DeleteFirewallGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteFirewallGroup")
}

