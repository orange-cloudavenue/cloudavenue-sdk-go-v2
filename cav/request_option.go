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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/jobs"
)

type (
	// requestOption is a function that modifies the request.
	requestOption struct {
		// Job indicates if the request is for a job.
		JobOpts *jobs.JobOptions
	}

	RequestOption func(*requestOption) error
)

// Create a request option struct to hold the request options.
// This struct will be populated with the options provided in reqOpt.
func newRequestOptions(opts ...RequestOption) (*requestOption, error) {
	// Create a new request option struct to hold the options.
	ro := &requestOption{}
	for _, opt := range opts {
		if err := opt(ro); err != nil {
			return nil, err
		}
	}
	return ro, nil
}

// * Job

// WithJob is a request option to parse the Job Response.
func WithJob(opts ...jobs.JobOption) RequestOption {
	return func(ro *requestOption) error {
		// This option is used to parse the job response.
		// It can be used to set the job settings or any other job-related options.
		jobOpts, err := jobs.NewJobOptions(opts...)
		if err != nil {
			return err
		}
		ro.JobOpts = jobOpts

		return nil
	}
}
