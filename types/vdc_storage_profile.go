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
		// Name is the name of the storage profiles to filter by.
		Name string `documentation:"Name of the storage profiles to filter by"`
		// VDCID is the unique identifier of the VDC to get storage profiles from.
		VdcID string `documentation:"ID of the VDC to get storage profiles from"` //nolint:revive
		// Name is the name of the VDC to get storage profiles from.
		VdcName string `documentation:"Name of the VDC to get storage profiles from"` //nolint:revive
	}

	ParamsAddStorageProfile struct {
		// VdcId is the unique identifier of the VDC to add the storage profile to.
		VdcID string //nolint:revive
		// VdcName is the name of the VDC to add the storage profile to.
		VdcName string //nolint:revive

		StorageProfiles []ParamsCreateVDCStorageProfile
	}

	ParamsDeleteStorageProfile struct {
		// VdcId is the unique identifier of the VDC to delete the storage profile
		VdcID string //nolint:revive
		// VdcName is the name of the VDC to delete the storage profile from
		VdcName string //nolint:revive
		// StorageProfile is the list of storage profiles to delete, in this case it is a single storage profile
		StorageProfiles []ParamsDeleteVDCStorageProfile
	}
)
