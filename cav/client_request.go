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
	"fmt"
	"time"

	"resty.dev/v3"
)

// NewRequest creates a new request using the resty client.
func (c *client) NewRequest(ctx context.Context, endpoint *Endpoint, _ ...RequestOption) (req *resty.Request, err error) {
	// Retrieve the subclient based on the provided client name.
	// This method identifies the subclient and returns it.
	sc, err := c.identifyClient(ctx, endpoint.SubClient)
	if err != nil {
		return nil, err
	}

	// Inject the client name into the context for retrieval in the other methods.
	ctxv := context.WithValue(ctx, contextKeyClientName, endpoint.SubClient)

	// Create and populate the request options.
	// TODO: actually the request options are not used in the current implementation.
	// They are intended to be used for setting custom options on the request.
	// reqOpts, err := newRequestOptions(opts...)
	// if err != nil {
	// 	return nil, err
	// }

	// TODO(azrod) thinking about removing the error from the new HTTPClient because
	// it is not really useful, we can just return the request with the context. No error
	// should be returned here, because the client should always be able to create a new HTTP.

	// Setup the HTTP client for the request.
	// This client is used to send the request and handle the response.
	hC, err := sc.NewHTTPClient(ctxv)
	if err != nil {
		return nil, err
	}

	// * Middlewares

	// ? Request Middlewares
	if endpoint.RequestMiddlewares != nil {
		// If the endpoint has request middlewares, set them on the HTTP client.
		// This allows for custom processing of the request before it is sent.
		for _, mw := range endpoint.RequestMiddlewares {
			hC.AddRequestMiddleware(mw)
		}
	}

	if isMockClient {
		// If the client is a mock client, we need to override the request URL to point to special prefix.
		// This is because the mock client uses a different URL structure for the mock endpoints.
		// The mock client will handle the request and return a mock response.

		hC.AddRequestMiddleware(resty.RequestMiddleware(func(_ *resty.Client, r *resty.Request) error {
			// Set the base URL to the mock endpoint URL.
			r.URL = fmt.Sprintf("%s%s", hC.BaseURL(), endpoint.MockPath())
			return nil
		}))
	}

	// ? Response Middlewares
	if endpoint.ResponseMiddlewares != nil {
		// If the endpoint has response middlewares, set them on the HTTP client.
		// This allows for custom processing of the response after it is received.
		for _, mw := range endpoint.ResponseMiddlewares {
			hC.AddResponseMiddleware(mw)
		}
	}

	// If JobOpts are provided, we need to create a request with job middleware.
	// This is used to handle job responses and status checks.
	if endpoint.JobOptions != nil {
		sCJob, ok := sc.(jobsInterface)
		if !ok {
			// If the subclient does not implement the jobs.Client interface,
			// we cannot create a job request.
			// Return an error indicating that the client does not support job options.
			// This is a design choice to ensure that only clients that support jobs can use this feature.
			// If you need to handle jobs, ensure that the client implements the jobs.Client interface.
			return nil, fmt.Errorf("client %s does not support job options", endpoint.SubClient)
		}

		// Create a new HTTP client specifically for job requests.
		// This is necessary because the initial client (hc) have a specific middleware defined below.
		// If the hc client has used in NewJobMiddleware, it will create an infinite loop.
		// So we create a new client for job requests.
		hCForJob, err := sc.NewHTTPClient(ctxv)
		if err != nil {
			return nil, err
		}

		// If the request is for a job, set the job middleware.
		// hC.SetResponseMiddlewares(
		// 	resty.AutoParseResponseMiddleware,
		// 	newJobMiddleware(hCForJob, sCJob, endpoint.JobOptions),
		// )
		hC.AddResponseMiddleware(newJobMiddleware(hCForJob, sCJob, endpoint.JobOptions))
	}

	var (
		retryCount           = 5
		retryWaitTime        = 60 * time.Second
		retryMaxWaitTime     = 5 * time.Second
		retryConditionsFuncs = make([]resty.RetryConditionFunc, 0)
		retryIdempotent      = false
	)

	if endpoint.RetryConditionsFuncs != nil {
		retryConditionsFuncs = append(retryConditionsFuncs, endpoint.RetryConditionsFuncs...)
	}

	if isMockClient {
		// If the client is a mock client, set the retry value to shorter values.
		retryCount = 1
		retryWaitTime = 5 * time.Millisecond
	}

	switch endpoint.Method {
	case MethodPOST, MethodPUT, MethodDELETE:
		// For POST, PUT, or PATCH requests, add retry hooks to check if the error return BUSY_ENTITY.
		var conflictRetry resty.RetryConditionFunc = func(resp *resty.Response, _ error) bool {
			if sc.idempotentRetryCondition()(resp, nil) {
				// ! Increment the retry count to allow for retries "unlimited" retry for the busy entity error.
				// This is because the busy entity has undefined max time to resolve.
				resp.Request.RetryCount++
				return true
			}

			return false
		}

		retryConditionsFuncs = append(retryConditionsFuncs, conflictRetry)
		retryIdempotent = true
	}

	// Create a new request with the context and options.
	// To know more about retry see https://resty.dev/docs/retry-mechanism/
	hR := hC.NewRequest().
		SetContext(ctxv).
		EnableRetryDefaultConditions().
		SetRetryCount(retryCount).
		SetRetryMaxWaitTime(retryMaxWaitTime).
		SetRetryWaitTime(retryWaitTime).
		AddRetryConditions(retryConditionsFuncs...).
		SetAllowNonIdempotentRetry(retryIdempotent)

	for _, q := range endpoint.QueryParams {
		if q.Value != "" {
			// If a value is provided for the query parameter, use it directly.
			hR.SetQueryParam(q.Name, q.Value)
		}
	}

	// Set the path parameters in the request.
	// This is done to replace the path parameters in the endpoint path template.
	for _, p := range endpoint.PathParams {
		if p.Value != "" {
			// If a value is provided for the path parameter, use it directly.
			hR.SetPathParam(p.Name, p.Value)
		}
	}

	return hR, nil
}

// Do executes the request and returns the response.
func (c *client) Do(ctx context.Context, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error) {
	resp, err := endpoint.RequestFunc(ctx, c, endpoint, opts...)
	if err != nil {
		return nil, err
	}

	// Retrieve the subclient based on the provided client name.
	// This method identifies the subclient and returns it.
	sc, err := c.identifyClient(ctx, endpoint.SubClient)
	if err != nil {
		return nil, err
	}

	if errAPI := sc.parseAPIError(endpoint.Description, resp); errAPI != nil {
		xlogger.Error("API error occurred", "operation", endpoint.Description, "error", errAPI)
		// If the response is an error, parse the API error and return it.
		return nil, errAPI
	}

	// If the response is successful, return the response.
	return resp, nil
}
