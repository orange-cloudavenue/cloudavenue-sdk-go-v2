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

// ListVDCGroup - List VDC Groups
//
// DocumentationURL: 
func ListVDCGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("ListVDCGroup")
}
// CreateVDCGroup - Create a VDC Group
//
// DocumentationURL: 
func CreateVDCGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateVDCGroup")
}
// UpdateVDCGroup - Update a VDC Group
//
// DocumentationURL: 
func UpdateVDCGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateVDCGroup")
}
// DeleteVDCGroup - Delete a VDC Group
//
// DocumentationURL: 
func DeleteVDCGroup() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteVDCGroup")
}

