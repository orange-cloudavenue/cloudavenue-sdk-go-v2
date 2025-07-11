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
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func init() {
	Endpoint{
		Name:             "JobCerberus",
		Method:           MethodGET,
		SubClient:        ClientCerberus,
		DocumentationURL: "https://swagger.cloudavenue.orange-business.com/#/Jobs/getJobById",
		PathTemplate:     "/api/customers/v1.0/jobs/{taskId}",
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
		QueryParams: []QueryParam{},
		RequestFunc: nil, // Will be set later in the Register function.
		requestInternalFunc: func(ctx context.Context, client *resty.Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error) {
			r := client.R().
				SetContext(ctx).
				SetHeader("Accept", "application/json;version="+cerberusVCDVersion)

			for _, opt := range opts {
				if err := opt(endpoint, r); err != nil {
					return nil, err
				}
			}

			return r.Get(endpoint.PathTemplate)
		},
		BodyRequestType:  nil, // No request body for this endpoint.
		BodyResponseType: cerberusJobAPIResponse{},
	}.Register()
}

// Ensure cerberus implements the jobs interface.
var _ jobsInterface = &cerberus{}

// cerberusJobCreatedAPIResponse represents the response body when a job is created
type cerberusJobCreatedAPIResponse struct {
	ID      string `json:"jobId" faker:"uuid_hyphenated"`
	Message string `json:"message" faker:"sentence"`
}

// cerberusJobAPIResponse represents an asynchronous operation in VCD.
type cerberusJobAPIResponse []struct {
	Actions []struct {
		Name    string `json:"name" faker:"word"`
		Status  string `json:"status" faker:"oneof:DONE"`
		Details string `json:"details" faker:"sentence"`
	} `json:"actions" faker:"slice_len=2"`
	Description string `json:"description" faker:"sentence"`
	Name        string `json:"name" faker:"word"`
	Status      string `json:"status" faker:"oneof:DONE"` // Status of the job.
}

// JobRefresh is a function type that defines how to refresh a job status.
func (v *cerberus) JobRefresh(httpC *resty.Client, resp *resty.Response, reqOpts []EndpointRequestOption) (job *Job, err error) {
	job, err = v.JobParser(resp)
	if err != nil {
		return job, err
	}

	ep, err := GetEndpoint("JobCerberus", MethodGET)
	if err != nil {
		return nil, errors.New("failed to get endpoint for JobCerberus: " + err.Error())
	}

	reqOpts = append(reqOpts,
		SetCustomRestyOption(func(r *resty.Request) { r.SetError(&cerberusError{}) }),
		WithPathParam(ep.PathParams[0], urn.ExtractUUID(job.ID)),
		OverrideSetResult(cerberusJobAPIResponse{}),
	)

	respR, err := ep.requestInternalFunc(resp.Request.Context(), httpC, ep, reqOpts...)
	if err != nil {
		return nil, errors.New("failed to refresh job status: " + err.Error())
	}

	return v.JobParser(respR)
}

// JobParser parses the job response body and extracts the job information.
func (v *cerberus) JobParser(resp *resty.Response) (job *Job, err error) {
	if resp == nil {
		return job, errors.New("no response to parse")
	}

	// The created job have different response structure
	// Cerberus does not respect the API convention for job creation.
	// It returns a HTTP 201 status code with a different response body.
	//
	// ! This is untestable because resp.Bytes() is indefinable in the mock.
	if resp.StatusCode() == http.StatusCreated {
		jobCreated := &cerberusJobCreatedAPIResponse{}
		if err := json.Unmarshal(resp.Bytes(), jobCreated); err == nil {
			// Continue only if the unmarshalling was successful.
			return &Job{
				ID:          jobCreated.ID,
				Description: jobCreated.Message,
			}, nil
		}
	}

	if apiR, ok := resp.Result().(*cerberusJobAPIResponse); ok {
		if len(*apiR) == 0 {
			return nil, &errors.APIError{
				StatusCode:    resp.StatusCode(),
				StatusMessage: "No job returned",
				Operation:     "Fetching job status",
				Message:       "The job response is empty",
				Duration:      resp.Duration(),
				Endpoint:      resp.Request.URL,
			}
		}

		job = &Job{
			// The taskId is used as the job ID.
			ID:          resp.Request.PathParams["taskId"],
			Name:        (*apiR)[0].Name,
			Description: (*apiR)[0].Description,
			HREF:        resp.Request.URL,
		}

		status, err := v.JobStatusParser((*apiR)[0].Status)
		if err != nil {
			return nil, errors.New("failed to parse job status: " + err.Error())
		}

		job.Status = status

		if (*apiR)[0].Status == "FAILED" {
			return job, &errors.APIError{
				StatusCode:    resp.StatusCode(),
				StatusMessage: status.String(),
				Operation:     "Fetching job status",
				Message:       (*apiR)[0].Description,
				Duration:      resp.Duration(),
				Endpoint:      resp.Request.URL,
			}
		}

		return job, nil
	}

	if err := v.ParseAPIError("JobParser", resp); err != nil {
		return nil, err
	}

	return nil, errors.New("failed to parse cerberus job response, unexpected type or empty response")
}

// Status returns the job status from the response body.
func (v *cerberus) JobStatusParser(status string) (s JobStatus, err error) {
	// CREATED, PENDING, IN_PROGRESS, FAILED, DONE
	switch strings.ToLower(status) {
	case "created":
		s = JobQueued
	case "pending":
		s = JobQueued
	case "in_progress":
		s = JobRunning
	case "failed":
		s = JobError
	case "done":
		s = JobSuccess
	default:
		return "", errors.New("unknown job status: " + status)
	}
	return s, nil
}
