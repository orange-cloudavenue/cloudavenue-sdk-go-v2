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
	const pathAdminOrgs = "/api/admin/orgs"
	const pathAdminOrg = "/api/admin/org/{orgId}"
	const pathAdminCatalog = "/api/admin/catalog/{catalogId}"
	const pathAdminCatalogID = "catalogId"
	const descCatalogID = "Catalog ID"

	// ListAdminOrgs
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminOrgs.html",
		Name:             "ListAdminOrgs",
		Description:      "List organizations (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrgs,
		ResponseType:     itypes.AdminOrgs{},
	}.Register()

	// GetAdminOrg
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminOrg.html",
		Name:             "GetAdminOrg",
		Description:      "Get an organization by ID (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminOrg,
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgID,
				Description: descOrgID,
				Required:    true,
			},
		},
		ResponseType: itypes.AdminOrg{},
	}.Register()

	// ListAdminVDCs
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminVdcs.html",
		Name:             "ListAdminVDCs",
		Description:      "List VDCs (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/admin/vdcs",
		ResponseType:     itypes.AdminVDCs{},
	}.Register()

	// GetAdminVDC
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminVdc.html",
		Name:             "GetAdminVDC",
		Description:      "Get a VDC by ID (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/admin/vdc/{vdcId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamVDCIDAdmin,
				Description: descVDCIDAdmin,
				Required:    true,
			},
		},
		ResponseType: itypes.AdminVDC{},
	}.Register()

	// ListAdminCatalogs
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminCatalogs.html",
		Name:             "ListAdminCatalogs",
		Description:      "List catalogs (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/admin/catalogs",
		ResponseType:     itypes.AdminCatalogs{},
	}.Register()

	// GetAdminCatalog
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminCatalog.html",
		Name:             "GetAdminCatalog",
		Description:      "Get a catalog by ID (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminCatalog,
		PathParams: []cav.PathParam{
			{
				Name:        pathAdminCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
		ResponseType: itypes.AdminCatalog{},
	}.Register()

	// CreateAdminCatalog
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-CreateCatalog.html",
		Name:             "CreateAdminCatalog",
		Description:      "Create a catalog (admin scope)",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/admin/org/{orgId}/catalog",
		PathParams: []cav.PathParam{
			{
				Name:        pathParamOrgIDAdmin,
				Description: descOrgIDAdmin,
				Required:    true,
			},
		},
		BodyRequestType: itypes.AdminCatalogRequest{},
		ResponseType:    itypes.AdminCatalog{},
	}.Register()

	// UpdateAdminCatalog
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-AdminCatalog.html",
		Name:             "UpdateAdminCatalog",
		Description:      "Update a catalog (admin scope)",
		Method:           cav.MethodPUT,
		Backend:          cav.BackendVMware,
		PathTemplate:     pathAdminCatalog,
		PathParams: []cav.PathParam{
			{
				Name:        pathAdminCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.AdminCatalogRequest{},
		ResponseType:    itypes.AdminCatalog{},
	}.Register()

	// DeleteAdminCatalog
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-AdminCatalog.html",
		Name:             "DeleteAdminCatalog",
		Description:      "Delete a catalog (admin scope)",
		Method:           cav.MethodDELETE,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/admin/catalog/{catalogId}",
		PathParams: []cav.PathParam{
			{
				Name:        pathAdminCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
	}.Register()

	// GetAdminCatalogACL
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminCatalogControlAccess.html",
		Name:             "GetAdminCatalogACL",
		Description:      "Get catalog ACL (admin scope)",
		Method:           cav.MethodGET,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/admin/catalog/{catalogId}/controlAccess",
		PathParams: []cav.PathParam{
			{
				Name:        pathAdminCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
		ResponseType: itypes.ControlAccessParams{},
	}.Register()

	// SetAdminCatalogACL
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-AdminCatalogControlAccess.html",
		Name:             "SetAdminCatalogACL",
		Description:      "Set catalog ACL (admin scope)",
		Method:           cav.MethodPOST,
		Backend:          cav.BackendVMware,
		PathTemplate:     "/api/admin/catalog/{catalogId}/action/controlAccess",
		PathParams: []cav.PathParam{
			{
				Name:        pathAdminCatalogID,
				Description: descCatalogID,
				Required:    true,
			},
		},
		BodyRequestType: itypes.ControlAccessParams{},
		ResponseType:    itypes.ControlAccessParams{},
	}.Register()
}
