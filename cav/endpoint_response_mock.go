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
	"reflect"

	"github.com/orange-cloudavenue/common-go/generator"
)

// defaultMockResponseFunc is the default mock response function for endpoints.
// Is used to generate response bodies for mock endpoints.
var defaultMockResponseFunc = func(ep *Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ep.restoreMockResponse() // Restore the original mock response after handling the request

		if ep.MockResponseFuncIsDefined() {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("Using custom mock response function for endpoint")
			ep.GetMockResponseFunc(ep)(w, r)
			return
		}

		// Here catch status code defined in mockResponseStatusCode if >= 300
		if ep.mockResponseStatusCode != nil {
			if *ep.mockResponseStatusCode >= 300 {
				xlogger.WithGroup("mock").With("endpoint", ep.Name).With("statusCode", *ep.mockResponseStatusCode).Debug("Mock response error with status code >= 300")
				var apiError any
				switch ep.SubClient {
				case ClientCerberus:
					apiError = &cerberusError{}
				case ClientVmware:
					apiError = &vmwareError{}
				}
				if err := generator.Struct(apiError); err != nil {
					xlogger.WithGroup("mock").With("endpoint", ep.Name).Error("Error generating mock data for endpoint:", slog.Any("error", err))
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// TODO apiError not used
				xlogger.WithGroup("mock").With("statusCode", *ep.mockResponseStatusCode).Debug("Mock response error")
				http.Error(w, http.StatusText(*ep.mockResponseStatusCode), *ep.mockResponseStatusCode)
				return
			}
		}

		// Construct the mock response body
		var newBody any

		if ep.MockResponseData != nil {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("Using mock response data for endpoint")
			// If mock response data is defined, use it directly
			newBody = ep.MockResponseData
		} else if ep.BodyResponseType != nil {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("No mock response data defined, generating mock data")

			// Reflect on the BodyResponseType to determine the type of the response body
			bodyType := reflect.TypeOf(ep.BodyResponseType)
			// If bodyType is a pointer, we need to dereference it to get the underlying type
			if bodyType.Kind() == reflect.Ptr {
				// Dereference the pointer
				bodyType = bodyType.Elem()
			}

			// Parse special case for bodyType is a Job type
			switch {
			case bodyType == reflect.TypeOf(Job{}):
				switch ep.SubClient {
				case ClientCerberus:
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(`{"jobId":"87ab1934-0146-4fb0-80bc-815fea03214d","message":"Job created successfully"}`)) //nolint:errcheck
					return

				case ClientVmware:
					w.Header().Add("Location", "/mock/cav/v1/jobvmware/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
					w.WriteHeader(http.StatusAccepted)
					return
				}

			default:
				// Set bodyType to a pointer to the struct type
				newBodyType := reflect.PointerTo(bodyType)
				// set new var body with the type of bodyType
				newBody = reflect.New(newBodyType).Interface()
				switch bodyType.Kind() {
				case reflect.Slice:
					xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("BodyResponseType is a slice, generating mock data for slice of structs", slog.String("type", newBodyType.String()))
					// If bodyType is a slice, we need to generate a slice of structs
					// We use the generator to generate a slice of structs

					// Add recovery to handle any panic during generation
					defer func() {
						if r := recover(); r != nil {
							xlogger.WithGroup("mock").With("endpoint", ep.Name).Error("Panic during mock data generation:", slog.Any("error", r))
							http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						}
					}()
					generator.Slice(newBody)

				default:
					xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("BodyResponseType is a struct, generating mock data for struct", slog.String("type", newBodyType.String()))
					if err := generator.Struct(newBody); err != nil {
						xlogger.WithGroup("mock").With("endpoint", ep.Name).Error("Error generating mock data for endpoint:", slog.Any("error", err))
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				}
			}
		}

		xlogger.WithGroup("mock").With("endpoint", ep.Name).Debug("Generated mock response body", slog.Any("body", newBody))

		bodyEncoded, err := json.Marshal(newBody)
		if err != nil {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Error("Error encoding body for endpoint:", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		}

		w.Header().Set("X-Cloud-Avenue-Mock", "true") // Indicate that this is a mock response. For what ? Because !
		_, err = w.Write(bodyEncoded)
		if err != nil {
			xlogger.WithGroup("mock").With("endpoint", ep.Name).Error("Error writing response body for endpoint:", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

var (
	GetDefaultMockResponseFunc  = defaultMockResponseFunc
	PostDefaultMockResponseFunc = defaultMockResponseFunc
)

// GetMockResponse retrieves the mock response for the endpoint.
func (e Endpoint) GetMockResponseFunc(ep *Endpoint) http.HandlerFunc {
	if e.MockResponseFuncIsDefined() {
		return e.MockResponseFunc
	}

	// Default mock response if not provided
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Mock response"}`)) //nolint:errcheck
	})
}

// MockResponseIsDefined checks if a mock response is defined for the endpoint.
func (e Endpoint) MockResponseFuncIsDefined() bool {
	return e.MockResponseFunc != nil
}

// SetMockResponse sets the mock response for the endpoint.
func (e *Endpoint) SetMockResponseFunc(mockResponse http.HandlerFunc) {
	if mockResponse == nil {
		xlogger.WithGroup("mock").With("endpoint", e.Name).Warn("Attempted to set nil mock response for endpoint, ignoring")
		return
	}
	e.MockResponseFunc = mockResponse
}

// GetMockResponseData retrieves the mock response data for the endpoint.
func (e Endpoint) GetMockResponse() (data any, statusCode *int) {
	return e.MockResponseData, e.mockResponseStatusCode
}

// SetMockResponse sets the mock response data and status code for the endpoint.
func (e *Endpoint) SetMockResponse(mockResponseData any, mockResponseStatusCode *int) {
	// Remove MockResponseFunc
	e.MockResponseFunc = nil // Clear the mock response function to use the default one
	// Set the mock response data and status code
	e.MockResponseData = mockResponseData
	e.mockResponseStatusCode = mockResponseStatusCode
}

// CleanMockResponse cleans the mock response for the endpoint.
func (e *Endpoint) CleanMockResponse() {
	xlogger.WithGroup("mock").With("endpoint", e.Name).Debug("Cleaning mock response for endpoint")
	e.MockResponseFunc = nil
	e.MockResponseData = nil
	e.mockResponseStatusCode = nil
}

func (e *Endpoint) RestoreMockResponse() {
	xlogger.WithGroup("mock").With("endpoint", e.Name).Debug("Restoring original mock response for endpoint")
	e.restoreMockResponse()
}

// restoreMockResponse restores the original mock response function and data.
func (e *Endpoint) restoreMockResponse() {
	e.MockResponseFunc = e.mockResponseFunc
	e.MockResponseData = e.mockResponseData
	e.mockResponseStatusCode = nil
}
