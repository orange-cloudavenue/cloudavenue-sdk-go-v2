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
