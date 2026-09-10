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

// ListCertificate - List certificate library items
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/get/
func ListCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("ListCertificate")
}

// GetCertificate - Get a certificate library item
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/id/get/
func GetCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("GetCertificate")
}

// CreateCertificate - Create a certificate library item
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/post/
func CreateCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateCertificate")
}

// UpdateCertificate - Update a certificate library item
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/id/put/
func UpdateCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateCertificate")
}

// DeleteCertificate - Delete a certificate library item
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/id/delete/
func DeleteCertificate() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteCertificate")
}
