// SPDX-FileCopyrightText: Copyright (c) 2025 Orange
// SPDX-License-Identifier: Mozilla Public License 2.0
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.
//

package cav

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// TODO refacto

const (
	mockOrg = "cav01ev01ocb0001234"
)

var pathPrefix = map[subClientName]string{
	ClientVmware:         "",
	ClientCerberus:       "",
	ClientNetbackup:      "/netbackup",
	subClientName("ihm"): "/ihm",
	subClientName("s3"):  "/s3",
}

func newMockClient() (Client, error) {
	// Mock implementation for testing purposes

	// Get All endpoints available in the endpoint package
	// Create an handler for each endpoint
	// Each handler should return a mock response
	// This is a placeholder for the actual implementation

	endpoints := GetEndpointsUncategorized()

	mux := chi.NewRouter()

	for _, ep := range endpoints {
		if ep.MockResponseFuncIsDefined() {
			xlogger.Debug("Registering mock responseFunc for endpoint", slog.String("endpoint", ep.Name), slog.String("method", ep.Method.String()))
			mux.MethodFunc(ep.Method.String(), ep.MockPath(), ep.GetMockResponseFunc(ep))
			continue
		}

		if ep.Method == MethodGET {
			mux.MethodFunc(ep.Method.String(), ep.MockPath(), GetDefaultMockResponseFunc(ep))
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

			if reflectBodyType == reflect.TypeOf(Job{}) {
				statusAccepted := http.StatusAccepted
				ep.SetMockResponseFunc(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Add("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
					w.WriteHeader(statusAccepted)
				}))
			}
		}

		mux.MethodFunc(ep.Method.String(), ep.MockPath(), GetDefaultMockResponseFunc(ep))
	}

	hts := httptest.NewServer(mux)

	xlogger.Debug("Mock server started", slog.String("url", hts.URL))

	nC, err := NewClient(
		mockOrg,
		WithCustomEndpoints(consoles.Services{
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
				Endpoint: hts.URL,
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

func buildPath(subClient subClientName, path string) string {
	if !strings.HasPrefix(path, pathPrefix[subClient]) {
		return pathPrefix[subClient] + path
	}
	return path
}

// * Not used for the moment, but can be used to set mock responses for endpoints.

// func setMockResponse(ep *Endpoint, mockResponseData any, mockResponseStatusCode *int) {
// 	if ep.MockResponseFuncIsDefined() {
// 		return
// 	}

// 	ep.SetMockResponse(mockResponseData, mockResponseStatusCode)
// }

// func cleanMockResponses() {
// 	endpoints := GetEndpointsUncategorized()
// 	for _, ep := range endpoints {
// 		if ep.MockResponseFuncIsDefined() {
// 			ep.CleanMockResponse()
// 		}
// 	}
// }
