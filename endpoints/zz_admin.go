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
// DocumentationURL: 
func ListAdminOrgs() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAdminOrgs")
}
// GetAdminOrg - Get an organization by ID (admin scope)
//
// DocumentationURL: 
func GetAdminOrg() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminOrg")
}
// ListAdminVDCs - List VDCs (admin scope)
//
// DocumentationURL: 
func ListAdminVDCs() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAdminVDCs")
}
// GetAdminVDC - Get a VDC by ID (admin scope)
//
// DocumentationURL: 
func GetAdminVDC() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminVDC")
}
// ListAdminCatalogs - List catalogs (admin scope)
//
// DocumentationURL: 
func ListAdminCatalogs() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAdminCatalogs")
}
// GetAdminCatalog - Get a catalog by ID (admin scope)
//
// DocumentationURL: 
func GetAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminCatalog")
}
// CreateAdminCatalog - Create a catalog (admin scope)
//
// DocumentationURL: 
func CreateAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateAdminCatalog")
}
// UpdateAdminCatalog - Update a catalog (admin scope)
//
// DocumentationURL: 
func UpdateAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateAdminCatalog")
}
// DeleteAdminCatalog - Delete a catalog (admin scope)
//
// DocumentationURL: 
func DeleteAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteAdminCatalog")
}
// GetAdminCatalogACL - Get catalog ACL (admin scope)
//
// DocumentationURL: 
func GetAdminCatalogACL() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminCatalogACL")
}
// SetAdminCatalogACL - Set catalog ACL (admin scope)
//
// DocumentationURL: 
func SetAdminCatalogACL() *cav.Endpoint {
	return cav.MustGetEndpoint("SetAdminCatalogACL")
}

