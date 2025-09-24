/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
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

//go:generate endpoint-generator -path organization.go -output org

func init() {
	// GetBMS
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/BMS/GetBMSList",
		Name:             "GetBMS",
		Description:      "Get Bare Metal Servers (BMS) details",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/bms/v2.0/servers",
		BodyResponseType: itypes.ApiResponseGetBMS{},
	}.Register()
}
