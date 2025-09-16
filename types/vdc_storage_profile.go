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
		Limit   int    `documentation:"Limit of the storage profile in GiB"`
		Used    int    `documentation:"Used storage of the storage profile in GiB"`
		Default bool   `documentation:"Indicates if the storage profile is the default one"`
	}
)

type (
	ParamsListStorageProfile struct {
		ID      string
		Name    string
		VdcID   string
		VdcName string
	}

	ParamsAddStorageProfile struct {
		VdcID           string
		VdcName         string
		StorageProfiles []ParamsCreateVDCStorageProfile

		disponibilityClass string
	}

	ParamsUpdateStorageProfile struct {
		VdcID           string
		VdcName         string
		StorageProfiles []ParamsUpdateVDCStorageProfile
	}

	ParamsDeleteStorageProfile struct {
		VdcID           string
		VdcName         string
		StorageProfiles []ParamsDeleteVDCStorageProfile
	}
)
