/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL-2.0/
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

func TestAwaitJobValidation(t *testing.T) {
	client, err := newMockClient()
	require.NoError(t, err)
	defer func() { require.NoError(t, client.Close()) }()

	_, err = AwaitJob[string](t.Context(), client, "", JobPollOptions{
		Timeout:         time.Millisecond,
		PollingInterval: time.Millisecond,
	}, func(_ *resty.Response) (string, error) { return "", nil })
	require.EqualError(t, err, "await job: job id is required")
}

func TestAwaitJobContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancelCause(t.Context())
	cancel(fmt.Errorf("caller canceled"))

	_, err := AwaitJob[string](ctx, nil, "id", JobPollOptions{
		Timeout:         time.Second,
		PollingInterval: time.Millisecond,
	}, func(_ *resty.Response) (string, error) { return "", nil })
	require.EqualError(t, err, "caller canceled")
}

func TestAwaitJobNilClient(t *testing.T) {
	_, err := AwaitJob[string](t.Context(), nil, "id", JobPollOptions{
		Timeout:         time.Second,
		PollingInterval: time.Millisecond,
	}, func(_ *resty.Response) (string, error) { return "", nil })
	require.EqualError(t, err, "await job id: client is nil")
}

type awaitJobTestClient struct {
	polls atomic.Int32
}

func (c *awaitJobTestClient) NewRequest(_ context.Context, _ *Endpoint, _ ...RequestOption) (*resty.Request, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *awaitJobTestClient) NewRequestWithBackend(_ context.Context, _ BackendTarget, _ *Endpoint, _ ...RequestOption) (*resty.Request, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *awaitJobTestClient) Logger() *slog.Logger {
	return slog.Default()
}

func (c *awaitJobTestClient) Do(_ context.Context, _ *Endpoint, _ ...EndpointRequestOption) (*resty.Response, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *awaitJobTestClient) DoWithBackend(_ context.Context, backend BackendTarget, endpoint *Endpoint, _ ...EndpointRequestOption) (*resty.Response, error) {
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

func (c *awaitJobTestClient) GetConsole() consoles.ConsoleName {
	return consoles.ConsoleName("test")
}

func (c *awaitJobTestClient) Close() error {
	return nil
}

func TestAwaitJobPollsImmediately(t *testing.T) {
	client := &awaitJobTestClient{}
	start := time.Now()

	result, err := AwaitJobOnBackend[string](t.Context(), client, BackendVMware, "job-id", JobPollOptions{
		Timeout:         time.Second,
		PollingInterval: 200 * time.Millisecond,
	}, func(_ *resty.Response) (string, error) { return "ok", nil })
	require.NoError(t, err)
	require.Equal(t, "ok", result)
	require.EqualValues(t, 1, client.polls.Load())
	require.Less(t, time.Since(start), 150*time.Millisecond)
}
