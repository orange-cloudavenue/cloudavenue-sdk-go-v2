/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package subclient

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
)

func TestMock_NewHTTPClient_Success(t *testing.T) {
	_ = httpclient.NewMockHTTPClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://mock.api/cloudapi/1.0.0/sessions",
		httpmock.NewStringResponder(200, "ok").
			HeaderSet(map[string][]string{
				auth.XVmwareAccessToken: {"token-vmware"},
			}),
	)

	// Mock the NewCloudavenueCredential function to return a mock auth
	cred := auth.NewMockAuth(map[string]string{
		"Mock-Header": "mock-value",
	})
	assert.NotNil(t, cred)

	vmw := &mockClient{}
	vmw.SetConsole(getMockConsole())
	vmw.SetCredential(cred)

	httpC, err := vmw.NewHTTPClient(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, httpC)
	assert.Equal(t, "mock-value", cred.Headers()["Mock-Header"])
}

func TestMock_ParseAPIError_NilResponse(t *testing.T) {
	vmw := &mockClient{}
	err := vmw.ParseAPIError(nil)
	assert.Nil(t, err)
}

func TestMock_ParseAPIError_NilError(t *testing.T) {
	vmw := &mockClient{}
	resp := &resty.Response{}
	err := vmw.ParseAPIError(resp)
	assert.Nil(t, err)
}

func TestMock_ParseAPIError_ErrorResponse(t *testing.T) {
	vmw := &mockClient{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: 400,
			Status:     http.StatusText(400),
			Body:       http.NoBody,
		},
		Body: http.NoBody,
		Request: &resty.Request{
			URL: "http://mock.api/mock",
			Error: &MockError{
				Message: "Bad Request",
			},
			Time: time.Now(),
		},
	}

	apiErr := vmw.ParseAPIError(resp)
	assert.NotNil(t, apiErr)
	assert.Equal(t, 400, apiErr.StatusCode)
	assert.Equal(t, "Bad Request", apiErr.Message)
	assert.Equal(t, resp.Duration(), apiErr.Duration)
	assert.Equal(t, resp.Request.URL, apiErr.Endpoint)
}

func TestMock_ParseAPIError_ErrorResponse_Unknown(t *testing.T) {
	vmw := &mockClient{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: 500,
			Status:     http.StatusText(500),
			Body:       http.NoBody,
		},
		Body: http.NoBody,
		Request: &resty.Request{
			URL:   "http://mock.api/mock",
			Error: "Internal Server Error",
		},
	}

	apiErr := vmw.ParseAPIError(resp)
	assert.Nil(t, apiErr)
}
