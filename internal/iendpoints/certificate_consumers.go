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

//go:generate endpoint-generator -path certificate_consumers.go -output certificate_consumers

func init() {
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/certLibraryItemId/consumers/get/",
		Name:             "ListCertificateConsumers",
		Description:      "List consumers of a certificate library item",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/certificateLibrary/{certLibraryItemId}/consumers",
		PathParams: []cav.PathParam{{
			Name:        "certLibraryItemId",
			Description: "ID of the certificate library item",
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, "urn=certificateLibraryItem")
			},
		}},
		QueryParams:  []cav.QueryParam{{Name: "filter", Description: "Filter to apply to the list of consumers."}, {Name: "pageSize", Description: "The number of items per page.", Value: "100"}},
		ResponseType: itypes.ApiEntityReferences{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/certLibraryItemId/consumers/post/",
		Name:             "AddCertificateConsumer",
		Description:      "Add consumer reference to a certificate library item",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/certificateLibrary/{certLibraryItemId}/consumers",
		PathParams: []cav.PathParam{{
			Name:        "certLibraryItemId",
			Description: "ID of the certificate library item",
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, "urn=certificateLibraryItem")
			},
		}},
		BodyRequestType: itypes.ApiEntityReference{},
		ResponseType:    itypes.ApiEntityReference{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/certLibraryItemId/consumers/put/",
		Name:             "SetCertificateConsumers",
		Description:      "Replace consumer references for a certificate library item",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/ssl/certificateLibrary/{certLibraryItemId}/consumers",
		PathParams: []cav.PathParam{{
			Name:        "certLibraryItemId",
			Description: "ID of the certificate library item",
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, "urn=certificateLibraryItem")
			},
		}},
		BodyRequestType: itypes.ApiEntityReferences{},
		ResponseType:    itypes.ApiEntityReferences{},
	}.Register()
}
