// SPDX-FileCopyrightText: Copyright (c) 2025 Orange
// SPDX-License-Identifier: Mozilla Public License 2.0
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.
//

package cav

import (
	"log"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

const (
	mockOrg = "cav01ev01ocb0001234"
)

var (
	pathPrefix = map[SubClientName]string{
		ClientVmware:         "/cloudapi",
		ClientCerberus:       "/api/customers",
		ClientNetbackup:      "/netbackup",
		SubClientName("ihm"): "/ihm",
		SubClientName("s3"):  "/s3",
	}
)

func newMockClient() (Client, error) {
	// Mock implementation for testing purposes

	// Get All endpoints available in the endpoint package
	// Create an handler for each endpoint
	// Each handler should return a mock response
	// This is a placeholder for the actual implementation

	endpoints := GetEndpointsUncategorized()

	mux := chi.NewRouter()

	for _, ep := range endpoints {
		switch ep.Method {
		case MethodGET:
			if ep.MockResponseFuncIsDefined() {
				mux.Get(buildPath(ep.SubClient, ep.PathTemplate), ep.GetMockResponseFunc())
				continue
			}

			log.Default().Println("No mock response defined for endpoint", ep.Name)
			mux.Get(buildPath(ep.SubClient, ep.PathTemplate), GetDefaultMockResponseFunc(ep))
		case MethodPOST:

			if ep.MockResponseFuncIsDefined() {
				mux.Post(buildPath(ep.SubClient, ep.PathTemplate), ep.GetMockResponseFunc())
				continue
			}

			log.Default().Println("No mock response defined for endpoint", ep.Name)
			mux.Post(buildPath(ep.SubClient, ep.PathTemplate), PostDefaultMockResponseFunc(ep))
		}
	}

	hts := httptest.NewServer(mux)

	log.Default().Println("Mock server started at", hts.URL)

	nC, err := NewClient(
		mockOrg,
		WithCustomEndpoints(consoles.Services{
			IHM: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/ihm",
			},
			APIVCD: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/cloudapi",
			},
			APICerberus: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/api/customers",
			},
			S3: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/s3",
			},
			Netbackup: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/netbackup",
			},
		}),
		WithCloudAvenueCredential("mockuser", "mockpassword"),
	)
	if err != nil {
		return nil, err
	}

	return nC, nil
}

func buildPath(subClient SubClientName, path string) string {
	if !strings.HasPrefix(path, pathPrefix[subClient]) {
		return pathPrefix[subClient] + path
	}
	return path
}

func setMockResponse(ep *Endpoint, mockResponseData any, mockResponseStatusCode *int) {
	if ep.MockResponseFuncIsDefined() {
		log.Default().Println("Mock response already defined for endpoint", ep.Name)
		return
	}

	ep.SetMockResponse(mockResponseData, mockResponseStatusCode)
	log.Default().Printf("Mock response set for endpoint %s with status code %d", ep.Name, mockResponseStatusCode)
}

func cleanMockResponses() {
	endpoints := GetEndpointsUncategorized()
	for _, ep := range endpoints {
		if ep.MockResponseFuncIsDefined() {
			ep.CleanMockResponse()
			log.Default().Printf("Mock response cleaned for endpoint %s", ep.Name)
		}
	}
}
