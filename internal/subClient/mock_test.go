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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/jobs"
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

	vmw := &MockClient{}
	vmw.SetConsole(getMockConsole())
	vmw.SetCredential(cred)

	httpC, err := vmw.NewHTTPClient(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, httpC)
	assert.Equal(t, "mock-value", cred.Headers()["Mock-Header"])
}

func TestMock_ParseAPIError_NilResponse(t *testing.T) {
	vmw := &MockClient{}
	err := vmw.ParseAPIError(nil)
	assert.Nil(t, err)
}

func TestMock_ParseAPIError_NilError(t *testing.T) {
	vmw := &MockClient{}
	resp := &resty.Response{}
	err := vmw.ParseAPIError(resp)
	assert.Nil(t, err)
}

func TestMock_ParseAPIError_ErrorResponse(t *testing.T) {
	vmw := &MockClient{}
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
	vmw := &MockClient{}
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

func TestMockJobClientWithJob_SetCredential_And_SetConsole(t *testing.T) {
	mockJobClient := &MockJobClientWithJob{
		Client: &MockClient{},
	}
	mockAuth := auth.NewMockAuth(map[string]string{"Test": "value"})
	mockConsole := getMockConsole()

	mockJobClient.SetCredential(mockAuth)
	assert.Equal(t, mockAuth, mockJobClient.Client.(*MockClient).credential)

	mockJobClient.SetConsole(mockConsole)
	assert.Equal(t, mockConsole, mockJobClient.Client.(*MockClient).console)
}

func TestMockWithJob_NewHTTPClient_Success(t *testing.T) {
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

	vmw := NewMockJobClient()
	vmw.SetConsole(getMockConsole())
	vmw.SetCredential(cred)

	httpC, err := vmw.NewHTTPClient(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, httpC)
	assert.Equal(t, "mock-value", cred.Headers()["Mock-Header"])
}

func TestMockWithJob_ParseAPIError_NilResponse(t *testing.T) {
	vmw := NewMockJobClient()
	err := vmw.ParseAPIError(nil)
	assert.Nil(t, err)
}

func TestMockWithJob_ParseAPIError_NilError(t *testing.T) {
	vmw := NewMockJobClient()
	resp := &resty.Response{}
	err := vmw.ParseAPIError(resp)
	assert.Nil(t, err)
}

func TestMockWithJob_ParseAPIError_ErrorResponse(t *testing.T) {
	vmw := NewMockJobClient()
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

func TestMockWithJob_ParseAPIError_ErrorResponse_Unknown(t *testing.T) {
	vmw := NewMockJobClient()
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

func TestMockJobClientWithJob_JobRefresh_Success(t *testing.T) {
	vmw := NewMockJobClient()

	if vmwJob, ok := vmw.(jobs.Client); ok {
		job, err := vmwJob.JobRefresh(nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, job)
		// assert.Equal(t, "http://mock.api/job/123", job.HREF)
	}
}
