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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path org.go -output org

func init() {
	// List Organizations
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/orgs/get/",
		Name:             "ListOrganizationsQuerry",
		Description:      "List organizations with query parameters",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/cloudapi/1.0.0/orgs",
		BodyResponseType: itypes.ApiResponseListOrgs{},
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "Filter to apply to the list of organizations. Format: key==value. Supported keys: name, id.",
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
				Value:       "10",
			},
			{
				Name:        "page",
				Description: "The page number to retrieve.",
				Value:       "1",
			},
		},
	}.Register()

	// GetOrganization from infraAPI
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Organization",
		Name:             "GetOrganization",
		Description:      "Get your organization information",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/configurations",
		BodyRequestType:  cav.Job{},
		BodyResponseType: itypes.ApiResponseGetOrg{},
	}.Register()

	// UpdateOrganization
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Organization",
		Name:             "UpdateOrganization",
		Description:      "Update an existing organization",
		Method:           cav.MethodPUT,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/configurations",
		BodyRequestType:  itypes.ApiResponseUpdateOrg{},
		BodyResponseType: cav.Job{},
		QueryParams: []cav.QueryParam{
			{
				Name:        "fullname",
				Description: "My Organization Full Name",
			},
			{
				Name:        "description",
				Description: "This is my organization description.",
			},
			{
				Name:        "customerMail",
				Description: "my-email-adress@domain.com",
			},
			{
				Name:        "internetBillingMode",
				Description: "PAYG",
			},
			{
				Name:        "isEnabled",
				Description: "true",
			},
			{
				Name:        "isSuspended",
				Description: "false",
			},
		},
	}.Register()
}
