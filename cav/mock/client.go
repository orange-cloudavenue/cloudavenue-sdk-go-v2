/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package mock

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

const (
	mockOrg = "cav01ev01ocb0001234"
)

var pathPrefix = map[cav.SubClientName]string{
	cav.ClientVmware:         "",
	cav.ClientCerberus:       "/api/customers",
	cav.ClientNetbackup:      "/netbackup",
	cav.SubClientName("ihm"): "/ihm",
	cav.SubClientName("s3"):  "/s3",
}

func NewClient() (cav.Client, error) {
	// Mock implementation for testing purposes

	// Get All endpoints available in the endpoint package
	// Create an handler for each endpoint
	// Each handler should return a mock response
	// This is a placeholder for the actual implementation

	endpoints := cav.GetEndpointsUncategorized()

	mux := chi.NewRouter()

	for _, ep := range endpoints {
		if ep.MockResponseFuncIsDefined() {
			log.Default().Printf("Registering mock responseFunc for endpoint %s with method %s", ep.Name, ep.Method)
			mux.MethodFunc(ep.Method.String(), buildPath(ep.SubClient, ep.PathTemplate), ep.GetMockResponseFunc())
			continue
		}

		if ep.Method == cav.MethodGET {
			mux.MethodFunc(ep.Method.String(), buildPath(ep.SubClient, ep.PathTemplate), cav.GetDefaultMockResponseFunc(ep))
			continue
		}

		// Methods POST/PUT/PATCH/DELETE require a body
		if ep.BodyResponseType != nil {
			// If the request body type is defined, we need to check if it is a pointer
			// and dereference it to get the actual type
			reflectBodyType := reflect.TypeOf(ep.BodyResponseType)
			if reflectBodyType.Kind() == reflect.Ptr {
				// If the request body type is a pointer, we need to dereference it
				reflectBodyType = reflectBodyType.Elem()
			}

			if reflectBodyType == reflect.TypeOf(cav.Job{}) {
				statusAccepted := http.StatusAccepted
				ep.SetMockResponseFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Add("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
					w.WriteHeader(statusAccepted)
				})
			}
		}

		mux.MethodFunc(ep.Method.String(), buildPath(ep.SubClient, ep.PathTemplate), cav.GetDefaultMockResponseFunc(ep))
	}

	hts := httptest.NewServer(mux)

	log.Default().Println("Mock server started at", hts.URL)

	nC, err := cav.NewClient(
		mockOrg,
		cav.WithCustomEndpoints(consoles.Services{
			IHM: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/ihm",
			},
			APIVCD: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL,
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
		cav.WithCloudAvenueCredential("mockuser", "mockpassword"),
	)
	if err != nil {
		return nil, err
	}

	return nC, nil
}

func buildPath(subClient cav.SubClientName, path string) string {
	if !strings.HasPrefix(path, pathPrefix[subClient]) {
		return pathPrefix[subClient] + path
	}
	return path
}

func SetMockResponse(ep *cav.Endpoint, mockResponseData any, mockResponseStatusCode *int) {
	if ep.MockResponseFuncIsDefined() {
		log.Default().Println("Mock response already defined for endpoint", ep.Name)
		return
	}

	ep.SetMockResponse(mockResponseData, mockResponseStatusCode)
	log.Default().Printf("Mock response set for endpoint %s with status code %d", ep.Name, mockResponseStatusCode)
}

func CleanMockResponses() {
	endpoints := cav.GetEndpointsUncategorized()
	for _, ep := range endpoints {
		if ep.MockResponseFuncIsDefined() {
			ep.CleanMockResponse()
			log.Default().Printf("Mock response cleaned for endpoint %s", ep.Name)
		}
	}
}

var GetEndpoint = cav.GetEndpoint
