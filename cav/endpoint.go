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
	"net/http"

	"resty.dev/v3"
)

type (
	API     string
	Version string
	Method  string

	Endpoint struct {
		// ID is the unique identifier of the endpoint.
		// It is used to uniquely identify the endpoint in the registry.
		ID string

		// // api is the api of the endpoint, e.g., "vdc", "edgegateway", "vapp"
		// api API `validate:"required"`

		// // version is the API version, e.g., "v1", "v2, "v3", etc.
		// // It is used to differentiate between different versions of the API.
		// version Version `validate:"required"`

		// Name is the name of the endpoint, e.g., "firewall", "loadBalancer", etc.
		// It is used to group endpoints by their functionality.
		// For example, all endpoints related to firewall operations can be grouped under the "firewall" name.
		Name string `validate:"required,disallow_space,case=PascalCase"` // e.g., "Firewall", "LoadBalancer"

		// Description is a brief description of the endpoint.
		// It provides additional information about the endpoint's purpose and functionality.
		// Description is used to provide context in the error messages.
		Description string `validate:"required"`

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

		// * Request

		// BodyRequestType is the golang type of the request body.
		// It is used to validate the body arguments passed to the endpoint.
		// BodyType is optional and can be used to specify the type of the request body
		// for POST, PUT, or PATCH requests.
		BodyRequestType any `validate:"-"`

		// RequestFunc is a function that takes a client and options and returns a resty.Response and an error.
		// This function is used to make the actual HTTP request to the endpoint.
		// It allows for customization of the request, such as setting headers, query parameters,
		// and body content.
		RequestFunc func(ctx context.Context, client Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error)

		// RetryConditions is a list of functions that can be used to customize the retry behavior of the request.
		// These functions can be used to modify the request, such as adding headers or query parameters,
		// or to handle specific conditions that require a retry.
		// To know more about retry see https://resty.dev/docs/retry-mechanism/
		RetryConditionsFuncs []resty.RetryConditionFunc

		// RequestMiddleware is a function that takes a resty.Request and returns a resty.Request.
		// This function is used to modify the request before it is sent.
		// It allows for adding headers, query parameters, and other modifications to the request.
		RequestMiddlewares []resty.RequestMiddleware

		// RequestInternalFunc is a function that takes a client and options and returns a resty.Response and an error.
		// This function is used to make the actual HTTP request (from internal package) to the endpoint.
		// It allows for customization of the request, such as setting headers, query parameters,
		// and body content.
		requestInternalFunc func(ctx context.Context, client *resty.Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error)

		// * Response

		// BodyResponseType is the golang type of the response body.
		// It is used to validate the response body returned by the endpoint.
		//
		// If your set `cav.Job{}` as BodyResponseType, the system will automatically
		// handle the job response and retrieve the job status until it is completed (success or error).
		BodyResponseType any `validate:"-"`

		// ResponseMiddleware is a function that takes a resty.Response and returns a resty.Response.
		// This function is used to modify the response after it is received.
		// It allows for handling the response, such as parsing the body, checking headers,
		// and handling errors.
		ResponseMiddlewares []resty.ResponseMiddleware

		// * Mock

		// mockResponse is the mock response that can be used for testing purposes.
		// It is optional and can be used to simulate a response from the endpoint without making an actual HTTP request.
		MockResponseFunc http.HandlerFunc

		// internalMockResponseFunc is used to store original mock response function.
		// It is used to restore the original mock response function after it has been overridden.
		mockResponseFunc http.HandlerFunc

		// MockResponseData is the mock response data that can be used for testing purposes.
		MockResponseData any

		// mockResponseData is used to store the original mock response data.
		// It is used to restore the original mock response data after it has been overridden.
		mockResponseData any

		// mockResponseStatusCode int
		// mockResponseStatusCode is the HTTP status code to return for the mock response.
		mockResponseStatusCode *int `validate:"omitempty"`

		// * Job

		// jobOptions is the options for the job.
		// It is used to specify the options for the job, such as the Timeout, PollingInterval, and ExtractorFunc
		JobOptions *JobOptions
	}

	// QueryParam represents a query parameter in the URL.
	// Query parameters are appended to the URL as key-value pairs.
	// For example, in the URL "/v1/edgeGateways/{gatewayId}/firewall/rules?active=true",
	// "active" is a query parameter that can be used to filter results.
	// Query parameters are used to add additional information to the URL for the endpoint.
	QueryParam struct {
		Name        string `validate:"required,disallow_space"`
		Description string `validate:"required"`

		// Required indicates whether the query parameter is required or not.
		Required bool

		// ValidatorFunc is a function that validates the value of the query parameter.
		// It is used to ensure that the value of the query parameter is valid before making the request.
		ValidatorFunc func(value string) error

		// TransformFunc is a function that transforms the value of the query parameter.
		// It is used to modify the value of the query parameter before making the request.
		// TransformFunc is called after the ValidatorFunc(if provided) and before the value is set in the request.
		TransformFunc func(value string) (string, error)

		// Ability to provides a value for the query parameter.
		// This is useful when the query parameter value is known at the time of registration.
		// If the value is provided Required, ValidatorFunc and TransformFunc are ignored.
		Value string
	}

	// PathParam represents a path parameter in the URL path.
	// Path parameters are placeholders in the URL that can be replaced with actual values.
	// For example, in the URL "/v1/edgeGateways/{gatewayId}/firewall/rules",
	// "{gatewayId}" is a path parameter that can be replaced with an actual gateway ID.
	// Path parameters are used to construct the final URL for the endpoint.
	PathParam struct {
		Name        string `validate:"required,disallow_space"`
		Description string `validate:"required"`

		// Required indicates whether the path parameter is required or not.
		Required bool

		// ValidatorFunc is a function that validates the value of the path parameter.
		// It is used to ensure that the value of the path parameter is valid before making the request.
		ValidatorFunc func(value string) error

		// TransformFunc is a function that transforms the value of the query parameter.
		// It is used to modify the value of the query parameter before making the request.
		// TransformFunc is called after the ValidatorFunc(if provided) and before the value is set in the request.
		TransformFunc func(value string) (string, error)

		// Ability to provides a value for the path parameter.
		// This is useful when the path parameter value is known at the time of registration.
		// For example, if the path parameter is {type}, you can provide a value like "firewall" or "loadBalancer".
		// If the value is provided Required, ValidatorFunc and TransformFunc are ignored.
		Value string
	}
)
