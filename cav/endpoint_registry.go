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
	"crypto/sha256"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

const (
	// API is the type for API names. (Type available in endpoint.go)
	// Version is the type for API versions. (Type available in endpoint.go)

	// * Exported api
	// APIVDC         API = "vdc"
	// APIEdgeGateway API = "edgegateway"
	// APIVApp        API = "vapp"
	APIOrg API = "org"

	// * versions
	VersionV1 Version = "v1"
	VersionV2 Version = "v2"

	// * Methods
	MethodGET    Method = "GET"
	MethodPOST   Method = "POST"
	MethodPUT    Method = "PUT"
	MethodDELETE Method = "DELETE"
	MethodPATCH  Method = "PATCH"
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
	Map: make(map[string]*Endpoint),
}

// Register registers an endpoint in the Endpoints map.
func (e Endpoint) Register() {
	// logger is not used in Register method, because it is called
	// by init() function, which is called before the logger is initialized.

	if err := validators.New().Struct(&e); err != nil {
		panic(err)
	}

	pc, _, _, ok := runtime.Caller(1)
	if ok {
		e.api, e.version = decodeCallerPackageName(runtime.FuncForPC(pc).Name())
	}

	if e.api == "" || e.version == "" {
		// If the API and version are not set, we try to
		panic("unable to determine API and version from caller context.")
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

	// Encode the endpoint to create a unique key
	encodedKey := encodeEndpoint(endpoint.api, endpoint.version, endpoint.Name, endpoint.Method)

	// Store the endpoint in the map using the encoded key
	e.Map[encodedKey] = endpoint
}

// encode encodes the API, version, name, and method into a sha256 string.
// This is used to create a unique identifier for the endpoint.
func encodeEndpoint(api API, version Version, name string, method Method) string {
	delimiter := "/"
	s := string(api) + delimiter + string(version) + delimiter + name + delimiter + string(method)

	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
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

// GetEndpoint retrieves an endpoint from the Endpoints map based on the provided api, version, name, and method.
func GetEndpoint(name string, method Method, opts ...EndpointRegistryOptions) (*Endpoint, error) {
	endpoints.mu.RLock()
	defer endpoints.mu.RUnlock()

	extraData := endpointRegistryOptions{}

	for _, opt := range opts {
		opt(&extraData)
	}

	if extraData.api == "" || extraData.version == "" {
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			extraData.api, extraData.version = decodeCallerPackageName(runtime.FuncForPC(pc).Name())
		}
		if extraData.api == "" || extraData.version == "" {
			return nil, errors.New("unable to determine API and version from caller context, use WithExtraProperties() to specify them explicitly")
		}
	}

	encodedKey := encodeEndpoint(extraData.api, extraData.version, name, method)
	if endpoint, ok := endpoints.Map[encodedKey]; ok {
		return endpoint, nil
	}

	return nil, errors.Newf("method %s not found for name %s in version %s of api %s", method, name, extraData.version, extraData.api)
}

func decodeCallerPackageName(funcName string) (API, Version) {
	// funcName == github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1.funcName

	// Remove prefix  (Result: /api/vdc/v1.funcName)
	pkg := strings.TrimPrefix(funcName, "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/")

	// Remove funcName from /api/vdc/v1.funcName (Result: /api/vdc/v1)
	pkg = strings.Split(pkg, ".")[0]

	switch {
	case strings.HasPrefix(pkg, "api/"):
		x := strings.SplitN(pkg, "/", 3)
		if len(x) == 3 {
			return API(x[1]), Version(x[2])
		}
	case strings.HasPrefix(pkg, "cav"):
		// If the package is "cav", we assume it's a core endpoint
		return API("cav"), VersionV1
	}

	return "", ""
}
