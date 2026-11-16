/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
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
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

func TestWaitJobValidation(t *testing.T) {
	client, err := newMockClient()
	require.NoError(t, err)
	defer func() { require.NoError(t, client.Close()) }()

	_, err = WaitJob(t.Context(), client, JobWaitRequest{
		Backend:         BackendVMware,
		Timeout:         time.Millisecond,
		PollingInterval: time.Millisecond,
	})
	require.EqualError(t, err, "await job: job id is required")
}

func TestWaitJobContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancelCause(t.Context())
	cancel(fmt.Errorf("caller canceled"))

	_, err := WaitJob(ctx, nil, JobWaitRequest{
		Backend:         BackendVMware,
		JobID:           "id",
		Timeout:         time.Second,
		PollingInterval: time.Millisecond,
	})
	require.EqualError(t, err, "caller canceled")
}

func TestWaitJobNilClient(t *testing.T) {
	_, err := WaitJob(t.Context(), nil, JobWaitRequest{
		Backend:         BackendVMware,
		JobID:           "id",
		Timeout:         time.Second,
		PollingInterval: time.Millisecond,
	})
	require.EqualError(t, err, "await job id: client is nil")
}

type waitJobTestClient struct {
	polls atomic.Int32
}

func (c *waitJobTestClient) NewRequest(_ context.Context, _ *Endpoint, _ ...RequestOption) (*resty.Request, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *waitJobTestClient) NewRequestWithBackend(_ context.Context, _ BackendTarget, _ *Endpoint, _ ...RequestOption) (*resty.Request, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *waitJobTestClient) Logger() *slog.Logger {
	return slog.Default()
}

func (c *waitJobTestClient) Do(_ context.Context, _ *Endpoint, _ ...EndpointRequestOption) (*resty.Response, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *waitJobTestClient) DoWithBackend(_ context.Context, backend BackendTarget, endpoint *Endpoint, _ ...EndpointRequestOption) (*resty.Response, error) {
	c.polls.Add(1)

	req := &resty.Request{}
	req.SetResult(&vmwareJobAPIResponse{Status: string(JobSuccess)})
	resp := &resty.Response{
		Request: req,
	}

	if backend != BackendVMware && endpoint.Backend != BackendVMware {
		return nil, fmt.Errorf("unexpected backend %d", backend)
	}

	return resp, nil
}

func (c *waitJobTestClient) GetConsole() consoles.ConsoleName {
	return consoles.ConsoleName("test")
}

func (c *waitJobTestClient) Close() error {
	return nil
}

func TestWaitJobPollsImmediately(t *testing.T) {
	client := &waitJobTestClient{}
	start := time.Now()

	job, err := WaitJob(t.Context(), client, JobWaitRequest{
		Backend:         BackendVMware,
		JobID:           "job-id",
		Timeout:         time.Second,
		PollingInterval: 200 * time.Millisecond,
	})
	require.NoError(t, err)
	require.Equal(t, JobSuccess, job.Status)
	require.EqualValues(t, 1, client.polls.Load())
	require.Less(t, time.Since(start), 150*time.Millisecond)
}

func TestWaitJobBackendRequired(t *testing.T) {
	_, err := WaitJob(t.Context(), &waitJobTestClient{}, JobWaitRequest{
		JobID:           "job-id",
		Timeout:         time.Second,
		PollingInterval: time.Millisecond,
	})
	require.EqualError(t, err, "await job job-id: backend is required")
}

func TestApplyJobWaitDefaults(t *testing.T) {
	t.Run("explicit request values win over client config", func(t *testing.T) {
		cc := &client{jobTimeout: 5 * time.Minute, jobPollInterval: 3 * time.Second}
		req := applyJobWaitDefaults(cc, JobWaitRequest{
			Timeout:         time.Second,
			PollingInterval: time.Second,
		})
		require.Equal(t, time.Second, req.Timeout)
		require.Equal(t, time.Second, req.PollingInterval)
	})

	t.Run("client config fills values the request leaves unset", func(t *testing.T) {
		cc := &client{jobTimeout: 5 * time.Minute, jobPollInterval: 3 * time.Second}
		req := applyJobWaitDefaults(cc, JobWaitRequest{Timeout: 2 * time.Second})
		require.Equal(t, 2*time.Second, req.Timeout)
		require.Equal(t, 3*time.Second, req.PollingInterval)
	})

	t.Run("request values used when client is not configured", func(t *testing.T) {
		req := applyJobWaitDefaults(&client{}, JobWaitRequest{
			Timeout:         2 * time.Second,
			PollingInterval: 50 * time.Millisecond,
		})
		require.Equal(t, 2*time.Second, req.Timeout)
		require.Equal(t, 50*time.Millisecond, req.PollingInterval)
	})

	t.Run("package defaults when neither is set", func(t *testing.T) {
		req := applyJobWaitDefaults(&waitJobTestClient{}, JobWaitRequest{})
		require.Equal(t, defaultJobWaitTimeout, req.Timeout)
		require.Equal(t, defaultJobWaitPollingInterval, req.PollingInterval)
	})
}

func TestWithJobWaitClientOptions(t *testing.T) {
	t.Run("rejects non-positive values", func(t *testing.T) {
		require.Error(t, WithJobTimeout(0)(&settings{}))
		require.Error(t, WithJobTimeout(-time.Second)(&settings{}))
		require.Error(t, WithJobPollInterval(0)(&settings{}))
		require.Error(t, WithJobPollInterval(-time.Second)(&settings{}))
	})

	t.Run("stores configured values", func(t *testing.T) {
		s := &settings{}
		require.NoError(t, WithJobTimeout(4*time.Minute)(s))
		require.NoError(t, WithJobPollInterval(500*time.Millisecond)(s))
		require.Equal(t, 4*time.Minute, s.JobTimeout)
		require.Equal(t, 500*time.Millisecond, s.JobPollInterval)
	})
}
