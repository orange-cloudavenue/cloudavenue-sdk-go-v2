/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// ListAdminOrgs - List organizations (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-OrganizationsFromQuery.html
func ListAdminOrgs() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAdminOrgs")
}

// GetAdminOrg - Get an organization by ID (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-Organization-AdminView.html
func GetAdminOrg() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminOrg")
}

// ListAdminCatalogs - List catalogs (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-CatalogsFromQuery.html
func ListAdminCatalogs() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAdminCatalogs")
}

// GetAdminCatalog - Get a catalog by ID (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-Catalog-AdminView.html
func GetAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminCatalog")
}

// CreateAdminCatalog - Create a catalog (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-SecuredCreateCatalog.html
func CreateAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateAdminCatalog")
}

// UpdateAdminCatalog - Update a catalog (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-Catalog.html
func UpdateAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateAdminCatalog")
}

// DeleteAdminCatalog - Delete a catalog (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/DELETE-Catalog.html
func DeleteAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteAdminCatalog")
}

// GetAdminCatalogACL - Get catalog ACL (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-ControlAccessParams-catalog.html
func GetAdminCatalogACL() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminCatalogACL")
}

// SetAdminCatalogACL - Set catalog ACL (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ControlAccess-catalog.html
func SetAdminCatalogACL() *cav.Endpoint {
	return cav.MustGetEndpoint("SetAdminCatalogACL")
}
