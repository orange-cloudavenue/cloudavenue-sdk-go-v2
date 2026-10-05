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

// ListCertificateConsumers - List consumers of a certificate library item
//
// DocumentationURL: 
func ListCertificateConsumers() *cav.Endpoint {
	return cav.MustGetEndpoint("ListCertificateConsumers")
}
// AddCertificateConsumer - Add consumer reference to a certificate library item
//
// DocumentationURL: 
func AddCertificateConsumer() *cav.Endpoint {
	return cav.MustGetEndpoint("AddCertificateConsumer")
}
// SetCertificateConsumers - Replace consumer references for a certificate library item
//
// DocumentationURL: 
func SetCertificateConsumers() *cav.Endpoint {
	return cav.MustGetEndpoint("SetCertificateConsumers")
}

