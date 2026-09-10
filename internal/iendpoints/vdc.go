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
		Name:             "ListVdc",
		Description:      "List VDCs",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/query",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "Filter to apply to the list of VDCs. Format: key==value. Allowed keys: name, id.",
				ValidatorFunc: func(value string) error {
					valueSplit := strings.Split(value, "==")
					if len(valueSplit) != 2 {
						return errors.New("filter must be in the format 'key==value'")
					}

					allowedKeys := []string{"name", "id"}
					if !slices.Contains(allowedKeys, valueSplit[0]) {
						return fmt.Errorf("filter key '%s' is not allowed", valueSplit[0])
					}

					return nil
				},
				TransformFunc: func(value string) (string, error) {
					// Add ( ) around the filter value
					return fmt.Sprintf("(%s)", value), nil
				},
			},
			{
				Name:        "pageSize",
				Description: "The number of items per page.",
				Value:       "100",
			},
			{
				Name:        "format",
				Description: "The format of the response.",
				Value:       "records",
			},
			{
				Name:        "type",
				Description: "The type of object to query",
				Value:       "orgVdc",
			},
		},
		ResponseType: itypes.ApiResponseListVDC{},
	}.Register()

	// GetVDC
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-Vdc.html",
		Name:             "GetVdc",
		Description:      "Get VDC",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/vdc/{vdc-id}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-id",
				Description: "The ID of the VDC.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdc")
				},
				TransformFunc: func(value string) (string, error) {
					// vdc-id require UUID format and not urn format
					return extractor.ExtractUUID(value)
				},
			},
		},
		ResponseType: itypes.ApiResponseGetVDC{},
	}.Register()

	// GetVDCMetadata
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-VdcMetadata.html",
		Name:             "GetVdcMetadata",
		Description:      "Get VDC Metadata",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/vdc/{vdc-id}/metadata",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-id",
				Description: "The ID of the VDC.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdc")
				},
				TransformFunc: func(value string) (string, error) {
					// vdc-id require UUID format and not urn format
					return extractor.ExtractUUID(value)
				},
			},
		},
		ResponseType: itypes.ApiResponseGetVDCMetadatas{},
	}.Register()

	// CreateVdc
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/createOrgVdc",
		Name:             "CreateVdc",
		Description:      "Create a new Org VDC",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/vdcs",
		BodyRequestType:  itypes.ApiRequestCreateVDC{},
		ResponseType:     cav.Job{},
	}.Register()

	// UpdateVdc
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/updateOrgVdc",
		Name:             "UpdateVdc",
		Description:      "Update an existing Org VDC",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/vdcs/{vdc-name}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-name",
				Description: "The name of the VDC to update.",
				Required:    true,
			},
		},
		BodyRequestType: itypes.ApiRequestUpdateVDC{},
		ResponseType:    cav.Job{},
	}.Register()

	// DeleteVdc
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/vDC/deleteOrgVdc",
		Name:             "DeleteVdc",
		Description:      "Delete an existing Org VDC",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/vdcs/{vdc-name}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdc-name",
				Description: "The name of the VDC to delete.",
				Required:    true,
			},
		},
		ResponseType: cav.Job{},
	}.Register()
}
