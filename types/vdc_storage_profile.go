/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

type (
	ModelListStorageProfiles struct {
		VDCS []ModelListStorageProfilesVDC `documentation:"List of VDCs to list storage profiles for"`
	}

	ModelListStorageProfilesVDC struct {
		ID              string                    `documentation:"ID of the VDC"`
		Name            string                    `documentation:"Name of the VDC"`
		StorageProfiles []ModelListStorageProfile `documentation:"List of storage profiles in the VDC"`
	}

	ModelListStorageProfile struct {
		ID      string `documentation:"ID of the storage profile"`
		Class   string `documentation:"Name of the storage profile"`
		Limit   int    `documentation:"Limit of the storage profile in MB"`
		Used    int    `documentation:"Used storage of the storage profile in MB"`
		Default bool   `documentation:"Is this storage profile the default one?"`
	}
)

type (
	// ParamsListStorageProfile defines the parameters for listing storage profiles in a VDC.
	ParamsListStorageProfile struct {
		// ID is the unique identifier of storage profiles from.
		ID string `documentation:"ID of the storage profile"`
		// Name is the name of the VDC to get storage profiles from.
		Name string `documentation:"Name of the storage profiles from"`
		// VDCID is the unique identifier of the VDC to get storage profiles from.
		VdcID string `documentation:"ID of the VDC to get storage profiles from"` //nolint:revive
		// Name is the name of the VDC to get storage profiles from.
		VdcName string `documentation:"Name of the VDC to get storage profiles from"` //nolint:revive
	}

	ParamsAddStorageProfile struct {
		// VdcId is the unique identifier of the VDC to add the storage profile to.
		VdcId string //nolint:revive
		// VdcName is the name of the VDC to add the storage profile to.
		VdcName string //nolint:revive

		StorageProfiles []ParamsCreateVDCStorageProfile
	}

	ParamsDeleteStorageProfile struct {
		// VdcId is the unique identifier of the VDC to delete the storage profile
		VdcId string //nolint:revive
		// VdcName is the name of the VDC to delete the storage profile from
		VdcName string //nolint:revive
		// StorageProfile is the list of storage profiles to delete, is this case it is a single storage profile
		StorageProfile []ParamsCreateVDCStorageProfile
	}
)

type (
	// * ListStorageProfiles
	apiResponseListStorageProfiles struct {
		StorageProfiles []apiResponseListStorageProfile `json:"record" fakesize:"1"`
	}

	apiResponseListStorageProfile struct {
		HREF                    string `json:"href" fake:"{href_uuid}"`
		ID                      string `json:"id" fake:"-"` // Because VMware returns an empty ID, we will extract it from the HREF
		Name                    string `json:"name" fake:"platinum3k_r1"`
		IsEnabled               bool   `json:"isEnabled" fake:"true"`
		IsDefaultStorageProfile bool   `json:"isDefaultStorageProfile" fake:"true"`

		// Values are in MB
		Limit int `json:"storageLimitMB" fake:"{number:100,1000}"` //nolint:tagliatelle
		Used  int `json:"storageUsedMB" fake:"{number:10,500}"`    //nolint:tagliatelle

		// Vdc information
		VdcId   string `json:"vdc" fake:"{urn:vdc}"`  //nolint:revive
		VdcName string `json:"vdcName" fake:"{word}"` //nolint:revive
	}
)

func (r *apiResponseListStorageProfiles) toModel() *ModelListStorageProfiles {
	// Use a map to group storage profiles by unique VDC ID + Name
	type ModelVDCKey struct {
		ID, Name string
	}
	vdcMap := make(map[ModelVDCKey]*ModelListStorageProfilesVDC)
	for _, apiSP := range r.StorageProfiles {
		key := ModelVDCKey{ID: apiSP.VdcId, Name: apiSP.VdcName}
		vdc, exists := vdcMap[key]
		if !exists {
			vdc = &ModelListStorageProfilesVDC{
				ID:              apiSP.VdcId,
				Name:            apiSP.VdcName,
				StorageProfiles: []ModelListStorageProfile{},
			}
			vdcMap[key] = vdc
		}
		vdc.StorageProfiles = append(vdc.StorageProfiles, ModelListStorageProfile{
			ID:      apiSP.ID,
			Class:   apiSP.Name,
			Limit:   apiSP.Limit,
			Used:    apiSP.Used,
			Default: apiSP.IsDefaultStorageProfile,
		})
	}

	// Convert map to slice
	vdcs := make([]ModelListStorageProfilesVDC, 0, len(vdcMap))
	for _, vdc := range vdcMap {
		vdcs = append(vdcs, *vdc)
	}

	return &ModelListStorageProfiles{
		VDCS: vdcs,
	}
}
