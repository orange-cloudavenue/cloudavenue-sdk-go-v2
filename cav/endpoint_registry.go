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

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

const (
	// HTTP methods supported by endpoint registry.
	MethodGET    Method = http.MethodGet
	MethodPOST   Method = http.MethodPost
	MethodPUT    Method = http.MethodPut
	MethodDELETE Method = http.MethodDelete
	MethodPATCH  Method = http.MethodPatch
)

type (
	endpointsMap struct {
		// mu protects Map from concurrent access.
		mu sync.RWMutex

		// Map stores endpoints by registered name.
		Map map[string]*Endpoint
	}
)

var endpoints = endpointsMap{
	mu:  sync.RWMutex{},
	Map: make(map[string]*Endpoint),
}

// Register validates and stores e in global endpoint registry.
func (e Endpoint) Register() {
	if err := validators.New().Struct(&e); err != nil {
		panic(err)
	}
	endpoints.register(&e)
}

// register stores endpoint in registry.
func (e *endpointsMap) register(endpoint *Endpoint) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.Map[endpoint.Name]; ok {
		panic(errors.Newf("endpoint %q already registered", endpoint.Name))
	}

	e.Map[endpoint.Name] = endpoint
}

// GetEndpointsUncategorized returns all registered endpoints.
func GetEndpointsUncategorized() []*Endpoint {
	endpoints.mu.RLock()
	defer endpoints.mu.RUnlock()

	var endpointsList []*Endpoint

	for _, endpoint := range endpoints.Map {
		endpointsList = append(endpointsList, endpoint)
	}

	return endpointsList
}

// MustGetEndpoint returns named endpoint or panics.
func MustGetEndpoint(name string) *Endpoint {
	endpoint, err := GetEndpoint(name)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// GetEndpoint returns endpoint registered under name.
func GetEndpoint(name string) (*Endpoint, error) {
	endpoints.mu.RLock()
	defer endpoints.mu.RUnlock()

	if endpoint, ok := endpoints.Map[name]; ok {
		return endpoint, nil
	}

	return nil, errors.Newf("endpoint %q not found", name)
}
