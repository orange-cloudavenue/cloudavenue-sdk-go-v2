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
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

const (
	mockOrg = "cav01ev01ocb0001234"
)

var pathPrefix = map[cav.SubClientName]string{
	cav.ClientVmware:         "",
	cav.ClientCerberus:       "",
	cav.ClientNetbackup:      "/netbackup",
	cav.SubClientName("ihm"): "/ihm",
	cav.SubClientName("s3"):  "/s3",
}

var logger = xlog.GetGlobalLogger()

func NewClient(opts ...OptionFunc) (cav.Client, error) {
	// Mock implementation for testing purposes

	// Get All endpoints available in the endpoint package
	// Create an handler for each endpoint
	// Each handler should return a mock response
	// This is a placeholder for the actual implementation

	Options := &Options{}
	for _, opt := range opts {
		if err := opt(Options); err != nil {
			return nil, err
		}
	}

	if Options.logger != nil {
		xlog.SetGlobalLogger(Options.logger)
		logger = Options.logger
	}

	endpoints := cav.GetEndpointsUncategorized()

	mux := chi.NewRouter()

	for _, ep := range endpoints {
		if ep.MockResponseFuncIsDefined() {
			logger.Debug("Registering mock responseFunc for endpoint", slog.String("endpoint", ep.Name), slog.String("method", ep.Method.String()))
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
				switch ep.SubClient {
				case cav.ClientCerberus:
					statusCreated := http.StatusCreated
					ep.SetMockResponseFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(statusCreated)
						w.Write([]byte(`{"jobId":"87ab1934-0146-4fb0-80bc-815fea03214d","message":"Job created successfully"}`))
					})

				case cav.ClientVmware:
					statusAccepted := http.StatusAccepted
					ep.SetMockResponseFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.Header().Add("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
						w.WriteHeader(statusAccepted)
					})
				}
			}

			if ep.MockResponseFuncIsDefined() {
				log.Default().Printf("Registering mock responseFunc for endpoint %s with method %s", ep.Name, ep.Method)
				mux.MethodFunc(ep.Method.String(), buildPath(ep.SubClient, ep.PathTemplate), ep.GetMockResponseFunc())
				continue
			}
		}

		mux.MethodFunc(ep.Method.String(), buildPath(ep.SubClient, ep.PathTemplate), cav.GetDefaultMockResponseFunc(ep))
	}

	hts := httptest.NewServer(mux)
	logger.Debug("Mock server created", slog.String("url", hts.URL))

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
		cav.WithCloudAvenueCredential("mockuser", "mockpassword"),
		cav.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	logger.Debug("Mock client created", slog.String("organization", mockOrg))

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
		logger.Debug("Mock response already defined for endpoint", slog.String("endpoint", ep.Name))
		return
	}

	ep.SetMockResponse(mockResponseData, mockResponseStatusCode)
	logger.Debug("Mock response set for endpoint", slog.String("endpoint", ep.Name), slog.Int("status_code", *mockResponseStatusCode))
}

func CleanMockResponses() {
	endpoints := cav.GetEndpointsUncategorized()
	for _, ep := range endpoints {
		if ep.MockResponseFuncIsDefined() {
			ep.CleanMockResponse()
			logger.Debug("Mock response cleaned for endpoint", slog.String("endpoint", ep.Name))
		}
	}
}

var GetEndpoint = cav.GetEndpoint
