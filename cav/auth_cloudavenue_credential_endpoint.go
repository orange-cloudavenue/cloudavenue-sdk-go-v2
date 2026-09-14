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
	"encoding/json"
	"net/http"

	"github.com/orange-cloudavenue/common-go/generator"
)

func init() {
	Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v39.1/cloudapi/1.0.0/sessions/post/",
		Name:             "SessionVmware",
		Description:      "Create a new VMware session",
		Method:           MethodPOST,
		Backend:          BackendVMware,
		PathTemplate:     "/cloudapi/1.0.0/sessions",
		PathParams:       []PathParam{},
		QueryParams:      []QueryParam{},
		BodyRequestType:  nil, // No request body for this endpoint.
		ResponseType:     apiResponseSessionVmware{},
	}.Register()
}

// NewSessionVmwareMockHandler returns mock VMware session handler.
func NewSessionVmwareMockHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add(cloudavenueCredentialXVmwareAccessToken, "mock-access-token")

		resp := apiResponseSessionVmware{}

		generator.MustStruct(&resp)

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
	}
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
