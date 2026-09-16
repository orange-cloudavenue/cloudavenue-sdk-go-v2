/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package inetworkobjects

import (
	"context"
	"log/slog"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/require"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

type fakeClient struct {
	do func(context.Context, *cav.Endpoint, ...cav.EndpointRequestOption) (*resty.Response, error)
}

func (f fakeClient) NewRequest(context.Context, *cav.Endpoint, ...cav.RequestOption) (*resty.Request, error) {
	return nil, nil
}

func (f fakeClient) NewRequestWithBackend(context.Context, cav.BackendTarget, *cav.Endpoint, ...cav.RequestOption) (*resty.Request, error) {
	return nil, nil
}

func (f fakeClient) Logger() *slog.Logger {
	return slog.Default()
}

func (f fakeClient) Do(ctx context.Context, endpoint *cav.Endpoint, opts ...cav.EndpointRequestOption) (*resty.Response, error) {
	return f.do(ctx, endpoint, opts...)
}

func (f fakeClient) DoWithBackend(context.Context, cav.BackendTarget, *cav.Endpoint, ...cav.EndpointRequestOption) (*resty.Response, error) {
	return nil, nil
}

func (f fakeClient) GetConsole() consoles.ConsoleName {
	return ""
}

func (f fakeClient) Close() error {
	return nil
}

func TestFindAppPortProfileRejectsUnexpectedGetResponseType(t *testing.T) {
	client := fakeClient{do: func(context.Context, *cav.Endpoint, ...cav.EndpointRequestOption) (*resty.Response, error) {
		return &resty.Response{Request: &resty.Request{Result: &itypes.APIResponseListAppPortProfile{}}}, nil
	}}

	profile, err := FindAppPortProfile(t.Context(), client, generator.MustGenerate("{urn:applicationPortProfile}"), "")
	require.Nil(t, profile)
	require.EqualError(t, err, "unexpected get app port profile response type *itypes.APIResponseListAppPortProfile")
}

func TestFindAppPortProfileRejectsUnexpectedListResponseType(t *testing.T) {
	client := fakeClient{do: func(context.Context, *cav.Endpoint, ...cav.EndpointRequestOption) (*resty.Response, error) {
		return &resty.Response{Request: &resty.Request{Result: &itypes.APIResponseAppPortProfile{}}}, nil
	}}

	profile, err := FindAppPortProfile(t.Context(), client, "app-1", "urn:vcloud:vdcGroup:12345678-1234-1234-1234-123456789012")
	require.Nil(t, profile)
	require.EqualError(t, err, "unexpected list app port profile response type *itypes.APIResponseAppPortProfile")
}
