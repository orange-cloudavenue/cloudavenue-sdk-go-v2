/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

type (
	ApiResponseListDraasOnPremise []string
	ApiRequestAddDraasOnPremiseIP string
)

func (r ApiResponseListDraasOnPremise) ToModel() *types.ModelListDraasOnPremise {
	return &types.ModelListDraasOnPremise{
		IPs: []string(r),
	}
}
