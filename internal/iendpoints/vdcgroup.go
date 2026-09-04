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

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path vdcgroup.go -output vdcgroup

func init() {
	// ListVDCGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/get/",
		Name:             "ListVDCGroup",
		Description:      "List VDC Groups",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups",
		QueryParams: []cav.QueryParam{
			{
				Name:        queryParamFilter,
				Description: "Filter to apply to the list of Vdcs. Format: key==value. Allowed keys: name, id.",
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
		},
		ResponseType: itypes.APIResponseListVDCGroup{},
	}.Register()

	// CreateVDCGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/post/",
		Name:             "CreateVDCGroup",
		Description:      "Create a VDC Group",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups",
		BodyRequestType:  itypes.APIRequestCreateVDCGroup{},
		ResponseType:     cav.Job{},
	}.Register()

	// UpdateVDCGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/put/",
		Name:             "UpdateVDCGroup",
		Description:      "Update a VDC Group",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCGroupID,
				Description: "ID of the Vdc Group to update",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnVDCGroup)
				},
			},
		},
		BodyRequestType: itypes.APIRequestUpdateVDCGroup{},
		ResponseType:    cav.Job{},
	}.Register()

	// DeleteVDCGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/delete/",
		Name:             "DeleteVDCGroup",
		Description:      "Delete a VDC Group",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCGroupID,
				Description: "ID of the Vdc Group to delete",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnVDCGroup)
				},
			},
		},
		QueryParams: []cav.QueryParam{
			{
				Name:        "force",
				Description: "Force delete the Vdc Group",
				Required:    false,
			},
		},
		ResponseType: cav.Job{},
	}.Register()
}
