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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIResponseGetVDCStorageProfiles_ArrayOrObject(t *testing.T) {
	const profile = `{"id": "urn:vcloud:vdcstorageProfile:33333333-3333-4333-8333-333333333333", "name": "gold", "limit": 500}`

	t.Run("plain array (real API)", func(t *testing.T) {
		var r APIResponseGetVDCStorageProfiles
		require.NoError(t, json.Unmarshal([]byte(`[`+profile+`]`), &r))
		require.Len(t, r.StorageProfiles, 1)
		assert.Equal(t, "gold", r.StorageProfiles[0].Name)
	})

	t.Run("wrapped object", func(t *testing.T) {
		var r APIResponseGetVDCStorageProfiles
		require.NoError(t, json.Unmarshal([]byte(`{"vdcStorageProfile": [`+profile+`]}`), &r))
		require.Len(t, r.StorageProfiles, 1)
		assert.Equal(t, "gold", r.StorageProfiles[0].Name)
	})

	t.Run("null", func(t *testing.T) {
		var r APIResponseGetVDCStorageProfiles
		require.NoError(t, json.Unmarshal([]byte(`null`), &r))
		assert.Empty(t, r.StorageProfiles)
	})
}
