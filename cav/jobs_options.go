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
	"time"

	"github.com/orange-cloudavenue/common-go/validators"
	"resty.dev/v3"
)

type (
	// JobOptions configures legacy job polling helpers.
	JobOptions struct {
		// Timeout is maximum wait time for job completion.
		Timeout time.Duration `default:"5m"`

		// PollInterval is delay between status checks.
		PollInterval time.Duration `default:"15s"`

		extractorFunc ExtractorFunc
	}

	// ExtractorFunc extracts side-channel data from job responses.
	// It does not return an error so polling flow is not interrupted.
	ExtractorFunc func(resp *resty.Response)

	// JobOption applies configuration to JobOptions.
	JobOption func(*JobOptions) error
)

// NewJobOptions creates JobOptions with defaults and applies opts.
func NewJobOptions(opts ...JobOption) (*JobOptions, error) {
	jO := &JobOptions{}

	if err := validators.New().Struct(jO); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(jO); err != nil {
			return nil, err
		}
	}

	return jO, nil
}

// WithCustomTimeout sets job timeout.
func WithCustomTimeout(timeout time.Duration) JobOption {
	return func(opts *JobOptions) error {
		opts.Timeout = timeout
		return nil
	}
}

// WithCustomPollInterval sets job poll interval.
func WithCustomPollInterval(interval time.Duration) JobOption {
	return func(opts *JobOptions) error {
		if interval <= 0 {
			return fmt.Errorf("poll interval must be greater than 0, got %s", interval)
		}
		opts.PollInterval = interval
		return nil
	}
}

// SetExtractorFunc sets response extractor used during job polling.
func SetExtractorFunc(extractorFunc ExtractorFunc) JobOption {
	return func(opts *JobOptions) error {
		opts.extractorFunc = extractorFunc
		return nil
	}
}
