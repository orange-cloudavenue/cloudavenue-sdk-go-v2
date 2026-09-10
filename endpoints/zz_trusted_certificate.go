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

// ListTrustedCertificate - List trusted certificates
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/get/
func ListTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("ListTrustedCertificate")
}

// GetTrustedCertificate - Get a trusted certificate
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/get/
func GetTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("GetTrustedCertificate")
}

// CreateTrustedCertificate - Create a trusted certificate
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/post/
func CreateTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateTrustedCertificate")
}

// UpdateTrustedCertificate - Update a trusted certificate
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/put/
func UpdateTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateTrustedCertificate")
}

// DeleteTrustedCertificate - Delete a trusted certificate
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/delete/
func DeleteTrustedCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteTrustedCertificate")
}
