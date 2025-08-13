/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

type (
	// * ListStorageProfiles
	ApiResponseListStorageProfiles struct {
		StorageProfiles []ApiResponseListStorageProfile `json:"record" fakesize:"1"`
	}

	ApiResponseListStorageProfile struct {
		HREF      string `json:"href" fake:"{href_uuid}"`
		ID        string `json:"id" fake:"{urn:vdcstorageProfile}"`
		Name      string `json:"name" fake:"platinum3k_r1"`
		IsEnabled bool   `json:"isEnabled" fake:"true"`
		Default   bool   `json:"default" fake:"true"`

		// Values are in MB
		Limit int `json:"storageLimitMB" fake:"{number:100,1000}"` //nolint:tagliatelle
		Used  int `json:"storageUsedMB" fake:"{number:10,500}"`    //nolint:tagliatelle
	}
)

func (r *ApiResponseListStorageProfiles) ToModel() *types.ModelListStorageProfiles {
	storageProfiles := make([]types.ModelListStorageProfile, 0, len(r.StorageProfiles))
	for _, sp := range r.StorageProfiles {
		storageProfiles = append(storageProfiles, types.ModelListStorageProfile{
			ID:      sp.ID,
			Class:   sp.Name,
			Limit:   sp.Limit,
			Used:    sp.Used,
			Default: sp.Default,
		})
	}
	return &types.ModelListStorageProfiles{StorageProfiles: storageProfiles}
}
