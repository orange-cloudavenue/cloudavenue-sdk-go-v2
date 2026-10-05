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

// ListTrustedCertificate - List trusted certificates
//
// DocumentationURL: 
func ListTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("ListTrustedCertificate")
}
// GetTrustedCertificate - Get a trusted certificate
//
// DocumentationURL: 
func GetTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("GetTrustedCertificate")
}
// CreateTrustedCertificate - Create a trusted certificate
//
// DocumentationURL: 
func CreateTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateTrustedCertificate")
}
// UpdateTrustedCertificate - Update a trusted certificate
//
// DocumentationURL: 
func UpdateTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateTrustedCertificate")
}
// DeleteTrustedCertificate - Delete a trusted certificate
//
// DocumentationURL: 
func DeleteTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteTrustedCertificate")
}

