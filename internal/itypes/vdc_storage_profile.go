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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

type (
	// * ListStorageProfiles
	APIResponseListStorageProfiles struct {
		StorageProfiles []APIResponseListStorageProfile `json:"values" fakesize:"1"`
	}

	APIResponseListStorageProfile struct {
		HREF                    string              `json:"href" fake:"{href_uuid}"`
		ID                      string              `json:"id,omitempty" fake:"{urn:vdcstorageProfile}"`
		Name                    string              `json:"name" fake:"platinum3k_r1"`
		IsEnabled               bool                `json:"isEnabled" fake:"true"`
		IsDefaultStorageProfile bool                `json:"isDefaultStoragePolicy" fake:"true"`
		OrgVDC                  *APIObjectReference `json:"orgVdcRef,omitempty"`

		// Values are in MB
		Limit int `json:"storageLimitMb" fake:"{number:100000,81920000}"` //nolint:tagliatelle
		Used  int `json:"storageUsedMB" fake:"{number:1000,100000}"`      //nolint:tagliatelle

		// VDC information
		VDCID   string `json:"vdc" fake:"{href_uuid}"`
		VDCName string `json:"vdcName" fake:"{word}"`
	}
)

// UnmarshalJSON accepts CloudAPI's storage policy fields and the former
// storage profile field names used by older fixtures.
func (r *APIResponseListStorageProfile) UnmarshalJSON(data []byte) error {
	var response struct {
		HREF                    string              `json:"href"`
		ID                      string              `json:"id"`
		Name                    string              `json:"name"`
		IsEnabled               bool                `json:"isEnabled"`
		IsDefaultStoragePolicy  bool                `json:"isDefaultStoragePolicy"`
		IsDefaultStorageProfile bool                `json:"isDefaultStorageProfile"`
		OrgVDC                  *APIObjectReference `json:"orgVdcRef"`
		LegacyOrgVDC            *APIObjectReference `json:"orgVdc"`
		Limit                   int                 `json:"storageLimitMb"`
		LegacyLimit             int                 `json:"storageLimitMB"`
		Used                    int                 `json:"storageUsedMB"`
		VDCID                   string              `json:"vdc"`
		VDCName                 string              `json:"vdcName"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	r.HREF = response.HREF
	r.ID = response.ID
	r.Name = response.Name
	r.IsEnabled = response.IsEnabled
	r.IsDefaultStorageProfile = response.IsDefaultStoragePolicy || response.IsDefaultStorageProfile
	r.OrgVDC = response.OrgVDC
	if r.OrgVDC == nil {
		r.OrgVDC = response.LegacyOrgVDC
	}
	r.Limit = response.Limit
	if r.Limit == 0 {
		r.Limit = response.LegacyLimit
	}
	r.Used = response.Used
	r.VDCID = response.VDCID
	r.VDCName = response.VDCName
	return nil
}

// UnmarshalJSON accepts CloudAPI's values collection and the former query
// endpoint's record collection. The latter keeps response decoding compatible
// for callers replaying older fixtures.
func (r *APIResponseListStorageProfiles) UnmarshalJSON(data []byte) error {
	var response struct {
		Values []APIResponseListStorageProfile `json:"values"`
		Record []APIResponseListStorageProfile `json:"record"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	r.StorageProfiles = response.Values
	if len(r.StorageProfiles) == 0 {
		r.StorageProfiles = response.Record
	}
	return nil
}

func (r *APIResponseListStorageProfiles) ToModel() *types.ModelListStorageProfiles {
	// Use a map to group storage profiles by unique VDC ID + Name
	type ModelVDCKey struct {
		ID, Name string
	}
	vdcMap := make(map[ModelVDCKey]*types.ModelListStorageProfilesVDC)
	for _, apiSP := range r.StorageProfiles {
		vdcID, vdcName := apiSP.VDCID, apiSP.VDCName
		if apiSP.OrgVDC != nil {
			if vdcID == "" {
				vdcID = apiSP.OrgVDC.ID
			}
			if vdcName == "" {
				vdcName = apiSP.OrgVDC.Name
			}
		}
		key := ModelVDCKey{ID: vdcID, Name: vdcName}
		vdc, exists := vdcMap[key]
		if !exists {
			vdc = &types.ModelListStorageProfilesVDC{
				ID:              vdcID,
				Name:            vdcName,
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
