/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is licensed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

// ModelBMS contains BMS and network details.
type ModelBMS struct {
	Network []ModelBMSNetworkDetails `documentation:"BMS network details"`
	BMS     []ModelBMSDetails        `documentation:"BMS configuration"`
}

// ModelBMSDetails contains BMS configuration.
type ModelBMSDetails struct {
	BMSType           string          `documentation:"BMS class of service"`
	Hostname          string          `documentation:"BMS hostname"`
	OS                string          `documentation:"Operating system"`
	BIOSConfiguration string          `documentation:"BIOS configuration"`
	Storage           ModelBMSStorage `documentation:"BMS storage"`
}

// ModelBMSNetworkDetails contains BMS network details.
type ModelBMSNetworkDetails struct {
	VLANID int    `documentation:"VLAN ID"`
	Subnet string `documentation:"Subnet IP address"`
	Prefix int    `documentation:"Subnet prefix"`
}

// ModelBMSStorage contains BMS storage groups.
type ModelBMSStorage struct {
	Local  []ModelBMSStorageDetails `documentation:"Local storage"`
	System []ModelBMSStorageDetails `documentation:"System storage"`
	Data   []ModelBMSStorageDetails `documentation:"Data storage"`
	Shared []ModelBMSStorageDetails `documentation:"Shared storage"`
}

// ModelBMSStorageDetails contains one storage device description.
type ModelBMSStorageDetails struct {
	StorageClass string `documentation:"Drive class"`
	Size         string `documentation:"Drive size"`
}

// ParamsUpdateNetworkService defines CAV standard service properties.
type ParamsUpdateNetworkService struct {
	ServiceID string `documentation:"Network service identifier"`
	Adm       *bool  `documentation:"Activate tooling zone connection"`
	Os        *bool  `documentation:"Activate object storage connection"`
	Tf        *bool  `documentation:"Activate Tufin connection"`
}
