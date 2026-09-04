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
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path vdc.go -output vdc

func init() {
	// ListVDC
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/38.1/doc/types/ReferenceType.html",
		Name:             "ListVDC",
		Description:      "List VDCs",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathQueryAPI,
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamFilter,
				Description: descFilterNameOrID,
				ValidatorFunc: func(value string) error {
					valueSplit := strings.Split(value, "==")
					if len(valueSplit) != 2 {
						return errors.New(errFilterFormatSingle)
					}

					allowedKeys := filterKeysNameOrID
					if !slices.Contains(allowedKeys, valueSplit[0]) {
						return fmt.Errorf(errFilterKeyNotAllowed, valueSplit[0])
					}

					return nil
				},
				TransformFunc: func(value string) (string, error) {
					// Add ( ) around the filter value
					return fmt.Sprintf("(%s)", value), nil
				},
			},
			{
				Name:        queryParamPageSize,
				Description: descPageSize,
				Value:       pageSize100,
			},
			{
				Name:        queryParamFormat,
				Description: descFormatResponse,
				Value:       formatRecords,
			},
			{
				Name:        queryParamType,
				Description: descTypeOfObjectQuery,
				Value:       typeOrgVDC,
			},
		},
		ResponseType: itypes.APIResponseListVDC{},
	}.Register()

	// GetVDC
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Vdc.html",
		Name:             "GetVDC",
		Description:      "Get VDC",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/vdc/{vdc-id}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCID,
				Description: descVDCID,
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnVDC)
				},
				TransformFunc: extractor.ExtractUUID,
			},
		},
		ResponseType: itypes.APIResponseGetVDC{},
	}.Register()

	// GetVDCMetadata
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VdcMetadata.html",
		Name:             "GetVDCMetadata",
		Description:      "Get VDC Metadata",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/vdc/{vdc-id}/metadata",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCID,
				Description: descVDCID,
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnVDC)
				},
				TransformFunc: extractor.ExtractUUID,
			},
		},
		ResponseType: itypes.APIResponseGetVDCMetadatas{},
	}.Register()

	// CreateVDC
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/createOrgVdc",
		Name:             "CreateVDC",
		Description:      "Create a new Org VDC",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/vdcs",
		BodyRequestType:  itypes.APIRequestCreateVDC{},
		ResponseType:     cav.Job{},
	}.Register()

	// UpdateVDC
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/updateOrgVdc",
		Name:             "UpdateVDC",
		Description:      "Update an existing Org VDC",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/vdcs/{vdc-name}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCName,
				Description: "The name of the VDC to update.",
				Required:    true,
			},
		},
		BodyRequestType: itypes.APIRequestUpdateVDC{},
		ResponseType:    cav.Job{},
	}.Register()

	// DeleteVDC
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/deleteOrgVdc",
		Name:             "DeleteVDC",
		Description:      "Delete an existing Org VDC",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/vdcs/{vdc-name}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCName,
				Description: "The name of the VDC to delete.",
				Required:    true,
			},
		},
		ResponseType: cav.Job{},
	}.Register()
}
