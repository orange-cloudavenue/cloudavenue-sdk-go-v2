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
	"net/http"
	"sync"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/common-go/validators"
)

const (
	// * Methods
	MethodGET    Method = http.MethodGet
	MethodPOST   Method = http.MethodPost
	MethodPUT    Method = http.MethodPut
	MethodDELETE Method = http.MethodDelete
	MethodPATCH  Method = http.MethodPatch
)

type (
	endpointsMap struct {
		// mu is a mutex to protect the endpoints map from concurrent access.
		mu sync.RWMutex

		// Map is a nested map structure to hold endpoints.
		// String keys is a sha256 encoded string of the API/Version/Name/Method.
		// Map is capitalized to avoid confusion with the map golang.
		Map map[string]*Endpoint
	}
)

var endpoints = endpointsMap{
	mu:  sync.RWMutex{},
	Map: make(map[string]*Endpoint),
}

// Register registers an endpoint in the Endpoints map.
func (e Endpoint) Register() {
	// logger is not used in Register method, because it is called
	// by init() function, which is called before the logger is initialized.

	if err := validators.New().Struct(&e); err != nil {
		panic(err)
	}
	// Set the endpoint in the Endpoints map
	endpoints.register(&e)
}

// register is a helper function to register an endpoint with the given parameters.
func (e *endpointsMap) register(endpoint *Endpoint) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// RequestFunc is a function that will be used to make the request.
	// If it is not set, we will use the default request function.
	if endpoint.RequestFunc == nil {
		switch endpoint.BodyResponseType.(type) {
		case Job, *Job: // If the endpoint is a job, we use the job middleware.
			endpoint.RequestFunc = defaultRequestFuncWithJob
		default:
			// Default RequestFunc
			endpoint.RequestFunc = defaultRequestFunc
		}
	}

	// Store mockResponseFunc in the internal to restore it later.
	if endpoint.MockResponseFunc != nil {
		endpoint.mockResponseFunc = endpoint.MockResponseFunc
	}

	// Store mockResponseData in the internal to restore it later.
	if endpoint.MockResponseData != nil {
		endpoint.mockResponseData = endpoint.MockResponseData
	}

	if _, ok := e.Map[endpoint.Name]; ok {
		panic(errors.Newf("endpoint %q already registered", endpoint.Name))
	}

	// Store the endpoint in the map using the encoded key
	e.Map[endpoint.Name] = endpoint
}

// GetEndpointsUncategorized retrieves all endpoints without categorization.
func GetEndpointsUncategorized() []*Endpoint {
	endpoints.mu.RLock()
	defer endpoints.mu.RUnlock()

	var endpointsList []*Endpoint

	// Iterate through the endpoints map and collect all endpoints
	for _, endpoint := range endpoints.Map {
		endpointsList = append(endpointsList, endpoint)
	}

	return endpointsList
}

// MustGetEndpoint retrieves an endpoint from the Endpoints map based on the provided api, version, name, and method.
// It panics if the endpoint is not found.
func MustGetEndpoint(name string) *Endpoint {
	endpoint, err := GetEndpoint(name)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// GetEndpoint retrieves an endpoint from the Endpoints map based on the provided api, version, name, and method.
func GetEndpoint(name string) (*Endpoint, error) {
	endpoints.mu.RLock()
	defer endpoints.mu.RUnlock()

	if endpoint, ok := endpoints.Map[name]; ok {
		return endpoint, nil
	}

	return nil, errors.Newf("endpoint %q not found", name)
}
