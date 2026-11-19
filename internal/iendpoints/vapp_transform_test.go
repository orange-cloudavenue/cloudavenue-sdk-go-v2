/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToLegacyVAppID(t *testing.T) {
	const (
		vapp   = "urn:vcloud:vapp:2feab1dc-1b61-4d4d-8f75-fab1f9a2acee"
		uuid   = "2feab1dc-1b61-4d4d-8f75-fab1f9a2acee"
		legacy = "vapp-2feab1dc-1b61-4d4d-8f75-fab1f9a2acee"
	)

	tests := []struct {
		name    string
		value   string
		want    string
		wantErr bool
	}{
		{name: "vapp urn", value: vapp, want: legacy},
		{name: "bare uuid", value: uuid, want: legacy},
		{name: "already prefixed", value: legacy, want: legacy},
		{name: "vdc urn rejected", value: "urn:vcloud:vdc:a76aec87-fa19-4ece-aab9-2b24138863d6", wantErr: true},
		{name: "arbitrary value rejected", value: "not-an-id", wantErr: true},
		{name: "empty value rejected", value: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toLegacyVAppID(tt.value)
			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
