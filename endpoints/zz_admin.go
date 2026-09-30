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
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminOrgs.html
func ListAdminOrgs() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAdminOrgs")
}
// GetAdminOrg - Get an organization by ID (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminOrg.html
func GetAdminOrg() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminOrg")
}
// ListAdminVDCs - List VDCs (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminVdcs.html
func ListAdminVDCs() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAdminVDCs")
}
// GetAdminVDC - Get a VDC by ID (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminVdc.html
func GetAdminVDC() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminVDC")
}
// ListAdminCatalogs - List catalogs (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminCatalogs.html
func ListAdminCatalogs() *cav.Endpoint {
	return cav.MustGetEndpoint("ListAdminCatalogs")
}
// GetAdminCatalog - Get a catalog by ID (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminCatalog.html
func GetAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminCatalog")
}
// CreateAdminCatalog - Create a catalog (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-CreateCatalog.html
func CreateAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateAdminCatalog")
}
// UpdateAdminCatalog - Update a catalog (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/PUT-AdminCatalog.html
func UpdateAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("UpdateAdminCatalog")
}
// DeleteAdminCatalog - Delete a catalog (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/DELETE-AdminCatalog.html
func DeleteAdminCatalog() *cav.Endpoint {
	return cav.MustGetEndpoint("DeleteAdminCatalog")
}
// GetAdminCatalogACL - Get catalog ACL (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/GET-AdminCatalogControlAccess.html
func GetAdminCatalogACL() *cav.Endpoint {
	return cav.MustGetEndpoint("GetAdminCatalogACL")
}
// SetAdminCatalogACL - Set catalog ACL (admin scope)
//
// DocumentationURL: https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/operations/POST-AdminCatalogControlAccess.html
func SetAdminCatalogACL() *cav.Endpoint {
	return cav.MustGetEndpoint("SetAdminCatalogACL")
}

