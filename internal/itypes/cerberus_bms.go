/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

type APIResponseBMS struct {
	Network []APIResponseBMSNetworkDetails `json:"network"`
	BMS     []APIResponseBMSDetails        `json:"bms"`
}

type APIResponseBMSDetails struct {
	BMSType           string                `json:"bmsType"`
	Hostname          string                `json:"hostname"`
	OS                string                `json:"os"`
	BIOSConfiguration string                `json:"biosConfiguration"`
	Storage           APIResponseBMSStorage `json:"storage"`
}

type APIResponseBMSNetworkDetails struct {
	VLANID int    `json:"vlanid"`
	Subnet string `json:"subnet"`
	Prefix int    `json:"prefix"`
}

type APIResponseBMSStorage struct {
	Local  []APIResponseBMSStorageDetails `json:"local"`
	System []APIResponseBMSStorageDetails `json:"system"`
	Data   []APIResponseBMSStorageDetails `json:"data"`
	Shared []APIResponseBMSStorageDetails `json:"shared"`
}

type APIResponseBMSStorageDetails struct {
	StorageClass string `json:"storageClass"`
	Size         string `json:"size"`
}

type APIRequestNetworkServiceUpdate struct {
	Adm *bool `json:"adm,omitempty"`
	Os  *bool `json:"os,omitempty"`
	Tf  *bool `json:"tf,omitempty"`
}
