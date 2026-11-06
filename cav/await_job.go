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

// JobAwaitRequest describes how to wait for a backend job and optionally
// extract a typed result from the terminal job response.
type JobAwaitRequest[R any] struct {
	Backend     BackendTarget
	JobID       string
	PollOptions JobPollOptions
	Extract     func(resp *resty.Response) (R, error)
}

// JobWaitRequest describes how to wait for a backend job when no response
// extraction is needed by the caller.
type JobWaitRequest struct {
	Backend     BackendTarget
	JobID       string
	PollOptions JobPollOptions
}

// AwaitJobOnBackend polls a backend-specific job until completion and extracts
// a typed result from the terminal response.
func AwaitJobOnBackend[R any](ctx context.Context, c Client, backend BackendTarget, jobID string, opts JobPollOptions, extract func(resp *resty.Response) (R, error)) (R, error) {
	return awaitJob(ctx, c, JobAwaitRequest[R]{
		Backend:     backend,
		JobID:       jobID,
		PollOptions: opts,
		Extract:     extract,
	})
}

// AwaitJob polls a job until completion, probing the backend automatically on
// the first request, then extracts a typed result from the terminal response.
func AwaitJob[R any](ctx context.Context, c Client, jobID string, opts JobPollOptions, extract func(resp *resty.Response) (R, error)) (R, error) {
	return awaitJob(ctx, c, JobAwaitRequest[R]{
		Backend:     0,
		JobID:       jobID,
		PollOptions: opts,
		Extract:     extract,
	})
}

// WaitJob polls a backend-specific job until completion and ignores the
// terminal response payload.
func WaitJob(ctx context.Context, c Client, req JobWaitRequest) error {
	_, err := awaitJob(ctx, c, JobAwaitRequest[struct{}]{
		Backend:     req.Backend,
		JobID:       req.JobID,
		PollOptions: req.PollOptions,
	})
	return err
}

func awaitJob[R any](ctx context.Context, c Client, req JobAwaitRequest[R]) (R, error) {
	var zero R

	if err := context.Cause(ctx); err != nil {
		return zero, err
	}

	if c == nil {
		return zero, fmt.Errorf("await job %s: client is nil", req.JobID)
	}
	if req.JobID == "" {
		return zero, fmt.Errorf("await job: job id is required")
	}
	if req.PollOptions.Timeout <= 0 {
		return zero, fmt.Errorf("await job %s: timeout must be greater than 0", req.JobID)
	}
	if req.PollOptions.PollingInterval <= 0 {
		return zero, fmt.Errorf("await job %s: polling interval must be greater than 0", req.JobID)
	}

	pollCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	timer := time.AfterFunc(req.PollOptions.Timeout, func() { cancel(pkgerrors.ErrJobTimeout) })
	defer timer.Stop()
	startedAt := time.Now()
	firstPoll := true

	for {
		if !firstPoll {
			select {
			case <-pollCtx.Done():
				return zero, context.Cause(pollCtx)
			case <-time.After(withJitter(req.PollOptions.PollingInterval, req.PollOptions.Jitter)):
			}
		}
		firstPoll = false

		resp, resolvedBackend, err := getJobStatus(pollCtx, c, req.JobID, req.Backend)
		if err != nil {
			return zero, fmt.Errorf("await job %s: %w", req.JobID, err)
		}

		job, err := parseJobResponse(&Response{Raw: resp}, resolvedBackend)
		if err != nil {
			return zero, fmt.Errorf("await job %s: parse job status: %w", req.JobID, err)
		}

		c.Logger().Debug("job poll", "job_id", req.JobID, "status", job.Status, "elapsed", time.Since(startedAt))

		switch job.Status {
		case JobSuccess:
			if req.Extract == nil {
				return zero, nil
			}
			result, extractErr := req.Extract(resp)
			if extractErr != nil {
				return zero, extractErr
			}
			return result, nil
		case JobError, JobAborted:
			return zero, fmt.Errorf("job %s: %w", req.JobID, pkgerrors.ErrJobFailed)
		}
	}
}
