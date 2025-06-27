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
	"context"

	"resty.dev/v3"

	"github.com/jarcoal/httpmock"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/jobs"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

// MockClient and MockJobClientWithJob are mock implementations of the Client interface
// There are used for testing purposes, allowing you to simulate API responses

var _ Client = &MockClient{}

type MockClient struct {
	client
}

type MockError struct {
	Message string `json:"message"`
}

var NewMockClient = func() Client {
	return &MockClient{}
}

func (m *MockClient) NewHTTPClient(_ context.Context) (*resty.Client, error) {
	if m.httpClient == nil {
		m.httpClient = httpclient.NewMockHTTPClient().
			SetBaseURL("https://mock-api.cloudavenue.com").
			SetHeaders(m.credential.Headers()).
			SetError(MockError{})
		httpmock.ActivateNonDefault(m.httpClient.Client())
	}
	return m.httpClient, nil
}

func (m *MockClient) SetCredential(a auth.Auth) {
	m.credential = a
}

func (m *MockClient) SetConsole(c consoles.Console) {
	m.console = c
}

func (m *MockClient) ParseAPIError(resp *resty.Response) *errors.APIError {
	if resp == nil || !resp.IsError() {
		return nil
	}

	// Mock error parsing logic
	if err, ok := resp.Error().(*MockError); ok {
		return &errors.APIError{
			StatusCode: resp.StatusCode(),
			Message:    err.Message,
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
		}
	}

	return nil
}

// * WithJob interface

var _ jobs.Client = &MockJobClientWithJob{}

type MockJobClientWithJob struct {
	// Set the desired job status for the mock client
	JobStatus jobs.Status
	Client
}

var NewMockJobClient = func() Client {
	return &MockJobClientWithJob{
		JobStatus: jobs.Queued,
		Client:    NewMockClient(),
	}
}

func (m *MockJobClientWithJob) SetCredential(a auth.Auth) {
	m.Client.SetCredential(a)
}

func (m *MockJobClientWithJob) SetConsole(c consoles.Console) {
	m.Client.SetConsole(c)
}

func (m *MockJobClientWithJob) NewHTTPClient(ctx context.Context) (*resty.Client, error) {
	return m.Client.NewHTTPClient(ctx)
}

func (m *MockJobClientWithJob) ParseAPIError(resp *resty.Response) *errors.APIError {
	return m.Client.ParseAPIError(resp)
}

func (m *MockJobClientWithJob) JobRefresh(_ *resty.Request, resp *resty.Response) (*jobs.Job, error) {
	return m.JobParser(resp)
}

func (m *MockJobClientWithJob) JobParser(_ *resty.Response) (*jobs.Job, error) {
	status, _ := m.JobStatusParser("")

	// Mock job parsing logic
	return &jobs.Job{
		ID:     "mock-job-id",
		Status: status,
	}, nil
}

func (m *MockJobClientWithJob) JobStatusParser(_ string) (jobs.Status, error) {
	return m.JobStatus, nil
}
