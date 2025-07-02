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

		pp.value = value
		endpoint.PathParams = append(endpoint.PathParams, pp)
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
				if p.ValidatorFunc != nil {
					if err := p.ValidatorFunc(value); err != nil {
						return errors.Newf("query param %s validation failed for endpoint %s %s %s %s: %v", qp.Name, endpoint.Category, endpoint.Version, endpoint.Name, endpoint.Method, err)
					}
				}
			}
		}

		qp.value = value
		endpoint.QueryParams = append(endpoint.QueryParams, qp)
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

		// Reflect BodyType and body to ensure they match
		if endpoint.BodyType != nil {
			reflectBodyType := reflect.TypeOf(endpoint.BodyType)
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
	return func(endpoint *Endpoint, req *resty.Request) error {
		option(req)
		return nil
	}
}
