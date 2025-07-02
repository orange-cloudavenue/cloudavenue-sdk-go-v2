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
	"context"
	"log"
	"net/http"
	"reflect"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"

	"sync"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/common-go/validators"
)

type (
	Category string
	Version  string

	Endpoint struct {
		// Category is the category of the endpoint, e.g., "vdc", "edgegateway", "vapp"
		Category Category `validate:"required"`

		// Version is the API version, e.g., "v1", "v2, "v3", etc.
		// It is used to differentiate between different versions of the API.
		Version Version `validate:"required"`

		// Name is the name of the endpoint, e.g., "firewall", "loadBalancer", etc.
		// It is used to group endpoints by their functionality.
		// For example, all endpoints related to firewall operations can be grouped under the "firewall" name.
		Name string `validate:"required"`

		// SubClient is the name of the sub-client that this endpoint belongs to.
		SubClient SubClientName `validate:"required"`

		// Method is the HTTP method used for the endpoint, e.g., "GET", "POST", "PUT", "DELETE".
		Method Method `validate:"required,oneof=GET POST PUT DELETE PATCH"`

		// PathTemplate is the URL path template for the endpoint.
		PathTemplate string `validate:"required"` // e.g., "/v1/edgeGateways/{gatewayId}/firewall/rules"

		// PathParams List of path parameters that can be used in the URL path.
		// These parameters are placeholders in the URL that can be replaced with actual values.
		// For example, in the URL "/v1/edgeGateways/{gatewayId}/firewall/rules",
		// "{gatewayId}" is a path parameter that can be replaced with an actual gateway ID.
		// PathParams are used to construct the final URL for the endpoint.
		PathParams []PathParam `validate:"dive"`

		// QueryParams List of query parameters that can be used in the URL.
		// These parameters are appended to the URL as key-value pairs.
		// For example, in the URL "/v1/edgeGateways/{gatewayId}/firewall/rules?active=true",
		// "active" is a query parameter that can be used to filter results.
		// QueryParams are used to add additional information to the URL for the endpoint.
		QueryParams []QueryParam `validate:"dive"`

		// DocumentationURL is the URL to the documentation for this endpoint.
		DocumentationURL string `validate:"required,url"` // e.g., "https://docs.xx.com/api/v1/xx"

		// RequestFunc is a function that takes a client and options and returns a resty.Response and an error.
		// This function is used to make the actual HTTP request to the endpoint.
		// It allows for customization of the request, such as setting headers, query parameters,
		// and body content.
		RequestFunc func(ctx context.Context, client Client, endpoint *Endpoint, opts ...RequestOption) (*resty.Response, error)

		// BodyType is the golang type of the request body.
		// It is used to validate the body areguments passed to the endpoint.
		// BodyType is optional and can be used to specify the type of the request body
		// for POST, PUT, or PATCH requests.
		BodyType any `validate:"omitempty"`

		// body is the request body that can be used in the request.
		// It is optional and can be used to send data in the request body for POST, PUT, or PATCH requests.
		// Body is only accessible through the RequestOption function
		body any

		// RequestInternalFunc is a function that takes a client and options and returns a resty.Response and an error.
		// This function is used to make the actual HTTP request (from internal package) to the endpoint.
		// It allows for customization of the request, such as setting headers, query parameters,
		// and body content.
		requestInternalFunc func(ctx context.Context, client *resty.Client, endpoint *Endpoint, opts ...RequestOption) (*resty.Response, error)

		// mockResponse is the mock response that can be used for testing purposes.
		// It is optional and can be used to simulate a response from the endpoint without making an actual HTTP request.
		mockResponseFunc func(w http.ResponseWriter, _ *http.Request)

		// mockResponseData is the mock response data that can be used for testing purposes.
		mockResponseData any

		// mockResponseStatusCode int
		// mockResponseStatusCode is the HTTP status code to return for the mock response.
		mockResponseStatusCode *int `validate:"omitempty,oneof=200 201 202 204 400 401 403 404 500"`
	}

	QueryParam struct {
		Name          string `validate:"required"`
		Description   string `validate:"required"`
		Required      bool
		ValidatorFunc func(value string) error

		// value is completed by the RequestOption
		value string
	}

	PathParam struct {
		Name          string `validate:"required"`
		Description   string `validate:"required"`
		Required      bool
		ValidatorFunc func(value string) error

		// value is completed by the RequestOption
		value string
	}

	Method string

	RequestOption func(*Endpoint, *resty.Request) error
)

const (
	// * Categories
	CategoryVDC            Category = "vdc"
	CategoryEdgeGateway    Category = "edgegateway"
	CategoryVApp           Category = "vapp"
	CategoryAuthentication Category = "authentication"

	// * Versions
	VersionV1 Version = "v1"
	VersionV2 Version = "v2"

	// * Methods
	MethodGET    Method = "GET"
	MethodPOST   Method = "POST"
	MethodPUT    Method = "PUT"
	MethodDELETE Method = "DELETE"
	MethodPATCH  Method = "PATCH"
)

var (
	// mu is a mutex to protect the Endpoints map from concurrent access.
	// It ensures that only one goroutine can modify the map at a time.
	// This is important because the Endpoints map is shared across multiple goroutines,
	// and concurrent modifications could lead to race conditions.
	mu = &sync.RWMutex{}
)

// map[category]map[version]map[object]map[method]Endpoint
// map[edgegateway]map[v1]map[firewall]map[GET|POST|PUT|DELETE]Endpoint
var endpoints = map[Category]map[Version]map[string]map[Method]*Endpoint{}

