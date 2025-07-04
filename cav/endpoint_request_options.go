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
	"reflect"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

type (
	// EndpointRequestOption is a function that modifies the request for an endpoint.
	// It takes an Endpoint and a resty.Request as parameters and returns an error.
	EndpointRequestOption func(*Endpoint, *resty.Request) error
)

func WithPathParam(pp PathParam, value string) EndpointRequestOption {
	return func(endpoint *Endpoint, req *resty.Request) error {
		if endpoint.PathParams == nil {
			return errors.Newf("endpoint %s %s %s %s has no path params", endpoint.api, endpoint.version, endpoint.Name, endpoint.Method)
		}

		for _, p := range endpoint.PathParams {
			if p.Name == pp.Name {
				if p.Required && value == "" {
					return errors.Newf("path param %s is required for endpoint %s %s %s %s", pp.Name, endpoint.api, endpoint.version, endpoint.Name, endpoint.Method)
				}
				if p.ValidatorFunc != nil && value != "" {
					if err := p.ValidatorFunc(value); err != nil {
						return errors.Newf("path param %s validation failed for endpoint %s %s %s %s: %v", pp.Name, endpoint.api, endpoint.version, endpoint.Name, endpoint.Method, err)
					}
				}
			}
		}

		req.SetPathParam(pp.Name, value)
		return nil
	}
}

func WithQueryParam(qp QueryParam, value string) EndpointRequestOption {
	return func(endpoint *Endpoint, req *resty.Request) error {
		if endpoint.QueryParams == nil {
			return errors.Newf("endpoint %s %s %s %s has no query params", endpoint.api, endpoint.version, endpoint.Name, endpoint.Method)
		}

		for _, p := range endpoint.QueryParams {
			if p.Name == qp.Name {
				if p.Required && value == "" {
					return errors.Newf("query param %s is required for endpoint %s %s %s %s", qp.Name, endpoint.api, endpoint.version, endpoint.Name, endpoint.Method)
				}
				if p.ValidatorFunc != nil && value != "" {
					if err := p.ValidatorFunc(value); err != nil {
						return errors.Newf("query param %s validation failed for endpoint %s %s %s %s: %v", qp.Name, endpoint.api, endpoint.version, endpoint.Name, endpoint.Method, err)
					}
				}
			}
		}

		req.SetQueryParam(qp.Name, value)
		return nil
	}
}

func OverrideSetResult(rt any) EndpointRequestOption {
	return func(endpoint *Endpoint, req *resty.Request) error {
		if rt == nil {
			return errors.Newf("result type cannot be nil for endpoint %s %s %s %s", endpoint.api, endpoint.version, endpoint.Name, endpoint.Method)
		}
		req.SetResult(rt)
		return nil
	}
}

func SetBody(body any) EndpointRequestOption {
	return func(endpoint *Endpoint, req *resty.Request) error {
		if body == nil {
			return errors.Newf("body cannot be nil for endpoint %s %s %s %s", endpoint.api, endpoint.version, endpoint.Name, endpoint.Method)
		}

		// Reflect BodyRequestType and body to ensure they match
		if endpoint.BodyRequestType != nil {
			reflectBodyType := reflect.TypeOf(endpoint.BodyRequestType)
			if reflectBodyType.Kind() == reflect.Ptr {
				reflectBodyType = reflectBodyType.Elem()
			}

			reflectBody := reflect.TypeOf(body)
			if reflectBody.Kind() == reflect.Ptr {
				reflectBody = reflectBody.Elem()
			}
			if reflectBody != reflectBodyType {
				return errors.Newf("body must be of type %s (not %s) for endpoint %s %s %s %s", reflectBodyType, reflectBody, endpoint.api, endpoint.version, endpoint.Name, endpoint.Method)
			}
		}

		req.SetBody(body)
		return nil
	}
}

func SetCustomRestyOption(option func(*resty.Request)) EndpointRequestOption {
	return func(_ *Endpoint, req *resty.Request) error {
		option(req)
		return nil
	}
}
