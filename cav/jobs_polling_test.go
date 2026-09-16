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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

func TestWithJitter(t *testing.T) {
	interval := 10 * time.Millisecond
	got := withJitter(interval, 0)
	require.Equal(t, interval, got)
}

func TestParseJobResponseUnsupportedClient(t *testing.T) {
	_, err := parseJobResponse(&Response{Raw: &resty.Response{}}, BackendNetBackup)
	require.EqualError(t, err, "backend 4 does not support jobs")
}

func TestGetJobEndpointNameUnsupportedBackend(t *testing.T) {
	_, err := getJobEndpointName(BackendNetBackup)
	require.EqualError(t, err, "backend 4 does not support jobs")
}
