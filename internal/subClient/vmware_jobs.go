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
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/jobs"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ jobs.Client = &vmware{}

// VmwareJobAPIResponse represents an asynchronous operation in VMware Cloud Director.
type VmwareJobAPIResponse struct {
	HREF             string       `json:"href,omitempty"`             // The URI of the entity.
	ID               string       `json:"id,omitempty"`               // The entity identifier, expressed in URN format. The value of this attribute uniquely identifies the entity, persists for the life of the entity, and is never reused.
	OperationKey     string       `json:"operationKey,omitempty"`     // Optional unique identifier to support idempotent semantics for create and delete operations.
	Name             string       `json:"name,omitempty"`             // The name of the entity.
	Status           string       `json:"status,omitempty"`           // The execution status of the task. One of queued, preRunning, running, success, error, aborted
	Operation        string       `json:"operation,omitempty"`        // A message describing the operation that is tracked by this task.
	OperationName    string       `json:"operationName,omitempty"`    // The short name of the operation that is tracked by this task.
	ServiceNamespace string       `json:"serviceNamespace,omitempty"` // Identifier of the service that created the task. It must not start with com.vmware.vcloud and the length must be between 1 and 128 symbols.
	StartTime        string       `json:"startTime,omitempty"`        // The date and time the system started executing the task. May not be present if the task has not been executed yet.
	EndTime          string       `json:"endTime,omitempty"`          // The date and time that processing of the task was completed. May not be present if the task is still being executed.
	ExpiryTime       string       `json:"expiryTime,omitempty"`       // The date and time at which the task resource will be destroyed and no longer available for retrieval. May not be present if the task has not been executed or is still being executed.
	CancelRequested  bool         `json:"cancelRequested,omitempty"`  // Whether user has requested this processing to be canceled.
	Description      string       `json:"description,omitempty"`      // Optional description.
	Error            *VmwareError `json:"error,omitempty"`            // Represents error information from a failed task.
	Progress         int          `json:"progress,omitempty"`         // Read-only indicator of task progress as an approximate percentage between 0 and 100. Not available for all tasks.
	Details          string       `json:"details,omitempty"`          // Detailed message about the task. Also contained by the Owner entity when task status is preRunning.
}

// JobRefresh is a function type that defines how to refresh a job status.
func (v *vmware) JobRefresh(newReq *resty.Request, resp *resty.Response) (job *jobs.Job, err error) {
	job, err = v.JobParser(resp)
	if err != nil {
		return job, err
	}

	respR, err := newReq.
		SetHeader("Accept", "application/*+json;version="+auth.VDCVersion).
		SetResult(VmwareJobAPIResponse{}).
		SetError(VmwareError{}).
		Get(job.HREF)

	if err != nil {
		return job, err
	}

	if respR.IsError() {
		if apiErr := v.ParseAPIError(respR); apiErr != nil {
			return nil, apiErr
		}
		return nil, respR.Err
	}

	return v.JobParser(respR)
}

// JobParser parses the job response body and extracts the job information.
func (v *vmware) JobParser(resp *resty.Response) (job *jobs.Job, err error) {
	if resp == nil {
		return job, errors.New("no response to parse")
	}

	// If the response is not of type VmwareJobAPIResponse, find the HREF in the header
	if href := resp.Header().Get("Location"); href != "" {
		job = &jobs.Job{
			HREF: href,
		}
		return job, nil
	}

	if apiR, ok := resp.Result().(*VmwareJobAPIResponse); ok && apiR.Status != "" {
		job = &jobs.Job{
			ID:          apiR.ID,
			Name:        apiR.Name,
			Description: apiR.Description,
			HREF:        apiR.HREF,
		}

		status, err := v.JobStatusParser(apiR.Status)
		if err != nil {
			return job, errors.New("failed to parse vmware job status: " + err.Error())
		}

		job.Status = status

		if apiR.Error != nil {
			return job, &errors.APIError{
				StatusCode:    apiR.Error.StatusCode,
				StatusMessage: apiR.Error.StatusMessage,
				Operation:     apiR.Operation,
				Message:       apiR.Error.Message,
				Duration:      resp.Duration(),
				Endpoint:      apiR.HREF,
			}
		}

		return job, nil
	}

	return job, errors.New("failed to parse vmware job response, expected VmwareJobAPIResponse type")
}

// Status returns the job status from the response body.
func (v *vmware) JobStatusParser(status string) (s jobs.Status, err error) {
	switch status {
	case "queued":
		s = jobs.Queued
	case "preRunning":
		s = jobs.Running // preRunning is considered as running
	case "running":
		s = jobs.Running
	case "success":
		s = jobs.Success
	case "error":
		s = jobs.Error
	case "aborted":
		s = jobs.Aborted
	default:
		return s, errors.New("unknown job status: " + status)
	}
	return s, nil
}
