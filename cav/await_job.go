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
	"fmt"
	"time"

	"resty.dev/v3"

	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

// JobPollOptions defines polling behavior for AwaitJob.
type JobPollOptions struct {
	Timeout         time.Duration
	PollingInterval time.Duration
	Jitter          time.Duration
}

// AwaitJob polls until job completes or cancellation/timeout stops polling.
// The extract callback receives the raw response to allow extraction of
// resource data from the raw response (e.g. created resource names).
func AwaitJob[R any](ctx context.Context, c Client, jobID string, opts JobPollOptions, extract func(resp *resty.Response) (R, error)) (R, error) {
	var zero R

	if err := context.Cause(ctx); err != nil {
		return zero, err
	}

	if c == nil {
		return zero, fmt.Errorf("await job %s: client is nil", jobID)
	}
	if jobID == "" {
		return zero, fmt.Errorf("await job: job id is required")
	}
	if extract == nil {
		return zero, fmt.Errorf("await job %s: extract is required", jobID)
	}
	if opts.Timeout <= 0 {
		return zero, fmt.Errorf("await job %s: timeout must be greater than 0", jobID)
	}
	if opts.PollingInterval <= 0 {
		return zero, fmt.Errorf("await job %s: polling interval must be greater than 0", jobID)
	}

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	timer := time.AfterFunc(opts.Timeout, func() { cancel(pkgerrors.ErrJobTimeout) })
	defer timer.Stop()
	startedAt := time.Now()

	for {
		select {
		case <-ctx.Done():
			return zero, context.Cause(ctx)
		case <-time.After(withJitter(opts.PollingInterval, opts.Jitter)):
		}

		resp, backend, err := getJobStatus(ctx, c, jobID)
		if err != nil {
			return zero, fmt.Errorf("await job %s: %w", jobID, err)
		}

		job, err := parseJobResponse(&Response{Raw: resp}, backend)
		if err != nil {
			return zero, fmt.Errorf("await job %s: parse job status: %w", jobID, err)
		}

		c.Logger().Debug("job poll", "job_id", jobID, "status", job.Status, "elapsed", time.Since(startedAt))

		switch job.Status {
		case JobSuccess:
			result, extractErr := extract(resp)
			if extractErr != nil {
				return zero, extractErr
			}
			return result, nil
		case JobError, JobAborted:
			return zero, fmt.Errorf("job %s: %w", jobID, pkgerrors.ErrJobFailed)
		}
	}
}
