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

//go:generate endpoint-generator -path organization.go -output org

func init() {
	// Get Organization from Vmware Cloud Director
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/orgs/get/",
		Name:             "GetOrganizationDetails",
		Description:      "Get organizations details from VMware Cloud Director",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/orgs",
		ResponseType:     itypes.APIResponseGetOrgs{},
	}.Register()

	// GetOrganization from infraAPI
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Organizations/get_api_customers_v2_0_configurations",
		Name:             "GetOrganization",
		Description:      "Get your organization information",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/configurations",
		ResponseType:     itypes.APIResponseGetOrg{},
	}.Register()

	// UpdateOrganization
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Organizations/put_api_customers_v2_0_configurations",
		Name:             "UpdateOrganization",
		Description:      "Update an existing organization",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     "/api/customers/v2.0/configurations",
		BodyRequestType:  itypes.APIRequestUpdateOrg{},
		ResponseType:     cav.Job{},
	}.Register()

	// GetCatalogAccessControl
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/catalogs/catalogUrn/accessControls/get/",
		Name:             "GetCatalogAccessControl",
		Description:      "List catalog access control grants",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCatalogAccessControl,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogURN,
				Description: descCatalogURN,
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnCatalog)
				},
			},
		},
		ResponseType: itypes.APIResponseCatalogAccessControlGrants{},
	}.Register()

	// SetCatalogAccessControl
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/catalogs/catalogUrn/accessControls/put/",
		Name:             "SetCatalogAccessControl",
		Description:      "Set catalog access control grants",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCatalogAccessControl,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogURN,
				Description: descCatalogURN,
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, urnCatalog)
				},
			},
		},
		BodyRequestType: itypes.APIRequestCatalogAccessControlGrants{},
		ResponseType:    itypes.APIResponseCatalogAccessControlGrants{},
	}.Register()
}
