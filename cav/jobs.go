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

const (
	defaultJobWaitTimeout         = 3 * time.Minute
	defaultJobWaitPollingInterval = 2 * time.Second
)

// JobWaitRequest describes a known backend job to wait for.
type JobWaitRequest struct {
	// Backend hosting the job. Must be BackendInfrapi (Cerberus) or BackendVMware.
	Backend BackendTarget

	// JobID is the backend job identifier.
	JobID string

	// Timeout is the maximum total wait time. When zero, a backend-oriented
	// default is used.
	Timeout time.Duration

	// PollingInterval is the delay between status checks. When zero, a default
	// interval is used.
	PollingInterval time.Duration

	// Jitter randomizes PollingInterval by +/- Jitter. Optional.
	Jitter time.Duration
}

// WaitJob polls a known backend job until it reaches a terminal state and
// returns that terminal Job.
//
// The job is returned with a nil error only when it succeeded. When the job
// ends in error or is aborted, the terminal Job is returned together with an
// error wrapping pkgerrors.ErrJobFailed. If the timeout elapses the error is
// pkgerrors.ErrJobTimeout; if ctx is canceled its cause is returned.
func WaitJob(ctx context.Context, c Client, req JobWaitRequest) (*Job, error) {
	if err := context.Cause(ctx); err != nil {
		return nil, err
	}

	if c == nil {
		return nil, fmt.Errorf("await job %s: client is nil", req.JobID)
	}
	if req.JobID == "" {
		return nil, fmt.Errorf("await job: job id is required")
	}
	if req.Backend == 0 {
		return nil, fmt.Errorf("await job %s: backend is required", req.JobID)
	}
	if req.Timeout == 0 {
		req.Timeout = defaultJobWaitTimeout
	}
	if req.Timeout < 0 {
		return nil, fmt.Errorf("await job %s: timeout must be greater than 0", req.JobID)
	}
	if req.PollingInterval == 0 {
		req.PollingInterval = defaultJobWaitPollingInterval
	}
	if req.PollingInterval < 0 {
		return nil, fmt.Errorf("await job %s: polling interval must be greater than 0", req.JobID)
	}

	pollCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	timer := time.AfterFunc(req.Timeout, func() { cancel(pkgerrors.ErrJobTimeout) })
	defer timer.Stop()
	startedAt := time.Now()

	for first := true; ; first = false {
		if !first {
			select {
			case <-pollCtx.Done():
				return nil, context.Cause(pollCtx)
			case <-time.After(withJitter(req.PollingInterval, req.Jitter)):
			}
		}

		job, err := fetchJob(pollCtx, c, req.Backend, req.JobID)
		if err != nil {
			return job, fmt.Errorf("await job %s: %w", req.JobID, err)
		}

		c.Logger().Debug("job poll", "job_id", req.JobID, "status", job.Status, "elapsed", time.Since(startedAt))

		switch job.Status {
		case JobSuccess:
			return job, nil
		case JobError, JobAborted:
			return job, fmt.Errorf("job %s: %w", req.JobID, pkgerrors.ErrJobFailed)
		}
	}
}

// ParseCerberusJob parses a Cerberus job from the given response.
func ParseCerberusJob(resp *resty.Response) (*Job, error) {
	return parseCerberusJob(resp)
}

// ParseVMwareJob parses a VMware job from the given response.
func ParseVMwareJob(resp *resty.Response) (*Job, error) {
	return parseVMwareJob(resp)
}
