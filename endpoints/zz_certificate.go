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

// ListCertificate - List certificate library items
//
// DocumentationURL: 
func ListCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("ListCertificate")
}
// GetCertificate - Get a certificate library item
//
// DocumentationURL: 
func GetCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("GetCertificate")
}
// CreateCertificate - Create a certificate library item
//
// DocumentationURL: 
func CreateCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateCertificate")
}
// UpdateCertificate - Update a certificate library item
//
// DocumentationURL: 
func UpdateCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateCertificate")
}
// DeleteCertificate - Delete a certificate library item
//
// DocumentationURL: 
func DeleteCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteCertificate")
}

