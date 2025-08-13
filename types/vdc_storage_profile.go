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
	// ParamsListStorageProfiles defines the parameters for listing storage profiles in a VDC.
	ParamsListStorageProfiles struct {
		// ID is the unique identifier of the VDC to get storage profiles from.
		ID string `documentation:"ID of the VDC to get storage profiles from"`
	}

	ParamsAddStorageProfile struct {
		// VdcId is the unique identifier of the VDC to add the storage profile to.
		VdcID string
		// VdcName is the name of the VDC to add the storage profile to.
		VdcName string

		StorageProfiles []ParamsCreateVDCStorageProfile
	}
)
