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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/jobs"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

// Test for JobRefresh: complete success
func TestVmware_JobRefresh_ResponseNil(t *testing.T) {
	vmw := NewVmwareClient()
	vmw.SetConsole(getMockConsole())

	if vmwJob, ok := vmw.(jobs.Client); ok {
		job, err := vmwJob.JobRefresh(nil, nil)
		assert.Error(t, err)
		assert.Nil(t, job)
	}
}

func mustResponder(responder httpmock.Responder, _ error) httpmock.Responder {
	return responder
}

// Test for JobRefresh: table-driven tests
func TestVmware_JobRefresh_TableDriven(t *testing.T) {
	type fields struct {
		jobResponders []httpmock.Responder
	}
	tests := []struct {
		name           string
		fields         fields
		statusExpected jobs.Status
		wantErr        bool
		wantNilJob     bool
		wantJobHref    string
		wantParseErr   bool
	}{
		{
			name: "success",
			fields: fields{
				jobResponders: []httpmock.Responder{
					mustResponder(httpmock.NewJsonResponder(200, VmwareJobAPIResponse{
						ID:          "123",
						Name:        "Test Job",
						Description: "This is a test job",
						HREF:        "http://mock.api/job/123",
						Status:      "success",
						Operation:   "test-operation",
					})),
				},
			},
			wantErr:        false,
			wantNilJob:     false,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: jobs.Success,
		},
		{
			name: "error",
			fields: fields{
				jobResponders: []httpmock.Responder{
					mustResponder(httpmock.NewJsonResponder(200, VmwareJobAPIResponse{
						ID:          "123",
						Name:        "Test Job",
						Description: "This is a test job",
						HREF:        "http://mock.api/job/123",
						Status:      "error",
						Operation:   "test-operation",
						Error: &VmwareError{
							StatusCode:    500,
							StatusMessage: "Internal Server Error",
							Message:       "An error occurred",
						},
					})),
				},
			},
			wantErr:        true,
			wantNilJob:     false,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: jobs.Error,
		},
		{
			name: "bad-job",
			fields: fields{
				jobResponders: []httpmock.Responder{
					mustResponder(httpmock.NewJsonResponder(200, VmwareJobAPIResponse{
						ID:          "123",
						Name:        "Test Job",
						Description: "This is a test job",
						HREF:        "http://mock.api/job/123",
						Status:      "unknown",
						Operation:   "test-operation",
					})),
				},
			},
			wantErr:        true,
			wantNilJob:     false,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: "",
		},
		{
			name: "bad-jobformat",
			fields: fields{
				jobResponders: []httpmock.Responder{
					httpmock.NewStringResponder(200, "unknown response format"),
				},
			},
			wantErr:        true,
			wantNilJob:     true,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: "",
		},
		{
			name: "error-500",
			fields: fields{
				jobResponders: []httpmock.Responder{
					httpmock.NewStringResponder(500, "Internal Server Error"),
				},
			},
			wantErr:        true,
			wantNilJob:     true,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: "",
		},
		{
			name: "error-vmware",
			fields: fields{
				jobResponders: []httpmock.Responder{
					mustResponder(httpmock.NewJsonResponder(500, VmwareError{
						StatusCode:    500,
						StatusMessage: "Internal Server Error",
						Message:       "An error occurred",
					})),
				},
			},
			wantErr:        true,
			wantNilJob:     true,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: "",
			wantParseErr:   true,
		},
		{
			name: "status-preRunning",
			fields: fields{
				jobResponders: []httpmock.Responder{
					mustResponder(httpmock.NewJsonResponder(200, VmwareJobAPIResponse{
						ID:          "123",
						Name:        "Test Job",
						Description: "This is a test job",
						HREF:        "http://mock.api/job/123",
						Status:      "preRunning", // preRunning is considered as running
						Operation:   "test-operation",
					})),
				},
			},
			wantErr:        false,
			wantNilJob:     false,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: jobs.Running, // pre-running is considered as running
		},
		{
			name: "status-queued",
			fields: fields{
				jobResponders: []httpmock.Responder{
					mustResponder(httpmock.NewJsonResponder(200, VmwareJobAPIResponse{
						ID:          "123",
						Name:        "Test Job",
						Description: "This is a test job",
						HREF:        "http://mock.api/job/123",
						Status:      "queued",
						Operation:   "test-operation",
					})),
				},
			},
			wantErr:        false,
			wantNilJob:     false,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: jobs.Queued,
		},
		{
			name: "status-running",
			fields: fields{
				jobResponders: []httpmock.Responder{
					mustResponder(httpmock.NewJsonResponder(200, VmwareJobAPIResponse{
						ID:          "123",
						Name:        "Test Job",
						Description: "This is a test job",
						HREF:        "http://mock.api/job/123",
						Status:      "running",
						Operation:   "test-operation",
					})),
				},
			},
			wantErr:        false,
			wantNilJob:     false,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: jobs.Running,
		},
		{
			name: "status-aborted",
			fields: fields{
				jobResponders: []httpmock.Responder{
					mustResponder(httpmock.NewJsonResponder(200, VmwareJobAPIResponse{
						ID:          "123",
						Name:        "Test Job",
						Description: "This is a test job",
						HREF:        "http://mock.api/job/123",
						Status:      "aborted",
						Operation:   "test-operation",
					})),
				},
			},
			wantErr:        false,
			wantNilJob:     false,
			wantJobHref:    "http://mock.api/job/123",
			statusExpected: jobs.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the NewCloudavenueCredential function to return a mock auth
			cred := auth.NewMockAuth(map[string]string{
				"Mock-Header": "mock-value",
			})

			// Setup
			vmw := NewVmwareClient()
			vmw.SetConsole(getMockConsole())
			vmw.SetCredential(cred)

			hC := httpclient.NewMockHTTPClient()
			defer httpmock.DeactivateAndReset()
			for _, responder := range tt.fields.jobResponders {
				httpmock.RegisterResponder("GET", "http://mock.api/job/123", responder)
			}

			// Mock resty.Request and resty.Response
			req := hC.R()
			req.SetHeader("Accept", "application/json")
			req.SetResult(VmwareJobAPIResponse{})
			req.SetError(VmwareError{})
			req.Result = &VmwareJobAPIResponse{
				ID:          "123",
				Name:        "Test Job",
				Description: "This is a test job",
				HREF:        "http://mock.api/job/123",
				Status:      "queued",
				Operation:   "test-operation",
			}
			req.URL = "http://mock.api/test"
			req.SetContext(t.Context())
			resp := &resty.Response{
				Request: req,
				RawResponse: &http.Response{
					Header: map[string][]string{
						"Location": {"http://mock.api/job/123"},
					},
				},
			}

			if vmwJob, ok := vmw.(jobs.Client); ok {
				job, err := vmwJob.JobRefresh(req, resp)
				if tt.wantNilJob {
					assert.Nil(t, job)
				} else {
					assert.NotNil(t, job)
					assert.Equal(t, "http://mock.api/job/123", job.HREF)
					assert.Equal(t, tt.statusExpected, job.Status)
				}
				if tt.wantErr {
					assert.Error(t, err)
					if tt.wantParseErr {
						apiErr := &errors.APIError{}
						if errors.As(err, &apiErr) {
							assert.Equal(t, 500, apiErr.StatusCode)
							assert.Equal(t, "Internal Server Error", apiErr.StatusMessage)
							assert.Equal(t, "An error occurred", apiErr.Message)
						}
					}
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
