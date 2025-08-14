/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"context"
	"encoding/json"
	"net/http"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/common-go/generator"
)

func init() {
	Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/sessions/post/",
		Name:             "SessionVmware",
		Description:      "Create a new VMware session",
		Method:           MethodPOST,
		SubClient:        ClientVmware,
		PathTemplate:     "/cloudapi/1.0.0/sessions",
		PathParams:       []PathParam{},
		QueryParams:      []QueryParam{},
		RequestFunc:      nil,
		requestInternalFunc: func(ctx context.Context, client *resty.Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error) {
			r := client.R().
				SetContext(ctx).
				SetHeader("Accept", "application/json;version="+vmwareVCDVersion).
				SetResult(&apiResponseSessionVmware{})

			for _, opt := range opts {
				if err := opt(endpoint, r); err != nil {
					return nil, err
				}
			}

			if isMockClient {
				// If the client is a mock client, we return a mock response.
				return r.Post(endpoint.MockPath())
			}

			return r.Post(endpoint.PathTemplate)
		},
		MockResponseFunc: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Add(cloudavenueCredentialXVmwareAccessToken, "mock-access-token")

			resp := apiResponseSessionVmware{}

			generator.MustStruct(&resp)

			// json encode
			w.Header().Set("Content-Type", "application/json")
			respJ, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = w.Write(respJ)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}),
	}.Register()
}

type (
	apiResponseSessionVmware struct {
		Org   apiResponseSessionVmwareOrg  `json:"org"`
		Site  apiResponseSessionVmwareSite `json:"site"`
		Roles []string                     `json:"roles" fake:"Organization Administrator"`
	}

	apiResponseSessionVmwareOrg struct {
		ID   string `json:"id" fake:"{urn:org}"`
		Name string `json:"name" fake:"cav01ev01ocb0001234"`
	}

	apiResponseSessionVmwareSite struct {
		ID   string `json:"id" fake:"{urn:site}"`
		Name string `json:"name" fake:"cav01ev01ocb0001234"`
	}
)
