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

// ListT0 - List T0
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func ListT0() *cav.Endpoint {
	return cav.MustGetEndpoint("ListT0")
}
// GetT0 - Get T0
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API
func GetT0() *cav.Endpoint {
	return cav.MustGetEndpoint("GetT0")
}

