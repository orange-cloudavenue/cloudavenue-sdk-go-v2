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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"resty.dev/v3"

	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
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

func TestWithJitter(t *testing.T) {
	interval := 10 * time.Millisecond
	got := withJitter(interval, 0)
	require.Equal(t, interval, got)
}

func TestParseJobResponseUnsupportedClient(t *testing.T) {
	_, err := parseJobResponse(&Response{Raw: &resty.Response{}}, BackendNetBackup)
	require.EqualError(t, err, "backend 4 does not support jobs")
}

func TestClassifyStatusCode(t *testing.T) {
	require.ErrorIs(t, classifyStatusCode(404), pkgerrors.ErrNotFound)
	require.NoError(t, classifyStatusCode(500))
}
