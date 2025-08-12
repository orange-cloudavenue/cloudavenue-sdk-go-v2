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
	"log/slog"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

const (
	mockOrg = "cav01ev01ocb0001234"
)

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

	// Here, for each endpoint, we build a response handler for the mock HTTP server
	for _, ep := range endpoints {
		logger.Debug("Registering mock endpoint", slog.String("name", ep.Name), slog.String("method", ep.Method.String()), slog.String("path", ep.MockPath()), slog.String("ID", ep.ID))
		mux.MethodFunc(ep.Method.String(), ep.MockPath(), cav.GetDefaultMockResponseFunc(ep))
	}

	hts := httptest.NewServer(mux)
	slog.SetDefault(logger)
	hts.Config.ErrorLog = log.Default()

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

func SetMockResponse(ep *cav.Endpoint, mockResponseData any, mockResponseStatusCode *int) {
	ep.SetMockResponse(mockResponseData, mockResponseStatusCode)
	logger.Debug("Mock response set for endpoint", slog.String("endpoint", ep.Name), slog.Int("status_code", *mockResponseStatusCode))
}

var GetEndpoint = cav.GetEndpoint
