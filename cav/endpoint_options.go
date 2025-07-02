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

func WithPathParam(pp PathParam, value string) RequestOption {
	return func(endpoint *Endpoint, req *resty.Request) error {
		if endpoint.PathParams == nil {
			return errors.Newf("endpoint %s %s %s %s has no path params", endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method)
		}

		for _, p := range endpoint.PathParams {
			if p.Name == pp.Name {
				if p.Required && value == "" {
					return errors.Newf("path param %s is required for endpoint %s %s %s %s", pp.Name, endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method)
				}
				if p.ValidatorFunc != nil && value != "" {
					if err := p.ValidatorFunc(value); err != nil {
						return errors.Newf("path param %s validation failed for endpoint %s %s %s %s: %v", pp.Name, endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method, err)
					}
				}
			}
		}

		req.SetPathParam(pp.Name, value)
		return nil
	}
}

func WithQueryParam(qp QueryParam, value string) RequestOption {
	return func(endpoint *Endpoint, req *resty.Request) error {
		if endpoint.QueryParams == nil {
			return errors.Newf("endpoint %s %s %s %s has no query params", endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method)
		}

		for _, p := range endpoint.QueryParams {
			if p.Name == qp.Name {
				if p.Required && value == "" {
					return errors.Newf("query param %s is required for endpoint %s %s %s %s", qp.Name, endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method)
				}
				if p.ValidatorFunc != nil && value != "" {
					if err := p.ValidatorFunc(value); err != nil {
						return errors.Newf("query param %s validation failed for endpoint %s %s %s %s: %v", qp.Name, endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method, err)
					}
				}
			}
		}

		req.SetQueryParam(qp.Name, value)
		return nil
	}
}

func OverrideSetResult(rt any) RequestOption {
	return func(endpoint *Endpoint, req *resty.Request) error {
		if rt == nil {
			return errors.Newf("result type cannot be nil for endpoint %s %s %s %s", endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method)
		}
		req.SetResult(rt)
		return nil
	}
}

func SetBody(body any) RequestOption {
	return func(endpoint *Endpoint, req *resty.Request) error {
		if body == nil {
			return errors.Newf("body cannot be nil for endpoint %s %s %s %s", endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method)
		}

		// Reflect BodyRequestType and body to ensure they match
		if endpoint.BodyRequestType != nil {
			reflectBodyType := reflect.TypeOf(endpoint.BodyRequestType)
			reflectBody := reflect.TypeOf(body)
			if reflectBody != reflectBodyType {
				return errors.Newf("body must be of type %s for endpoint %s %s %s %s", reflectBodyType, endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method)
			}
		}

		req.SetBody(body)
		return nil
	}
}

func SetCustomRestyOption(option func(*resty.Request)) RequestOption {
	return func(_ *Endpoint, req *resty.Request) error {
		option(req)
		return nil
	}
}
