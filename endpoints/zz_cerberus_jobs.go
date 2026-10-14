/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is licensed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package endpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// GetJobs - Get current root job id of an organization
//
// DocumentationURL: https://swagger.cloudavenue.orange-business.com/#/Jobs/getJobs
func GetJobs() *cav.Endpoint {
	return cav.MustGetEndpoint("GetJobs")
}
