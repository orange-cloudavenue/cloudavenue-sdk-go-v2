/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
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
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Network%20%26%20connectivity/getNetworkHierarchy
func ListT0() *cav.Endpoint {
	return cav.MustGetEndpoint("ListT0")
}
