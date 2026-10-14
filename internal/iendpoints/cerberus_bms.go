/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

//go:generate endpoint-generator -path cerberus_bms.go -output cerberus_bms

func init() {
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/BMS/GetBMSList",
		Name:             "GetBMSList",
		Description:      "Get list of BMS and their description",
		Method:           cav.MethodGET,
		Backend:          cav.BackendInfrapi,
		PathTemplate:     pathCerberusBMS,
		ResponseType:     itypes.APIResponseBMS{},
	}.Register()
}
