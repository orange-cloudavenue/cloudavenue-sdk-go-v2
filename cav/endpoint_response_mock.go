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
	"log/slog"
	"net/http"

	"github.com/go-faker/faker/v4"
)

var defaultMockResponseFunc = func(ep *Endpoint) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		returnErrFromStatusCodeExpected(w, ep.mockResponseStatusCode)

		var body any

		if ep.mockResponseData != nil {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("Using mock response data for endpoint")
			// If mock response data is defined, use it directly
			body = ep.mockResponseData
		} else if ep.BodyResponseType != nil {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("No mock response data defined, generating mock data")
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("Generating mock response data for endpoint")
			body = ep.BodyResponseType
			if err := faker.FakeData(&body); err != nil {
				xlogger.WithGroup("mock").With("endpoint", ep.Name).Error("Error generating mock data for endpoint:", slog.Any("error", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		bodyEncoded, err := json.Marshal(body)
		if err != nil {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Error("Error encoding body for endpoint:", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")

		// Case used to set custom status code beetween 200 and 299
		// If mockResponseStatusCode is defined, use it, otherwise default to 200
		if ep.mockResponseStatusCode != nil {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).With("statusCode", *ep.mockResponseStatusCode).Debug("Setting mock response status code")
			w.WriteHeader(*ep.mockResponseStatusCode)
		} else {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("No mock response status code defined, using 200 OK")
			w.WriteHeader(http.StatusOK)
		}

		w.Write(bodyEncoded) //nolint:errcheck
	}
}

var (
	GetDefaultMockResponseFunc  = defaultMockResponseFunc
	PostDefaultMockResponseFunc = defaultMockResponseFunc
)

func returnErrFromStatusCodeExpected(w http.ResponseWriter, statusCode *int) {
	if statusCode == nil {
		return
	}

	if *statusCode >= 200 && *statusCode < 300 {
		return
	}

	xlogger.WithGroup("mock").With("statusCode", *statusCode).Debug("Mock response error")
	http.Error(w, http.StatusText(*statusCode), *statusCode)
}

// GetMockResponse retrieves the mock response for the endpoint.
func (e Endpoint) GetMockResponseFunc() func(w http.ResponseWriter, _ *http.Request) {
	if e.mockResponseFunc != nil {
		return e.mockResponseFunc
	}

	// Default mock response if not provided
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Mock response"}`)) //nolint:errcheck
	}
}

// MockResponseIsDefined checks if a mock response is defined for the endpoint.
func (e Endpoint) MockResponseFuncIsDefined() bool {
	return e.mockResponseFunc != nil
}

// SetMockResponse sets the mock response for the endpoint.
func (e *Endpoint) SetMockResponseFunc(mockResponse func(w http.ResponseWriter, _ *http.Request)) {
	if mockResponse == nil {
		xlogger.WithGroup("mock").With("endpoint", e.Name).Warn("Attempted to set nil mock response for endpoint, ignoring")
		return
	}
	e.mockResponseFunc = mockResponse
}

// GetMockResponseData retrieves the mock response data for the endpoint.
func (e Endpoint) GetMockResponse() (data any, statusCode *int) {
	return e.mockResponseData, e.mockResponseStatusCode
}

// SetMockResponse sets the mock response data and status code for the endpoint.
func (e *Endpoint) SetMockResponse(mockResponseData any, mockResponseStatusCode *int) {
	e.mockResponseData = mockResponseData
	e.mockResponseStatusCode = mockResponseStatusCode
}

// CleanMockResponse cleans the mock response for the endpoint.
func (e *Endpoint) CleanMockResponse() {
	e.mockResponseFunc = nil
	e.mockResponseData = nil
	e.mockResponseStatusCode = nil
}
