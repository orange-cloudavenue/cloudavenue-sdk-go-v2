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
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path certificate.go -output certificate

func init() {
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/get/",
		Name:             "ListCertificate",
		Description:      "List certificate library items",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/certificateLibrary",
		QueryParams:      []cav.QueryParam{{Name: "filter", Description: "Filter to apply to the list of certificates."}, {Name: "pageSize", Description: "The number of items per page.", Value: "100"}},
		ResponseType:     itypes.ApiResponseListCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/id/get/",
		Name:             "GetCertificate",
		Description:      "Get a certificate library item",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/certificateLibrary/{id}",
		PathParams: []cav.PathParam{{
			Name:        "id",
			Description: "ID of the certificate library item",
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, "urn=certificateLibraryItem")
			},
		}},
		ResponseType: itypes.ApiResponseCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/post/",
		Name:             "CreateCertificate",
		Description:      "Create a certificate library item",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/certificateLibrary",
		BodyRequestType:  itypes.ApiRequestCertificate{},
		ResponseType:     itypes.ApiResponseCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/id/put/",
		Name:             "UpdateCertificate",
		Description:      "Update a certificate library item",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/certificateLibrary/{id}",
		PathParams: []cav.PathParam{{
			Name:        "id",
			Description: "ID of the certificate library item",
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, "urn=certificateLibraryItem")
			},
		}},
		BodyRequestType: itypes.ApiRequestCertificate{},
		ResponseType:    itypes.ApiResponseCertificate{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/id/delete/",
		Name:             "DeleteCertificate",
		Description:      "Delete a certificate library item",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/certificateLibrary/{id}",
		PathParams: []cav.PathParam{{
			Name:        "id",
			Description: "ID of the certificate library item",
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, "urn=certificateLibraryItem")
			},
		}},
	}.Register()
}
