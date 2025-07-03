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
	"testing"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

type dummyBody struct {
	Foo string
}

func TestSetBody_NilBody(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "POST",
		BodyRequestType: dummyBody{},
	}
	req := resty.New().R()
	opt := SetBody(nil)
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for nil body, got nil")
	}
}

func TestSetBody_BodyRequestTypeMismatch(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "POST",
		BodyRequestType: dummyBody{},
	}
	req := resty.New().R()
	opt := SetBody("not a struct")
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for type mismatch, got nil")
	}
}

func TestSetBody_BodyRequestTypeMatch(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "POST",
		BodyRequestType: dummyBody{},
	}
	req := resty.New().R()
	body := dummyBody{Foo: "bar"}
	opt := SetBody(body)
	err := opt(endpoint, req)
	if err != nil {
		t.Errorf("expected no error for correct type, got: %v", err)
	}
	if req.Body == nil {
		t.Error("expected body to be set in request")
	}
}

func TestSetBody_NoBodyRequestType(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "POST",
		BodyRequestType: nil,
	}
	req := resty.New().R()
	body := dummyBody{Foo: "bar"}
	opt := SetBody(body)
	err := opt(endpoint, req)
	if err != nil {
		t.Errorf("expected no error when BodyRequestType is nil, got: %v", err)
	}
	if req.Body == nil {
		t.Error("expected body to be set in request")
	}
}

func TestWithQueryParam_NoQueryParams(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
		QueryParams: nil,
	}
	req := resty.New().R()
	qp := QueryParam{Name: "foo"}
	opt := WithQueryParam(qp, "bar")
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for nil QueryParams, got nil")
	}
}

func TestWithQueryParam_RequiredMissing(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
		QueryParams: []QueryParam{{Name: "foo", Required: true}},
	}
	req := resty.New().R()
	qp := QueryParam{Name: "foo", Required: true}
	opt := WithQueryParam(qp, "")
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for missing required query param, got nil")
	}
}

func TestWithQueryParam_ValidatorFails(t *testing.T) {
	validator := func(_ string) error {
		return errors.Newf("invalid value")
	}
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
		QueryParams: []QueryParam{{Name: "foo", ValidatorFunc: validator}},
	}
	req := resty.New().R()
	qp := QueryParam{Name: "foo", ValidatorFunc: validator}
	opt := WithQueryParam(qp, "bad")
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for validator failure, got nil")
	}
}

func TestWithQueryParam_Success(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
		QueryParams: []QueryParam{{Name: "foo"}},
	}
	req := resty.New().R()
	qp := QueryParam{Name: "foo"}
	opt := WithQueryParam(qp, "bar")
	err := opt(endpoint, req)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if req.QueryParams.Get("foo") != "bar" {
		t.Error("expected query param to be set in request")
	}
}

func TestOverrideSetResult_NilResult(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
	}
	req := resty.New().R()
	opt := OverrideSetResult(nil)
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for nil result type, got nil")
	}
}

func TestOverrideSetResult_Success(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
	}
	req := resty.New().R()
	result := struct{ Foo string }{}
	opt := OverrideSetResult(&result)
	err := opt(endpoint, req)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	// No direct way to check SetResult, but no panic/error is sufficient
}

func TestWithPathParam_NoPathParams(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
		PathParams: nil,
	}
	req := resty.New().R()
	pp := PathParam{Name: "foo"}
	opt := WithPathParam(pp, "bar")
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for nil PathParams, got nil")
	}
}

func TestWithPathParam_RequiredMissing(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
		PathParams: []PathParam{{Name: "foo", Required: true}},
	}
	req := resty.New().R()
	pp := PathParam{Name: "foo", Required: true}
	opt := WithPathParam(pp, "")
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for missing required path param, got nil")
	}
}

func TestWithPathParam_ValidatorFails(t *testing.T) {
	validator := func(_ string) error {
		return errors.Newf("invalid value")
	}
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
		PathParams: []PathParam{{Name: "foo", ValidatorFunc: validator}},
	}
	req := resty.New().R()
	pp := PathParam{Name: "foo", ValidatorFunc: validator}
	opt := WithPathParam(pp, "bad")
	err := opt(endpoint, req)
	if err == nil || err.Error() == "" {
		t.Error("expected error for validator failure, got nil")
	}
}

func TestWithPathParam_Success(t *testing.T) {
	endpoint := &Endpoint{
		Category: "cat", Version: "v1", Name: "name", Method: "GET",
		PathParams: []PathParam{{Name: "foo"}},
	}
	req := resty.New().R()
	pp := PathParam{Name: "foo"}
	opt := WithPathParam(pp, "bar")
	err := opt(endpoint, req)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if req.PathParams["foo"] != "bar" {
		t.Error("expected path param to be set in request")
	}
}
