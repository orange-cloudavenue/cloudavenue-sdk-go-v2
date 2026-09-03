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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/orange-cloudavenue/common-go/urn"
	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

//go:generate endpoint-generator -path subclient_cerberus_jobs.go -filename zz_cav_cerberus_jobs.go -output cav_cerberus_jobs.go

func init() {
	Endpoint{
		Name:             "GetJobCerberus",
		Description:      "Get Cerberus Job",
		Method:           MethodGET,
		Backend:          BackendInfrapi,
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
		QueryParams:     []QueryParam{},
		BodyRequestType: nil, // No request body for this endpoint.
		ResponseType:    CerberusJobAPIResponse{},
	}.Register()
}

// Ensure cerberus implements jobsInterface.
var _ jobsInterface = &cerberus{}

// cerberusJobCreatedAPIResponse is returned when Cerberus creates job.
type cerberusJobCreatedAPIResponse struct {
	ID      string `json:"jobId" fake:"{uuid}"`
	Message string `json:"message" fake:"{sentence}"`
}

// CerberusJobCreatedAPIResponse aliases cerberusJobCreatedAPIResponse.
type CerberusJobCreatedAPIResponse = cerberusJobCreatedAPIResponse

// CerberusJobAPIResponse describes Cerberus job lookup response.
type CerberusJobAPIResponse []struct {
	Actions     []CerberusJobAPIResponseAction `json:"actions" fakesize:"3"`
	Description string                         `json:"description" fake:"{sentence}"`
	Name        string                         `json:"name" fake:"{word}"`
	Status      string                         `json:"status" fake:"DONE"` // Status of the job.
}

type CerberusJobAPIResponseAction struct {
	Name    string `json:"name" fake:"{word}"`
	Status  string `json:"status" fake:"DONE"`
	Details string `json:"details" fake:"{sentence}"`
}

// JobRefresh refreshes Cerberus job status.
func (v *cerberus) JobRefresh(httpC *resty.Client, resp *resty.Response, reqOpts []EndpointRequestOption) (job *Job, err error) {
	job, err = v.JobParser(resp)
	if err != nil {
		return job, err
	}

	ep, err := GetEndpoint("GetJobCerberus")
	if err != nil {
		return nil, errors.New("failed to get endpoint for JobCerberus: " + err.Error())
	}

	reqOpts = append(
		reqOpts,
		SetCustomRestyOption(
			func(r *resty.Request) {
				r.SetResultError(&cerberusError{})
				r.SetResult(&CerberusJobAPIResponse{})
			},
		),
		WithPathParam(ep.PathParams[0], urn.ExtractUUID(job.ID)),
	)

	r := httpC.R().
		SetContext(resp.Request.Context()).
		SetHeader("Accept", "application/json;version="+cerberusVCDVersion)

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

// JobParser parses Cerberus job data from response.
func (v *cerberus) JobParser(resp *resty.Response) (job *Job, err error) {
	if resp == nil {
		return job, errors.New("no response to parse")
	}

	// Cerberus returns a different body shape for HTTP 201 job creation responses.
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

	if apiR, ok := resp.Result().(*CerberusJobAPIResponse); ok {
		if len(*apiR) == 0 {
			return nil, &errors.APIError{
				StatusCode:    resp.StatusCode(),
				StatusMessage: "No job returned",
				Operation:     "Fetching job status",
				Message:       "The job response is empty",
				Duration:      resp.Duration(),
				Endpoint:      resp.Request.URL,
				Err:           classifyStatusCode(resp.StatusCode()),
			}
		}

		job = &Job{
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
				Err:           errors.ErrJobFailed,
			}
		}

		return job, nil
	}

	if err := v.parseAPIError("JobParser", resp); err != nil {
		return nil, err
	}

	return nil, errors.New("failed to parse cerberus job response: unexpected type or empty response")
}

// JobStatusParser maps Cerberus status strings to JobStatus.
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
