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
	// ListVdcGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/get/",
		Name:             "ListVdcGroup",
		Description:      "List Vdc Groups",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "Filter to apply to the list of Vdcs. Format: key==value. Allowed keys: name, id.",
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
		},
		ResponseType: itypes.ApiResponseListVdcGroup{},
	}.Register()

	// CreateVdcGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/post/",
		Name:             "CreateVdcGroup",
		Description:      "Create a Vdc Group",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups",
		BodyRequestType:  itypes.ApiRequestCreateVdcGroup{},
		ResponseType:     cav.Job{},
	}.Register()

	// UpdateVdcGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/put/",
		Name:             "UpdateVdcGroup",
		Description:      "Update a Vdc Group",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcGroupId",
				Description: "ID of the Vdc Group to update",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdcGroup")
				},
			},
		},
		BodyRequestType: itypes.ApiRequestUpdateVdcGroup{},
		ResponseType:    cav.Job{},
	}.Register()

	// DeleteVdcGroup
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/vdcGroups/vdcGroupId/delete/",
		Name:             "DeleteVdcGroup",
		Description:      "Delete a Vdc Group",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/vdcGroups/{vdcGroupId}",
		PathParams: []cav.PathParam{
			{
				Name:        "vdcGroupId",
				Description: "ID of the Vdc Group to delete",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdcGroup")
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
