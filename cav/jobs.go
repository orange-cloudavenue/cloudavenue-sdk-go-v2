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
	"fmt"
	"log/slog"

	"resty.dev/v3"
)

type (
	jobsInterface interface {
		// JobRefresh refreshes the job status.
		// Error return a custom error if the job fails. errors.JobError
		JobRefresh(httpC *resty.Client, resp *resty.Response, reqOpts []EndpointRequestOption) (job *Job, err error)

		// JobParser parses the job response.
		// This method is used to parse the job response and extract the job information.
		JobParser(resp *resty.Response) (job *Job, err error)

		// JobStatusParser returns the job status.
		// This method is used to get the job status from the response body.
		JobStatusParser(status string) (s JobStatus, err error)
	}
)

// NewJobMiddleware creates a new job middleware for handling job responses.
func newJobMiddleware(httpC *resty.Client, c jobsInterface, jobOpts *JobOptions) resty.ResponseMiddleware {
	return func(_ *resty.Client, resp *resty.Response) error {
		if jobOpts == nil {
			return fmt.Errorf("job options cannot be nil, use NewJobOptions() to create a valid job options")
		}

		if jobOpts.extractorFunc != nil {
			// If an extractor function is provided, use it to extract extra data from the response.
			// This allows for custom handling of the response data.
			// Add custom middleware to the request to handle job responses.
			httpC.AddResponseMiddleware(extractorFuncMiddleware(jobOpts.extractorFunc))
		}

		// Create a new request for job status checking.
		// This request will be used to poll the job status until it is terminated.
		// The request will be configured with retry conditions and timeout settings.
		reqOpts := []EndpointRequestOption{}
		reqOpts = append(reqOpts,
			// Set the context from the original response request.
			SetCustomRestyOption(func(r *resty.Request) { r.SetContext(resp.Request.Context()) }),
			// Set retry conditions for the job status check.
			// This will retry the request based on the job status and error conditions.
			SetCustomRestyOption(func(r *resty.Request) { r.SetRetryConditions(jobRetryCondition(c)) }),
			SetCustomRestyOption(func(r *resty.Request) { r.SetRetryWaitTime(jobOpts.PollInterval) }),
			SetCustomRestyOption(func(r *resty.Request) { r.SetRetryMaxWaitTime(jobOpts.Timeout) }),
			// Set the maximum number of retries based on the timeout and poll interval. E.g. if timeout is 5 minutes and poll interval is 15 seconds,
			// the maximum number of retries will be 20 (5 minutes / 15 seconds = 20).
			// This ensures that the job status is checked periodically until it is completed.
			SetCustomRestyOption(func(r *resty.Request) { r.SetRetryCount(int(jobOpts.Timeout / jobOpts.PollInterval)) }),
			SetCustomRestyOption(func(r *resty.Request) { r.SetTimeout(jobOpts.Timeout) }),
		)

		// Use the subclient's JobRefresh method to refresh the job status.
		// This method will handle the job response and return the updated job status.
		job, err := c.JobRefresh(httpC, resp, reqOpts)

		if job != nil {
			xlogger.Debug("Job completed", slog.String("jobID", job.ID), slog.String("status", job.Status.String()))
		}

		// If an error occurs while refreshing the job status, return an err.
		// In the subclient the error will be handled and returned as a custom error type (e.g., errors.APIError)
		// if the error is related by the job or a generic error if the job is not related.
		return err
	}
}

var jobRetryCondition = func(c jobsInterface) resty.RetryConditionFunc {
	return func(r *resty.Response, err error) bool {
		xlogger.Debug("Checking job status", slog.Int("attempt", r.Request.Attempt), slog.Duration("retryWaitTime", r.Request.RetryWaitTime))

		if err != nil {
			xlogger.Error("Error occurred while waiting for job response", slog.String("error", err.Error()))
			return false // Stop retrying if an error occurs
		}

		job, err := c.JobParser(r)
		if err != nil {
			xlogger.Error("Failed to parse job response", slog.String("error", err.Error()))
			return false // Stop retrying if parsing fails
		}

		if job == nil {
			xlogger.Warn("Job response is nil, stopping retries")
			return false // Stop retrying if the job response is nil
		}

		xlogger.Debug("Job response status", slog.String("status", job.Status.String()))
		return !job.Status.IsTerminated() // Continue retrying if the job is not terminated
	}
}

var extractorFuncMiddleware = func(extractorFunc func(resp *resty.Response)) resty.ResponseMiddleware {
	return func(_ *resty.Client, resp *resty.Response) error {
		if extractorFunc != nil {
			// Use the extractor function to extract data from the response.
			extractorFunc(resp)
		}
		return nil
	}
}
