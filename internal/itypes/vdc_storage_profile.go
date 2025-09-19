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
		HREF                    string `json:"href" fake:"{href_uuid}"`
		ID                      string `json:"id" fake:"-"` // Because VMware returns an empty ID, we will extract it from the HREF
		Name                    string `json:"name" fake:"platinum3k_r1"`
		IsEnabled               bool   `json:"isEnabled" fake:"true"`
		IsDefaultStorageProfile bool   `json:"isDefaultStorageProfile" fake:"true"`

		// Values are in MB
		Limit int `json:"storageLimitMB" fake:"{number:100000,81920000}"`
		Used  int `json:"storageUsedMB" fake:"{number:1000,100000}"`

		// Vdc information
		VdcID   string `json:"vdc" fake:"{href_uuid}"`
		VdcName string `json:"vdcName" fake:"{word}"`
	}
)

func (r *ApiResponseListStorageProfiles) ToModel() *types.ModelListStorageProfiles {
	// Use a map to group storage profiles by unique VDC ID + Name
	type ModelVDCKey struct {
		ID, Name string
	}
	vdcMap := make(map[ModelVDCKey]*types.ModelListStorageProfilesVDC)
	for _, apiSP := range r.StorageProfiles {
		key := ModelVDCKey{ID: apiSP.VdcID, Name: apiSP.VdcName}
		vdc, exists := vdcMap[key]
		if !exists {
			vdc = &types.ModelListStorageProfilesVDC{
				ID:              apiSP.VdcID,
				Name:            apiSP.VdcName,
				StorageProfiles: []types.ModelListStorageProfile{},
			}
			vdcMap[key] = vdc
		}
		vdc.StorageProfiles = append(vdc.StorageProfiles, types.ModelListStorageProfile{
			ID:      apiSP.ID,
			Class:   apiSP.Name,
			Limit:   apiSP.Limit,
			Used:    apiSP.Used,
			Default: apiSP.IsDefaultStorageProfile,
		})
	}

	// Convert map to slice
	vdcs := make([]types.ModelListStorageProfilesVDC, 0, len(vdcMap))
	for _, vdc := range vdcMap {
		vdcs = append(vdcs, *vdc)
	}

	return &types.ModelListStorageProfiles{
		VDCS: vdcs,
	}
}
