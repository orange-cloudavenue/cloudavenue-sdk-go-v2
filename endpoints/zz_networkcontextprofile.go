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

// ListNetworkContextProfile - List Network Context Profiles
//
// DocumentationURL: 
func ListNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("ListNetworkContextProfile")
}
// GetNetworkContextProfile - Get a Network Context Profile
//
// DocumentationURL: 
func GetNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("GetNetworkContextProfile")
}
// CreateNetworkContextProfile - Create a Network Context Profile
//
// DocumentationURL: 
func CreateNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateNetworkContextProfile")
}
// UpdateNetworkContextProfile - Update a Network Context Profile
//
// DocumentationURL: 
func UpdateNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateNetworkContextProfile")
}
// DeleteNetworkContextProfile - Delete a Network Context Profile
//
// DocumentationURL: 
func DeleteNetworkContextProfile() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteNetworkContextProfile")
}
// GetNetworkContextProfileAttributes - Get the static reference catalog (App IDs, Domain Names) of attributes usable in Network Context Profiles
//
// DocumentationURL: 
func GetNetworkContextProfileAttributes() *cav.Endpoint {
	return cav.MustGetEndpoint("GetNetworkContextProfileAttributes")
}

