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

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

func getMockConsole() consoles.Console {
	c, _ := consoles.FindBySiteID("mock")
	return c
}

// func getMockAuth() auth.Auth {
// 	return auth.NewCloudavenueCredential
// }

func TestVmware_NewHTTPClient_Success(t *testing.T) {
	client := httpclient.NewMockHTTPClient()
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
	cred, err := auth.NewCloudavenueCredential(client, getMockConsole(), "mockorg001", "user", "pass")
	assert.NoError(t, err)
	assert.NotNil(t, cred)

	vmw := &vmware{}
	vmw.SetConsole(getMockConsole())
	vmw.SetCredential(cred)

	httpC, err := vmw.NewHTTPClient(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, httpC)
	assert.Equal(t, "token-vmware", cred.Headers()["Authorization"][7:])
}

func TestVmware_NewHTTPClient_RefreshError(t *testing.T) {
	client := httpclient.NewMockHTTPClient()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://mock.api/cloudapi/1.0.0/sessions",
		httpmock.NewStringResponder(500, "Internal Server Error"),
	)

	cred, err := auth.NewCloudavenueCredential(client, getMockConsole(), "mockorg001", "user", "pass")
	assert.NoError(t, err)
	assert.NotNil(t, cred)

	vmw := &vmware{}
	vmw.SetConsole(getMockConsole())
	vmw.SetCredential(cred)

	httpC, err := vmw.NewHTTPClient(t.Context())
	assert.Error(t, err)
	assert.Nil(t, httpC)
}

func TestVmware_ParseAPIError_NilResponse(t *testing.T) {
	vmw := &vmware{}
	err := vmw.ParseAPIError(nil)
	assert.Nil(t, err)
}

func TestVmware_ParseAPIError_NilError(t *testing.T) {
	vmw := &vmware{}
	resp := &resty.Response{}
	err := vmw.ParseAPIError(resp)
	assert.Nil(t, err)
}

func TestVmware_ParseAPIError_ErrorResponse(t *testing.T) {
	vmw := &vmware{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: 400,
			Status:     http.StatusText(400),
			Body:       http.NoBody,
		},
		Body: http.NoBody,
		Request: &resty.Request{
			URL: "http://mock.api/vmware",
			Error: &vmwareError{
				Message:        "Bad Request",
				MinorErrorCode: "1234",
			},
		},
	}

	apiErr := vmw.ParseAPIError(resp)
	assert.NotNil(t, apiErr)
	assert.Equal(t, 400, apiErr.StatusCode)
	assert.Equal(t, "Bad Request", apiErr.Message)
	assert.Equal(t, resp.Duration(), apiErr.Duration)
	assert.Equal(t, resp.Request.URL, apiErr.Endpoint)
}

func TestVmware_ParseAPIError_ErrorResponse_Unknown(t *testing.T) {
	vmw := &vmware{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: 500,
			Status:     http.StatusText(500),
			Body:       http.NoBody,
		},
		Body: http.NoBody,
		Request: &resty.Request{
			URL:   "http://mock.api/vmware",
			Error: "Internal Server Error",
		},
	}

	apiErr := vmw.ParseAPIError(resp)
	assert.Nil(t, apiErr)
}
