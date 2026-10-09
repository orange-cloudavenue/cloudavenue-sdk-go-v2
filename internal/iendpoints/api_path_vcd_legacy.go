/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

const (
	// pathQueryAPI is the legacy VMware query endpoint.
	pathQueryAPI = "/api/query"
	// pathAdminCatalogs lists catalogs through the legacy admin XML API.
	pathAdminCatalogs = "/api/admin/catalogs"
	// pathAdminCatalog addresses one catalog through the legacy admin XML API.
	pathAdminCatalog = "/api/admin/catalog/{catalogId}"
	// pathAdminOrgCatalogs creates catalogs under an organization in the legacy admin XML API.
	pathAdminOrgCatalogs = "/api/admin/org/{orgId}/catalog"
	// pathAdminCatalogControlAccess gets catalog access control through the legacy XML API.
	pathAdminCatalogControlAccess = "/api/catalog/{catalogId}/controlAccess"
	// pathAdminCatalogControlAccessAction sets catalog access control through the legacy XML API.
	pathAdminCatalogControlAccessAction = "/api/catalog/{catalogId}/action/controlAccess"
	// pathAdminOrg addresses one organization in the legacy admin XML API.
	pathAdminOrg = "/api/admin/org/{orgId}"
	// pathAdminOrgs lists organizations through the legacy admin XML API.
	pathAdminOrgs = "/api/admin/orgs"
	// These legacy AdminOrg XML user paths are retained for compatibility only;
	// IAM user operations use VMware CloudAPI.
	// pathAdminOrgUsers lists and creates users in the legacy admin XML API.
	// pathAdminOrgUserByID addresses one user in the legacy admin XML API.
	// pathAdminTakeOwnershipOrgUserByID takes ownership of one user in the legacy admin XML API.
	// pathAdminOrgUserEnable enables one user in the legacy admin XML API.
	// pathAdminOrgUserDisable disables one user in the legacy admin XML API.
	// pathAdminOrgUserUnlock unlocks one user in the legacy admin XML API.
	// pathAdminOrgUserChangePassword changes one user password in the legacy admin XML API.
	// pathVAppByID addresses one vApp in the legacy XML API.
	pathVAppByID = "/api/vapp/{vapp-id}"
	// pathVAppRemoveAllNetworks removes all networks from one vApp.
	pathVAppRemoveAllNetworks = "/api/vapp/{vapp-id}/action/removeAllNetworks"
	// pathVAppUndeploy undeploys one vApp.
	pathVAppUndeploy = "/api/vapp/{vapp-id}/action/undeploy"
	// pathVAppCreate creates a vApp under a VDC in the legacy XML API.
	pathVAppCreate = "/api/vdc/{vdc-id}/action/createVApp"
	// pathVDCMetadata gets VDC metadata through the legacy XML API.
	pathVDCMetadata = "/api/vdc/{vdc-id}/metadata"
	// pathVAppLeaseSettings gets or updates organization vApp lease settings in the legacy XML API.
	pathVAppLeaseSettings = "/api/org/{orgId}/vAppLeaseSettings"
)

const (
	pathParamVDCID              = "vdc-id"
	pathParamVAppID             = "vapp-id"
	pathParamOrgID              = "orgId"
	pathParamUserID             = "userId"
	pathParamID                 = "id"
	pathParamEdgeGatewayIDAdmin = "edgeGatewayId"
	pathParamCatalogID          = "catalogId"
)
