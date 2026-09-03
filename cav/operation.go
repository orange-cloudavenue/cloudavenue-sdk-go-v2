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

	"resty.dev/v3"
)

// BackendTarget identifies backend used by an operation.
type BackendTarget int

const (
	BackendInfrapi BackendTarget = iota + 1
	BackendVMware
	BackendOSE
	BackendNetBackup
)

// Response wraps raw resty responses for typed operations.
type Response struct {
	Raw *resty.Response
}

// Result returns decoded response payload when present.
func (r *Response) Result() any {
	if r == nil || r.Raw == nil {
		return nil
	}
	return r.Raw.Result()
}

// StatusCode returns HTTP status code when response is present.
func (r *Response) StatusCode() int {
	if r == nil || r.Raw == nil {
		return 0
	}
	return r.Raw.StatusCode()
}

// Operation defines compile-time typed leaf operation behavior.
type Operation[P any, R any] struct {
	Name           string
	Backend        BackendTarget
	Endpoint       *Endpoint
	Validate       func(P) error
	RequestOptions func(P) ([]EndpointRequestOption, error)
	Transform      func(P) (any, error)
	Extract        func(*Response, P) (R, error)
}

// PartialSuccessError wraps errors raised after side effects already happened.
// Callers can inspect Result even when Error returns non-nil.
type PartialSuccessError[R any] struct {
	Result R
	Err    error
}

// Execute validates params, performs request, and extracts typed result.
func Execute[P, R any](ctx context.Context, c Client, op Operation[P, R], params P) (R, error) {
	var zero R

	if op.Name == "" {
		return zero, fmt.Errorf("operation name is required")
	}
	if op.Endpoint == nil {
		return zero, fmt.Errorf("%s: endpoint is required", op.Name)
	}
	if op.Extract == nil {
		return zero, fmt.Errorf("%s: extract is required", op.Name)
	}

	if op.Validate != nil {
		if err := op.Validate(params); err != nil {
			return zero, fmt.Errorf("%s: validate: %w", op.Name, err)
		}
	}

	var body any
	if op.Transform != nil {
		var err error
		body, err = op.Transform(params)
		if err != nil {
			return zero, fmt.Errorf("%s: transform: %w", op.Name, err)
		}
	}

	var requestOpts []EndpointRequestOption
	if op.RequestOptions != nil {
		var err error
		requestOpts, err = op.RequestOptions(params)
		if err != nil {
			return zero, fmt.Errorf("%s: request options: %w", op.Name, err)
		}
	}

	if body != nil {
		requestOpts = append(requestOpts, SetBody(body))
	}

	backend := op.Endpoint.Backend
	if op.Backend != 0 {
		if op.Backend != op.Endpoint.Backend {
			return zero, fmt.Errorf("%s: operation backend %d does not match endpoint backend %d", op.Name, op.Backend, op.Endpoint.Backend)
		}
		backend = op.Backend
	}

	resp, err := doOperationRequest(ctx, c, backend, op.Endpoint, requestOpts...)
	if err != nil {
		return zero, fmt.Errorf("%s: %w", op.Name, err)
	}

	result, err := op.Extract(resp, params)
	if err != nil {
		return zero, fmt.Errorf("%s: extract: %w", op.Name, err)
	}

	return result, nil
}

// Error implements error.
func (e *PartialSuccessError[R]) Error() string {
	if e == nil || e.Err == nil {
		return "partial success"
	}
	return e.Err.Error()
}

// Unwrap returns wrapped error.
func (e *PartialSuccessError[R]) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func doOperationRequest(ctx context.Context, c Client, backend BackendTarget, endpoint *Endpoint, opts ...EndpointRequestOption) (*Response, error) {
	raw, err := c.DoWithBackend(ctx, backend, endpoint, opts...)
	if err != nil {
		return nil, err
	}

	return &Response{Raw: raw}, nil
}
