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
		// api is the api of the endpoint, e.g., "vdc", "edgegateway", "vapp"
		api API `validate:"required"`

		// version is the API version, e.g., "v1", "v2, "v3", etc.
		// It is used to differentiate between different versions of the API.
		version Version `validate:"required"`

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

		// BodyRequestType is the golang type of the request body.
		// It is used to validate the body arguments passed to the endpoint.
		// BodyType is optional and can be used to specify the type of the request body
		// for POST, PUT, or PATCH requests.
		BodyRequestType any `validate:"required_if=Method POST PUT PATCH"`

		// BodyResponseType is the golang type of the response body.
		// It is used to validate the response body returned by the endpoint.
		//
		// If your set `cav.Job{}` as BodyResponseType, the system will automatically
		// handle the job response and retrieve the job status until it is completed (success or error).
		BodyResponseType any `validate:"omitempty"`

		// * Request

		// RequestFunc is a function that takes a client and options and returns a resty.Response and an error.
		// This function is used to make the actual HTTP request to the endpoint.
		// It allows for customization of the request, such as setting headers, query parameters,
		// and body content.
		RequestFunc func(ctx context.Context, client Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error)

		// RetryHooksFuncs is a list of functions that can be used to customize the retry behavior of the request.
		// These functions can be used to modify the request, such as adding headers or query parameters,
		// or to handle specific conditions that require a retry.
		// To know more about retry see https://resty.dev/docs/retry-mechanism/
		RetryHooksFuncs []resty.RetryHookFunc

		// RequestMiddleware is a function that takes a resty.Request and returns a resty.Request.
		// This function is used to modify the request before it is sent.
		// It allows for adding headers, query parameters, and other modifications to the request.
		RequestMiddlewares []resty.RequestMiddleware

		// RequestInternalFunc is a function that takes a client and options and returns a resty.Response and an error.
		// This function is used to make the actual HTTP request (from internal package) to the endpoint.
		// It allows for customization of the request, such as setting headers, query parameters,
		// and body content.
		requestInternalFunc func(ctx context.Context, client *resty.Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error)

		// * Mock

		// mockResponse is the mock response that can be used for testing purposes.
		// It is optional and can be used to simulate a response from the endpoint without making an actual HTTP request.
		mockResponseFunc func(w http.ResponseWriter, _ *http.Request)

		// mockResponseData is the mock response data that can be used for testing purposes.
		mockResponseData any

		// mockResponseStatusCode int
		// mockResponseStatusCode is the HTTP status code to return for the mock response.
		mockResponseStatusCode *int `validate:"omitempty"`

		// * Job

		// jobOptions is the options for the job.
		// It is used to specify the options for the job, such as the Timeout, PollingInterval, and ExtractorFunc
		JobOptions *JobOptions
	}

	QueryParam struct {
		Name          string `validate:"required"`
		Description   string `validate:"required"`
		Required      bool
		ValidatorFunc func(value string) error
	}

	PathParam struct {
		Name          string `validate:"required"`
		Description   string `validate:"required"`
		Required      bool
		ValidatorFunc func(value string) error
	}
)
