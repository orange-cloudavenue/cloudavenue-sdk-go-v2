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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

func TestCerberusJobStatusParser(t *testing.T) {
	v := &cerberus{}
	tests := []struct {
		input    string
		expected JobStatus
		wantErr  bool
	}{
		{"created", JobQueued, false},
		{"pending", JobQueued, false},
		{"in_progress", JobRunning, false},
		{"failed", JobError, false},
		{"done", JobSuccess, false},
		{"unknown", JobStatus(""), true},
	}
	for _, tt := range tests {
		got, err := v.JobStatusParser(tt.input)
		if tt.wantErr {
			assert.Error(t, err, tt.input)
		} else {
			assert.NoError(t, err, tt.input)
			assert.Equal(t, tt.expected, got, tt.input)
		}
	}
}

// * This test is commented out because it requires a specific response structure that is not mockable.
// Due to resty.Bytes() being undefined in the mock
//
// func TestCerberusJobParser_201Created(t *testing.T) {
// 	v := &cerberus{}
// 	resp := &resty.Response{
// 		RawResponse: &http.Response{
// 			StatusCode: http.StatusCreated,
// 		},
// 		Request: &resty.Request{
// 			PathParams: map[string]string{"taskId": "d3c42a20-96b9-4452-91dd-f71b71dfe314"},
// 			URL:        "http://example.com/job",
// 			Result:     &cerberusJobCreatedAPIResponse{ID: "d3c42a20-96b9-4452-91dd-f71b71dfe314", Message: "Job created successfully"},
// 		},
// 	}
// 	job, err := v.JobParser(resp)
//
// 	assert.NoError(t, err)
// 	assert.Equal(t, "d3c42a20-96b9-4452-91dd-f71b71dfe314", job.ID)
// 	assert.Equal(t, "Job created successfully", job.Description)
// 	assert.Equal(t, JobQueued, job.Status)
// }

func TestCerberusJobParser_NormalResponse(t *testing.T) {
	v := &cerberus{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
			Result: &CerberusJobAPIResponse{
				{
					Name:        "test-job",
					Description: "desc",
					Status:      "done",
				},
			},
		},
	}
	job, err := v.JobParser(resp)

	assert.NoError(t, err)
	assert.Equal(t, "id-123", job.ID)
	assert.Equal(t, "test-job", job.Name)
	assert.Equal(t, "desc", job.Description)
	assert.Equal(t, "http://example.com/job", job.HREF)
	assert.Equal(t, JobSuccess, job.Status)
}

func TestCerberusJobParser_FailedStatus(t *testing.T) {
	v := &cerberus{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
			Result: &CerberusJobAPIResponse{
				{
					Name:        "test-job",
					Description: "desc",
					Status:      "FAILED",
				},
			},
		},
	}
	job, err := v.JobParser(resp)

	assert.Error(t, err)
	assert.NotNil(t, job)
}

func TestCerberusJobParser_EmptyResponse(t *testing.T) {
	v := &cerberus{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
			Result:     &CerberusJobAPIResponse{},
		},
	}
	job, err := v.JobParser(resp)
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestCerberusJobParser_NilResponse(t *testing.T) {
	v := &cerberus{}
	job, err := v.JobParser(nil)
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestCerberusJobParser_UnknownJobStatus(t *testing.T) {
	v := &cerberus{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
			Result: &CerberusJobAPIResponse{
				{
					Name:        "test-job",
					Description: "desc",
					Status:      "unknown_status",
				},
			},
		},
	}
	job, err := v.JobParser(resp)
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestCerberusJobParser_CerberusErrorResponse(t *testing.T) {
	v := &cerberus{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusInternalServerError,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
		},
	}
	job, err := v.JobParser(resp)
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestCerberusJobParser_InvalidResponseType(t *testing.T) {
	v := &cerberus{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
			Result:     "invalid response type",
		},
	}
	job, err := v.JobParser(resp)
	assert.Error(t, err)
	assert.Nil(t, job)
}
