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

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

func TestVmwareJobStatusParser(t *testing.T) {
	v := &vmware{}
	tests := []struct {
		input    string
		expected JobStatus
		wantErr  bool
	}{
		{"queued", JobQueued, false},
		{"preRunning", JobRunning, false},
		{"running", JobRunning, false},
		{"error", JobError, false},
		{"success", JobSuccess, false},
		{"aborted", JobAborted, false},
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

func TestVmwareJobParser_201Created(t *testing.T) {
	v := &vmware{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Location": []string{"/api/task/d3c42a20-96b9-4452-91dd-f71b71dfe314"},
			},
		},
	}
	job, err := v.JobParser(resp)

	assert.NoError(t, err)
	assert.Equal(t, "d3c42a20-96b9-4452-91dd-f71b71dfe314", job.ID)
}

func TestVmwareJobParser_201Created_BadLocationHeader(t *testing.T) {
	v := &vmware{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Location": []string{"invalid-location"},
			},
		},
	}
	job, err := v.JobParser(resp)

	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestVmwareJobParser_NormalResponse(t *testing.T) {
	v := &vmware{}

	data := &vmwareJobAPIResponse{}
	_ = faker.FakeData(data)

	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": data.ID},
			URL:        "http://example.com/job",
			Result:     data,
		},
	}
	job, err := v.JobParser(resp)

	assert.NoError(t, err)
	assert.Equal(t, data.ID, job.ID)
	assert.Equal(t, data.Name, job.Name)
	assert.Equal(t, data.Description, job.Description)
	assert.Equal(t, data.HREF, job.HREF)
	assert.Equal(t, JobSuccess, job.Status)
}

func TestVmwareJobParser_FailedStatus(t *testing.T) {
	v := &vmware{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
			Result: &vmwareJobAPIResponse{
				Status: "error",
				Error: &vmwareError{
					StatusCode:    http.StatusInternalServerError,
					StatusMessage: "Internal Server Error",
					Message:       "An error occurred while processing the job",
				},
			},
		},
	}
	job, err := v.JobParser(resp)

	assert.Error(t, err)
	assert.NotNil(t, job)
}

func TestVmwareJobParser_EmptyResponse(t *testing.T) {
	v := &vmware{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
			Result:     &vmwareJobAPIResponse{},
		},
	}
	job, err := v.JobParser(resp)
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestVmwareJobParser_NilResponse(t *testing.T) {
	v := &vmware{}
	job, err := v.JobParser(nil)
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestVmwareJobParser_UnknownJobStatus(t *testing.T) {
	v := &vmware{}
	resp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
		Request: &resty.Request{
			PathParams: map[string]string{"taskId": "id-123"},
			URL:        "http://example.com/job",
			Result: &vmwareJobAPIResponse{
				Status: "unknown_status",
			},
		},
	}
	job, err := v.JobParser(resp)
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestVmwareJobParser_VmwareErrorResponse(t *testing.T) {
	v := &vmware{}
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

func TestVmwareJobParser_InvalidResponseType(t *testing.T) {
	v := &vmware{}
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
