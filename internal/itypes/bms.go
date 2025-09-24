/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

type (
	// ApiResponseGetBMS from infraAPI
	ApiResponseGetBMS struct {
		BMSserver  []ApiResponseGetBMSServer  `json:"bms" fakesize:"1"`
		BMSNetwork []ApiResponseGetBMSNetwork `json:"network" fakesize:"1"`
	}

	ApiResponseGetBMSServer struct {
		BMSType           string                   `json:"bmsType" fake:"{randomstring:10}"` // e.g. "bms.hc1.large"
		Hostname          string                   `json:"hostname" fake:"{domain_name}"`
		OperatingSystem   string                   `json:"os" fake:"{randomstring:10}"` // e.g. "Ubuntu 20.04"
		BiosConfiguration string                   `json:"biosConfiguration" fake:"{randomstring:10}"`
		Storage           ApiResponseGetBMSStorage `json:"storage" fakesize:"1"`
	}

	ApiResponseGetBMSStorage struct {
		Local  []ApiResponseGetBMSStorageDetails `json:"local" fakesize:"1"`
		System []ApiResponseGetBMSStorageDetails `json:"system" fakesize:"1"`
		Data   []ApiResponseGetBMSStorageDetails `json:"data" fakesize:"1"`
		Shared []ApiResponseGetBMSStorageDetails `json:"shared" fakesize:"1"`
	}

	ApiResponseGetBMSStorageDetails struct {
		Class string `json:"class" fake:"{randomstring:10}"` // e.g. "NVMe" or "SSD"
		Size  string `json:"size" fake:"{number:100,2000}"`  // e.g. "500" (in GiB)
	}

	ApiResponseGetBMSNetwork struct {
		VlanId string `json:"vlanId" fake:"{id}"` // e.g. "2900"
		Subnet string `json:"subnet" fake:"{ipv4_network}"`
		Prefix int    `json:"prefix" fake:"24"`
	}
)

// From infraAPI
// func (r *ApiResponseGetOrg) ToModel() *types.ModelGetOrganization {
// 	return &types.ModelGetOrganization{
// 		Name:                r.Name,
// 		FullName:            r.FullName,
// 		Description:         r.Description,
// 		Email:               r.CustomerMail,
// 		InternetBillingMode: r.InternetBillingMode,
// 	}
// }
