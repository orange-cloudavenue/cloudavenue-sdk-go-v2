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

// CreateOAuthClient - Register an OAuth client for a tenant
//
// DocumentationURL: https://techdocs.broadcom.com/fr/fr/vmware-cis/cloud-director/vmware-cloud-director/10-5/map-for-vmware-cloud-director-tenant-portal-guide-10-5/managing-users-groups-and-roles-in-vcd-tenant/managing-service-accounts-tenant/managing-service-accounts-tenant.html#GUID-C0CFD1AB-C249-4121-BEA6-05A8984ED4E0-en_GUID-C45B5C4F-1929-4344-AFD9-CEB61958B6BC
func CreateOAuthClient() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateOAuthClient")
}

// CreateOAuthToken - Exchange a tenant OAuth assertion or device code for a token
//
// DocumentationURL: https://techdocs.broadcom.com/fr/fr/vmware-cis/cloud-director/vmware-cloud-director/10-5/map-for-vmware-cloud-director-tenant-portal-guide-10-5/managing-users-groups-and-roles-in-vcd-tenant/managing-service-accounts-tenant/managing-service-accounts-tenant.html#GUID-C8CE5A63-8769-4DAE-8770-9ABE298B62F2-en_GUID-1A4671D8-D454-4F4B-8D27-192930DC2D36
func CreateOAuthToken() *cav.Endpoint {
	return cav.MustGetEndpoint("CreateOAuthToken")
}
