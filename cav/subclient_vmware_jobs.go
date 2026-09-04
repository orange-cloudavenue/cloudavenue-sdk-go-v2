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
	"strings"

	"github.com/orange-cloudavenue/common-go/urn"
	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

func init() {
	Endpoint{
		Name:             "GetJobVmware",
		Description:      "Get VMware Job",
		Method:           MethodGET,
		Backend:          BackendVMware,
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/types/TaskType.html",
		PathTemplate:     "/api/task/{taskId}",
		PathParams: []PathParam{
			{
				Name:        "taskId",
				Description: "The identifier of the task to retrieve.",
				Required:    true,
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "required,uuid4")
				},
			},
		},
		QueryParams:     []QueryParam{},
		BodyRequestType: nil, // No request body for this endpoint.
		ResponseType:    vmwareJobAPIResponse{},
	}.Register()
}

// Ensure vmware implements jobsInterface.
var _ jobsInterface = &vmware{}

// vmwareJobAPIResponse describes VMware task response.
type vmwareJobAPIResponse struct {
	HREF             string       `json:"href,omitempty" fake:"{href_uuid}"` // The URI of the entity.
	ID               string       `json:"id,omitempty" fake:"{uuid}"`        // The entity identifier, expressed in URN format. The value of this attribute uniquely identifies the entity, persists for the life of the entity, and is never reused.
	OperationKey     string       `json:"operationKey,omitempty"`            // Optional unique identifier to support idempotent semantics for create and delete operations.
	Name             string       `json:"name,omitempty" fake:"{word}"`      // The name of the entity.
	Status           string       `json:"status,omitempty" fake:"success"`   // The execution status of the task. One of queued, preRunning, running, success, error, aborted
	Operation        string       `json:"operation,omitempty"`               // A message describing the operation that is tracked by this task.
	OperationName    string       `json:"operationName,omitempty"`           // The short name of the operation that is tracked by this task.
	ServiceNamespace string       `json:"serviceNamespace,omitempty"`        // Identifier of the service that created the task. It must not start with com.vmware.vcloud and the length must be between 1 and 128 symbols.
	StartTime        string       `json:"startTime,omitempty"`               // The date and time the system started executing the task. May not be present if the task has not been executed yet.
	EndTime          string       `json:"endTime,omitempty"`                 // The date and time that processing of the task was completed. May not be present if the task is still being executed.
	ExpiryTime       string       `json:"expiryTime,omitempty"`              // The date and time at which the task resource will be destroyed and no longer available for retrieval. May not be present if the task has not been executed or is still being executed.
	CancelRequested  bool         `json:"cancelRequested,omitempty"`         // Whether user has requested this processing to be canceled.
	Description      string       `json:"description,omitempty" `            // Optional description.
	Error            *vmwareError `json:"error,omitempty" fake:"-"`          // Represents error information from a failed task.
	Progress         int          `json:"progress,omitempty"`                // Read-only indicator of task progress as an approximate percentage between 0 and 100. Not available for all tasks.
	Details          string       `json:"details,omitempty"`                 // Detailed message about the task. Also contained by the Owner entity when task status is preRunning.
}

// JobRefresh refreshes VMware job status.
func (v *vmware) JobRefresh(httpC *resty.Client, resp *resty.Response, reqOpts []EndpointRequestOption) (job *Job, err error) {
	job, err = v.JobParser(resp)
	if err != nil {
		return job, err
	}

	ep, err := GetEndpoint("GetJobVmware")
	if err != nil {
		return nil, errors.New("failed to get endpoint for GetJobVmware: " + err.Error())
	}

	reqOpts = append(
		reqOpts,
		SetCustomRestyOption(func(r *resty.Request) { r.SetResultError(&vmwareError{}) }),
		WithPathParam(ep.PathParams[0], urn.ExtractUUID(job.ID)),
		OverrideSetResult(vmwareJobAPIResponse{}),
	)

	r := httpC.R().
		SetContext(resp.Request.Context()).
		SetHeader("Accept", "application/*+json;version="+vmwareVCDVersion)

	for _, opt := range reqOpts {
		if err := opt(ep, r); err != nil {
			return nil, err
		}
	}

	respR, err := r.Get(ep.PathTemplate)
	if err != nil {
		return nil, errors.New("failed to refresh job status: " + err.Error())
	}

	return v.JobParser(respR)
}

// JobParser parses VMware job data from response.
func (v *vmware) JobParser(resp *resty.Response) (job *Job, err error) {
	if resp == nil {
		return job, errors.New("no response to parse")
	}

	// Job creation responses expose task location in header instead of body.
	if href := resp.Header().Get("Location"); href != "" {
		job = &Job{
			ID: func() string {
				parts := strings.Split(href, "/")
				if len(parts) > 0 {
					if validators.New().Var(parts[len(parts)-1], "uuid4") == nil {
						return parts[len(parts)-1]
					}
				}
				return ""
			}(),
		}

		if job.ID == "" {
			return nil, errors.New("failed to parse vmware job ID from response header")
		}

		return job, nil
	}

	if apiR, ok := resp.Result().(*vmwareJobAPIResponse); ok && apiR.Status != "" {
		job = &Job{
			ID:          apiR.ID,
			Name:        apiR.Name,
			Description: apiR.Description,
			HREF:        apiR.HREF,
		}

		status, err := v.JobStatusParser(apiR.Status)
		if err != nil {
			return nil, errors.New("failed to parse vmware job status: " + err.Error())
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
				Err:           classifyStatusCode(apiR.Error.StatusCode),
			}
		}

		return job, nil
	}

	if err := v.parseAPIError("JobParser", resp); err != nil {
		return nil, err
	}

	return nil, errors.New("failed to parse vmware job response: unexpected type or empty response")
}

// JobStatusParser maps VMware status strings to JobStatus.
func (v *vmware) JobStatusParser(status string) (s JobStatus, err error) {
	switch status {
	case "queued":
		s = JobQueued
	case "preRunning":
		s = JobRunning
	case "running":
		s = JobRunning
	case "success":
		s = JobSuccess
	case "error":
		s = JobError
	case "aborted":
		s = JobAborted
	default:
		return s, errors.New("unknown job status: " + status)
	}
	return s, nil
}
