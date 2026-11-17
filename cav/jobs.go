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

	// Timeout is the maximum total wait time. An explicit value wins over a
	// client-level WithJobTimeout and the package default; when zero, the
	// client-level value applies, falling back to the package default.
	Timeout time.Duration

	// PollingInterval is the delay between status checks. Resolution order
	// mirrors Timeout: explicit per-request value, then client-level
	// WithJobPollInterval, then the package default.
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
	if req.Timeout < 0 {
		return nil, fmt.Errorf("await job %s: timeout must be greater than 0", req.JobID)
	}
	if req.PollingInterval < 0 {
		return nil, fmt.Errorf("await job %s: polling interval must be greater than 0", req.JobID)
	}

	req = applyJobWaitDefaults(c, req)

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

// applyJobWaitDefaults resolves the effective timeout and polling interval for
// a job wait. Priority: an explicit per-request value wins; otherwise a
// client-level value configured with WithJobTimeout / WithJobPollInterval;
// otherwise the package defaults.
func applyJobWaitDefaults(c Client, req JobWaitRequest) JobWaitRequest {
	var cfgTimeout, cfgInterval time.Duration
	if cc, ok := c.(*client); ok {
		cfgTimeout, cfgInterval = cc.jobTimeout, cc.jobPollInterval
	}
	req.Timeout = firstPositive(req.Timeout, cfgTimeout, defaultJobWaitTimeout)
	req.PollingInterval = firstPositive(req.PollingInterval, cfgInterval, defaultJobWaitPollingInterval)
	return req
}

// firstPositive returns the first value greater than zero.
func firstPositive(values ...time.Duration) time.Duration {
	for _, v := range values {
		if v > 0 {
			return v
		}
	}
	return 0
}
