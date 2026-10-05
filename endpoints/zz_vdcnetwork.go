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

// ListVDCNetwork - List Org VDC Networks (routed and isolated)
//
// DocumentationURL: 
func ListVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVDCNetwork")
}
// GetVDCNetwork - Get an Org VDC Network (routed or isolated)
//
// DocumentationURL: 
func GetVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("GetVDCNetwork")
}
// CreateVDCNetwork - Create an Org VDC Network (routed or isolated)
//
// DocumentationURL: 
func CreateVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVDCNetwork")
}
// UpdateVDCNetwork - Update an Org VDC Network (routed or isolated)
//
// DocumentationURL: 
func UpdateVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVDCNetwork")
}
// DeleteVDCNetwork - Delete an Org VDC Network (routed or isolated)
//
// DocumentationURL: 
func DeleteVDCNetwork() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVDCNetwork")
}
// GetNetworkDhcpConfig - Get the DHCP configuration of an Org VDC Network
//
// DocumentationURL: 
func GetNetworkDhcpConfig() *cav.Endpoint {
	return cav.MustGetEndpoint("GetNetworkDhcpConfig")
}
// UpdateNetworkDhcpConfig - Update the DHCP configuration of an Org VDC Network
//
// DocumentationURL: 
func UpdateNetworkDhcpConfig() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateNetworkDhcpConfig")
}
// DeleteNetworkDhcpConfig - Delete the DHCP configuration of an Org VDC Network
//
// DocumentationURL: 
func DeleteNetworkDhcpConfig() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteNetworkDhcpConfig")
}

