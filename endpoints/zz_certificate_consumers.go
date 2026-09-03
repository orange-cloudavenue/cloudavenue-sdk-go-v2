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

// ListCertificateConsumers - List consumers of a certificate library item
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/certLibraryItemId/consumers/get/
func ListCertificateConsumers() *cav.Endpoint {
	return cav.MustGetEndpoint("ListCertificateConsumers")
}

// AddCertificateConsumer - Add consumer reference to a certificate library item
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/certLibraryItemId/consumers/post/
func AddCertificateConsumer() *cav.Endpoint {
	return cav.MustGetEndpoint("AddCertificateConsumer")
}

// SetCertificateConsumers - Replace consumer references for a certificate library item
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/certLibraryItemId/consumers/put/
func SetCertificateConsumers() *cav.Endpoint {
	return cav.MustGetEndpoint("SetCertificateConsumers")
}
