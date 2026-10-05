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
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func ListTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("ListTrustedCertificate")
}
// GetTrustedCertificate - Get a trusted certificate
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func GetTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("GetTrustedCertificate")
}
// CreateTrustedCertificate - Create a trusted certificate
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func CreateTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateTrustedCertificate")
}
// UpdateTrustedCertificate - Update a trusted certificate
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func UpdateTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateTrustedCertificate")
}
// DeleteTrustedCertificate - Delete a trusted certificate
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/
func DeleteTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteTrustedCertificate")
}