// map[api]map[version]

// Register registers an endpoint in the Endpoints map.
func (e Endpoint) Register() error {
	mu.Lock()
	defer mu.Unlock()

	if err := validators.New().Struct(&e); err != nil {
		// Handle validation errors
		return err
	}

	initEndpoint(e)

	log.Default().Printf("Registered endpoint: %s %s %s", e.Method, e.PathTemplate, e.Name)

	if e.RequestFunc == nil {
		// Default RequestFunc if not provided
		e.RequestFunc = DefaultRequestFunc
	}

	// TODO
	if e.BodyType != nil {
		log.Default().Print("====>", reflect.TypeOf(e.BodyType).PkgPath())
	}

	// Set the endpoint in the Endpoints map
	endpoints[e.Category][e.Version][e.Name][e.Method] = &e

	return nil
}

// GetMockResponse retrieves the mock response for the endpoint.
func (e Endpoint) GetMockResponseFunc() func(w http.ResponseWriter, _ *http.Request) {
	if e.mockResponseFunc != nil {
		return e.mockResponseFunc
	}

	// Default mock response if not provided
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Mock response"}`))
	}
}

// MockResponseIsDefined checks if a mock response is defined for the endpoint.
func (e Endpoint) MockResponseFuncIsDefined() bool {
	return e.mockResponseFunc != nil
}

// SetMockResponse sets the mock response for the endpoint.
func (e *Endpoint) SetMockResponseFunc(mockResponse func(w http.ResponseWriter, _ *http.Request)) {
	if mockResponse == nil {
		log.Default().Println("Mock response is nil, not setting it for endpoint:", e.Name)
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

// GetEndpoints retrieves all endpoints for a given category and version.
func GetEndpoints() map[Category]map[Version]map[string]map[Method]*Endpoint {
	mu.RLock()
	defer mu.RUnlock()

	return endpoints

	// // Create a copy of the endpoints map to avoid concurrent modification issues
	// endpointsCopy := make(map[Category]map[Version]map[string]map[Method]*Endpoint)
	// for category, versions := range endpoints {
	// 	endpointsCopy[category] = make(map[Version]map[string]map[Method]Endpoint)
	// 	for version, objects := range versions {
	// 		endpointsCopy[category][version] = make(map[string]map[Method]Endpoint)
	// 		for name, methods := range objects {
	// 			endpointsCopy[category][version][name] = make(map[Method]Endpoint)
	// 			for method, endpoint := range methods {
	// 				endpointsCopy[category][version][name][method] = endpoint
	// 			}
	// 		}
	// 	}
	// }

	// return endpointsCopy
}

// GetEndpointsUncategorized retrieves all endpoints without categorization.
func GetEndpointsUncategorized() []*Endpoint {
	mu.RLock()
	defer mu.RUnlock()

	var endpointsList []*Endpoint

	// Iterate through the endpoints map and collect all endpoints
	for _, versions := range endpoints {
		for _, objects := range versions {
			for _, methods := range objects {
				for _, endpoint := range methods {
					endpointsList = append(endpointsList, endpoint)
				}
			}
		}
	}

	return endpointsList
}

// GetEndpoint retrieves an endpoint from the Endpoints map based on the provided category, version, name, and method.
func GetEndpoint(category Category, version Version, name string, method Method) (*Endpoint, error) {
	mu.RLock()
	defer mu.RUnlock()

	// Check if the category exists
	if _, ok := endpoints[category]; !ok {
		return nil, errors.Newf("category %s not found", category)
	}

	// Check if the version exists in the category
	if _, ok := endpoints[category][version]; !ok {
		return nil, errors.Newf("version %s not found in category %s", version, category)
	}

	// Check if the name exists in the version
	if _, ok := endpoints[category][version][name]; !ok {
		return nil, errors.Newf("name %s not found in version %s of category %s", name, version, category)
	}

	// Check if the method exists in the name
	if endpoint, ok := endpoints[category][version][name][method]; ok {
		return endpoint, nil
	}

	return nil, errors.Newf("method %s not found for name %s in version %s of category %s", method, name, version, category)
}

// initCategory initializes the category in the Endpoints map if it does not exist.
func initCategory(category Category) {
	if _, ok := endpoints[category]; !ok {
		endpoints[category] = make(map[Version]map[string]map[Method]*Endpoint)
	}
}

// initVersion initializes the version in the category map if it does not exist.
func initVersion(category Category, version Version) {
	if _, ok := endpoints[category][version]; !ok {
		endpoints[category][version] = make(map[string]map[Method]*Endpoint)
	}
}

// initMethod initializes the method in the object map if it does not exist.
func initMethod(category Category, version Version, name string, method Method) {
	if _, ok := endpoints[category][version][name][method]; !ok {
		endpoints[category][version][name][method] = &Endpoint{
			Category:         category,
			Version:          version,
			Name:             name,
			Method:           method,
			PathTemplate:     "",
			PathParams:       []PathParam{},
			QueryParams:      []QueryParam{},
			DocumentationURL: "",
		}
	}
}

// initName initializes the name in the object map if it does not exist.
func initName(category Category, version Version, name string) {
	if _, ok := endpoints[category][version][name]; !ok {
		endpoints[category][version][name] = make(map[Method]*Endpoint)
	}
}

// init Method initializes the method in the object map if it does not exist.
func initEndpoint(endpoint Endpoint) {
	initCategory(endpoint.Category)
	initVersion(endpoint.Category, endpoint.Version)
	initName(endpoint.Category, endpoint.Version, endpoint.Name)
	initMethod(endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method)
}
