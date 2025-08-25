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
	"encoding/json"
	"net/http"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate endpoint-generator -path draas.go -output draas

func init() {
	// * ListDraasOnPremiseIP
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/VCDA/getVcdaIPs",
		Name:             "ListDraasOnPremiseIp",
		Description:      "List of on premise IP addresses allowed for this organization's draas offer",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/vcda/ips",
		BodyResponseType: itypes.ApiResponseListDraasOnPremise{},
		MockResponseFunc: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			ips := itypes.ApiResponseListDraasOnPremise{generator.MustGenerate("{ipv4address}"), generator.MustGenerate("{ipv4address}")}
			j, _ := json.Marshal(ips)
			// Return a mock response
			w.Header().Set("Content-Type", "application/json")
			// ignore write body error for mock response
			w.Write(j) //nolint:errcheck
		}),
	}.Register()

	// * AddDraasOnPremiseIp
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/VCDA/postVcdaIPs",
		Name:             "AddDraasOnPremiseIp",
		Description:      "Allow a new on premise IP address for this organization's draas offer",
		Method:           cav.MethodPOST,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/vcda/ips/{ip}",
		PathParams: []cav.PathParam{
			{
				Name:        "ip",
				Description: "The on premise IP address to allow",
				Required:    true,
				ValidatorFunc: func(v string) error {
					return validators.New().Var(v, "ipv4")
				},
			},
		},
	}.Register()

	// * RemoveDraasOnPremiseIp
	cav.Endpoint{
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/VCDA/deleteVcdaIPs",
		Name:             "RemoveDraasOnPremiseIp",
		Description:      "Remove an on premise IP address from this organization's draas offer",
		Method:           cav.MethodDELETE,
		SubClient:        cav.ClientCerberus,
		PathTemplate:     "/api/customers/v2.0/vcda/ips/{ip}",
		PathParams: []cav.PathParam{
			{
				Name:        "ip",
				Description: "The on premise IP address to remove",
				Required:    true,
				ValidatorFunc: func(v string) error {
					return validators.New().Var(v, "ipv4")
				},
			},
		},
	}.Register()
}
