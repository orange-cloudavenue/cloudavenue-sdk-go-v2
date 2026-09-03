/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path trusted_certificate.go -output trusted_certificate

func init() {
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/get/",
		Name:             "ListTrustedCertificate",
		Description:      "List trusted certificates",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/trustedCertificates",
		QueryParams:      []cav.QueryParam{{Name: "filter", Description: "Filter to apply to the list of trusted certificates."}, {Name: "pageSize", Description: "The number of items per page.", Value: "100"}},
		ResponseType:     itypes.ApiResponseListTrustedCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/get/",
		Name:             "GetTrustedCertificate",
		Description:      "Get a trusted certificate",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/trustedCertificates/{trustedCertificate}",
		PathParams: []cav.PathParam{{
			Name:        "trustedCertificate",
			Description: "ID of the trusted certificate",
			Required:    true,
		}},
		ResponseType: itypes.ApiResponseTrustedCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/post/",
		Name:             "CreateTrustedCertificate",
		Description:      "Create a trusted certificate",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/trustedCertificates",
		BodyRequestType:  itypes.ApiRequestTrustedCertificate{},
		ResponseType:     itypes.ApiResponseTrustedCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/put/",
		Name:             "UpdateTrustedCertificate",
		Description:      "Update a trusted certificate",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/trustedCertificates/{trustedCertificate}",
		PathParams: []cav.PathParam{{
			Name:        "trustedCertificate",
			Description: "ID of the trusted certificate",
			Required:    true,
		}},
		BodyRequestType: itypes.ApiRequestTrustedCertificate{},
		ResponseType:    itypes.ApiResponseTrustedCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/delete/",
		Name:             "DeleteTrustedCertificate",
		Description:      "Delete a trusted certificate",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/trustedCertificates/{trustedCertificate}",
		PathParams: []cav.PathParam{{
			Name:        "trustedCertificate",
			Description: "ID of the trusted certificate",
			Required:    true,
		}},
	}.Register()
}
