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
	// Get organization details from VMware CloudAPI.
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/cloudapi/1.0.0/orgs/get/",
		Name:             "GetOrganizationDetails",
		Description:      "Get organizations details from VMware Cloud Director",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCloudAPIOrgs,
		ResponseType:     itypes.APIResponseGetOrgs{},
	}.Register()

	// Get organization configuration from Infrapi.
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API",
		Name:             "GetOrganization",
		Description:      "Get your organization information",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusConfigurations,
		ResponseType:     itypes.APIResponseGetOrg{},
	}.Register()

	// Update organization configuration through Infrapi.
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API",
		Name:             "UpdateOrganization",
		Description:      "Update an existing organization",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusConfigurations,
		BodyRequestType:  itypes.APIRequestUpdateOrg{},
		ResponseType:     cav.Job{},
	}.Register()

	// GetCatalogAccessControl
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/",
		Name:             "GetCatalogAccessControl",
		Description:      "List catalog access control grants",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCatalogAccessControl,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogUrn,
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
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/39.1/",
		Name:             "SetCatalogAccessControl",
		Description:      "Set catalog access control grants",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathCatalogAccessControl,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogUrn,
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
