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
		QueryParams:      []cav.QueryParam{{Name: queryParamFilter, Description: "Filter to apply to the list of trusted certificates."}, {Name: queryParamPageSize, Description: descPageSize, Value: pageSize100}},
		ResponseType:     itypes.APIResponseListTrustedCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/get/",
		Name:             "GetTrustedCertificate",
		Description:      "Get a trusted certificate",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTrustedCertificates,
		PathParams: []cav.PathParam{{
			Name:        pathParamTrustedCertificate,
			Description: descTrustedCertificateID,
			Required:    true,
		}},
		ResponseType: itypes.APIResponseTrustedCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/post/",
		Name:             "CreateTrustedCertificate",
		Description:      "Create a trusted certificate",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/trustedCertificates",
		BodyRequestType:  itypes.APIRequestTrustedCertificate{},
		ResponseType:     itypes.APIResponseTrustedCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/put/",
		Name:             "UpdateTrustedCertificate",
		Description:      "Update a trusted certificate",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTrustedCertificates,
		PathParams: []cav.PathParam{{
			Name:        pathParamTrustedCertificate,
			Description: descTrustedCertificateID,
			Required:    true,
		}},
		BodyRequestType: itypes.APIRequestTrustedCertificate{},
		ResponseType:    itypes.APIResponseTrustedCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/trustedCertificates/trustedCertificate/delete/",
		Name:             "DeleteTrustedCertificate",
		Description:      "Delete a trusted certificate",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathTrustedCertificates,
		PathParams: []cav.PathParam{{
			Name:        pathParamTrustedCertificate,
			Description: descTrustedCertificateID,
			Required:    true,
		}},
	}.Register()
}
