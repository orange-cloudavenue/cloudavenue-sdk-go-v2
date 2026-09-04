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
		PathTemplate:     pathCertificateLibraryConsumers,
		PathParams: []cav.PathParam{{
			Name:        pathParamCertLibraryItemID,
			Description: descCertificateLibraryItemID,
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, urnCertificateLibraryItem)
			},
		}},
		QueryParams:  []cav.QueryParam{{Name: queryParamFilter, Description: "Filter to apply to the list of consumers."}, {Name: queryParamPageSize, Description: descPageSize, Value: pageSize100}},
		ResponseType: itypes.APIEntityReferences{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/certLibraryItemId/consumers/post/",
		Name:             "AddCertificateConsumer",
		Description:      "Add consumer reference to a certificate library item",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCertificateLibraryConsumers,
		PathParams: []cav.PathParam{{
			Name:        pathParamCertLibraryItemID,
			Description: descCertificateLibraryItemID,
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, urnCertificateLibraryItem)
			},
		}},
		BodyRequestType: itypes.APIEntityReference{},
		ResponseType:    itypes.APIEntityReference{},
	}.Register()

	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/ssl/certificateLibrary/certLibraryItemId/consumers/put/",
		Name:             "SetCertificateConsumers",
		Description:      "Replace consumer references for a certificate library item",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCertificateLibraryConsumers,
		PathParams: []cav.PathParam{{
			Name:        pathParamCertLibraryItemID,
			Description: descCertificateLibraryItemID,
			Required:    true,
			ValidatorFunc: func(value string) error {
				return validators.New().Var(value, urnCertificateLibraryItem)
			},
		}},
		BodyRequestType: itypes.APIEntityReferences{},
		ResponseType:    itypes.APIEntityReferences{},
	}.Register()
}
