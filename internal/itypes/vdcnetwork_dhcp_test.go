/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDhcpConfigFalseValuesAreSerialized(t *testing.T) {
	payload, err := json.Marshal(DhcpConfig{
		Enabled: false,
		DhcpPools: []DhcpPool{{
			Enabled: false,
		}},
	})
	require.NoError(t, err)

	assertJSONField(t, payload, "enabled", false)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(payload, &decoded))
	pools, ok := decoded["dhcpPools"].([]any)
	require.True(t, ok)
	require.Len(t, pools, 1)
	pool, ok := pools[0].(map[string]any)
	require.True(t, ok)
	require.Contains(t, pool, "enabled")
	require.Equal(t, false, pool["enabled"])
}

func assertJSONField(t *testing.T, payload []byte, name string, want any) {
	t.Helper()

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(payload, &decoded))
	require.Contains(t, decoded, name)
	require.Equal(t, want, decoded[name])
}
