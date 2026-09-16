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
	APIResponseListDraasOnPremise []string
	APIRequestAddDraasOnPremiseIP string
)

func (r APIResponseListDraasOnPremise) ToModel() *types.ModelListDraasOnPremise {
	return &types.ModelListDraasOnPremise{
		IPs: []string(r),
	}
}
