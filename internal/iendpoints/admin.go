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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path admin.go -output admin

func init() {
	// ListAdminOrgs is the read-only, XML-backed AdminOrg organization view.
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-OrganizationsFromQuery.html",
		Name:             "ListAdminOrgs",
		Description:      "List organizations (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminOrgs,
		ResponseType:     itypes.AdminOrgs{},
	}.Register()

	// GetAdminOrg is the read-only, XML-backed AdminOrg organization view.
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-Organization-AdminView.html",
		Name:             "GetAdminOrg",
		Description:      "Get an organization by ID (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminOrg,
		PathParams:       []cav.PathParam{{Name: pathParamOrgID, Description: descOrgID, Required: true}},
		ResponseType:     itypes.AdminOrg{},
	}.Register()

	const descCatalogID = "Catalog ID"
	// ListAdminCatalogs
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-CatalogsFromQuery.html",
		Name:             "ListAdminCatalogs",
		Description:      "List catalogs (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminCatalogs,
		ResponseType:     itypes.AdminCatalogs{},
	}.Register()

	// GetAdminCatalog
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-Catalog-AdminView.html",
		Name:             "GetAdminCatalog",
		Description:      "Get a catalog by ID (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminCatalog,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
		ResponseType: itypes.AdminCatalog{},
	}.Register()

	// CreateAdminCatalog
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-SecuredCreateCatalog.html",
		Name:             "CreateAdminCatalog",
		Description:      "Create a catalog (admin scope)",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminOrgCatalogs,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgIDAdmin,
				Required:    true,
			},
		},
		BodyRequestType: itypes.AdminCatalogRequest{},
		ResponseType:    itypes.AdminCatalog{},
	}.Register()

	// UpdateAdminCatalog
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-Catalog.html",
		Name:             "UpdateAdminCatalog",
		Description:      "Update a catalog (admin scope)",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminCatalog,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.AdminCatalogRequest{},
		ResponseType:    itypes.AdminCatalog{},
	}.Register()

	// DeleteAdminCatalog
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/DELETE-Catalog.html",
		Name:             "DeleteAdminCatalog",
		Description:      "Delete a catalog (admin scope)",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminCatalog,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
	}.Register()

	// GetAdminCatalogACL
	cav.Endpoint{
		DocumentationURL: docURLVMware,
		Name:             "GetAdminCatalogACL",
		Description:      "Get catalog ACL (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminCatalogControlAccess,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
		ResponseType: itypes.ControlAccessParams{},
	}.Register()

	// SetAdminCatalogACL
	cav.Endpoint{
		DocumentationURL: docURLVMware,
		Name:             "SetAdminCatalogACL",
		Description:      "Set catalog ACL (admin scope)",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		Headers:          map[string]string{headerAccept: headerXML, headerContentType: headerXML},
		PathTemplate:     pathAdminCatalogControlAccessAction,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.ControlAccessParams{},
		ResponseType:    itypes.ControlAccessParams{},
	}.Register()
}
