/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package jobs

import (
	"fmt"
	"time"

	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"
)

type (
	// JobOptions defines the options for job operations.
	JobOptions struct {
		// Timeout specifies the maximum duration to wait for a job to complete.
		// If the job does not complete within this time, an error will be returned.
		// Default is 5 minutes
		// If you want to wait indefinitely, set this to -1.
		Timeout time.Duration `default:"5m"`

		// PollInterval specifies the interval between job status checks.
		// Default is 15 seconds.
		// This value should be less than the timeout.
		PollInterval time.Duration `default:"15s"`

		extractorFunc ExtractorFunc
	}

	// ExtractorFunc defines a function type for extracting data from a resty.Response.
	//
	// The function does not return an error for not interrupting the job flow.
	//
	// Example usage:
	//
	//	var dataToExtract *string
	//
	//	extractorFunc := func(resp *resty.Response) {
	//		if v, ok := resp.Result().(*subclient.VmwareJobAPIResponse); ok {
	//			dataToExtract = new(string)
	//			*dataToExtract = v.ID
	//		}
	//	}
	ExtractorFunc func(resp *resty.Response)

	// JobOption defines a function that modifies the JobOptions.
	JobOption func(*JobOptions) error
)

// NewJobOptions creates a new JobOptions instance with default values.
// Default values can be overridden by passing options.
// Default values:
//   - Timeout: 5 minutes
//   - PollInterval: 15 seconds
func NewJobOptions(opts ...JobOption) (*JobOptions, error) {
	jO := &JobOptions{}

	if err := validators.New().Struct(jO); err != nil {
		return nil, err
	}

	// Override default values with provided options.
	for _, opt := range opts {
		if err := opt(jO); err != nil {
			return nil, err
		}
	}

	return jO, nil
}

// WithCustomTimeout sets the maximum duration to wait for a job to complete.
//
// Example:
//
//	opts, err := NewJobOptions(WithCustomTimeout(10 * time.Minute))
//	if err != nil {
//	    // handle error
//	}
//	// opts.Timeout will be 10 minutes
func WithCustomTimeout(timeout time.Duration) JobOption {
	return func(opts *JobOptions) error {
		opts.Timeout = timeout
		return nil
	}
}

// WithCustomPollInterval sets the interval between job status checks.
//
// Example:
//
//	opts, err := NewJobOptions(WithCustomPollInterval(30 * time.Second))
//	if err != nil {
//	    // handle error
//	}
//	// opts.PollInterval will be 30 seconds
func WithCustomPollInterval(interval time.Duration) JobOption {
	return func(opts *JobOptions) error {
		if interval <= 0 {
			return fmt.Errorf("poll interval must be greater than 0, got %s", interval)
		}
		opts.PollInterval = interval
		return nil
	}
}

// SetExtractorFunc sets a custom extractor function for parsing job responses.
// The extractor function should take a resty.Response and a target type to populate.
//
// Example:
//
//	extractor := func(resp *resty.Response) {
//	    // custom extraction logic
//	}
//	opts, err := NewJobOptions(SetExtractorFunc(extractor))
//	if err != nil {
//	    // handle error
//	}
func SetExtractorFunc(extractorFunc ExtractorFunc) JobOption {
	return func(opts *JobOptions) error {
		opts.extractorFunc = extractorFunc
		return nil
	}
}
