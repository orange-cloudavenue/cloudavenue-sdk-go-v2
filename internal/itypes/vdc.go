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
	// * List
	APIResponseListVDC struct {
		Records []APIResponseListVDCRecord `json:"record" fakesize:"2"`
	}

	APIResponseListVDCRecord struct {
		HREF                    string `json:"href" fake:"{href_uuid}"`
		ID                      string `json:"id" fake:"{urn:vdc}"`
		Name                    string `json:"name" fake:"mockvdc-{word}"`
		Description             string `json:"description" fake:"{sentence}"`
		NumberOfVMS             int    `json:"numberOfVms" fake:"{number:0,10}"`
		NumberOfRunningVMS      int    `json:"numberOfRunningVms" fake:"{number:0,10}"`
		NumberOfVAPPS           int    `json:"numberOfDeployedVApps" fake:"{number:0,10}"`
		NumberOfStorageProfiles int    `json:"numberOfStorageProfiles" fake:"{number:1,5}"`
		NumberOfDisks           int    `json:"numberOfDisks" fake:"{number:0,10}"`
	}

	// * Get
	APIResponseGetVDC struct {
		ID          string `json:"id" fake:"{urn:vdc}"`
		Name        string `json:"name" fake:"mockvdc-{word}"`
		Description string `json:"description" fake:"{sentence}"`

		IsEnabled bool `json:"isEnabled"`

		ComputeCapacity APIResponseGetVDCComputeCapacity `json:"computeCapacity"`
		Networks        APIResponseGetVDCNetworks        `json:"availableNetworks"`
		StorageProfiles APIResponseGetVDCStorageProfiles `json:"vdcStorageProfiles"`

		VCPUInMhz int `json:"vcpuInMhz2" fake:"2200"`
	}

	APIResponseGetVDCStorageProfiles struct {
		StorageProfiles []APIResponseGetVDCStorageProfile `json:"vdcStorageProfile" fakesize:"1"`
	}

	APIResponseGetVDCStorageProfile struct {
		ID   string `json:"id" fake:"{urn:vdcstorageProfile}"`
		Name string `json:"name" fake:"platinum3k_r1"`
	}

	APIResponseGetVDCNetworks struct {
		Networks []APIResponseGetVDCNetwork `json:"network" fakesize:"1"`
	}

	APIResponseGetVDCNetwork struct {
		ID   string `json:"id" fake:"{urn:network}"`
		Name string `json:"name" fake:"mocknetwork-{word}"`
	}

	APIResponseGetVDCComputeCapacity struct {
		CPU    APIResponseGetVDCComputeCapacityDetails `json:"cpu"`
		Memory APIResponseGetVDCComputeCapacityDetails `json:"memory"`
	}

	APIResponseGetVDCComputeCapacityDetails struct {
		Units     string `json:"units"`
		Limit     int    `json:"limit"`
		Allocated int    `json:"allocated"`
		Used      int    `json:"used"`
	}

	// * GetVDCMetadata

	APIResponseGetVDCMetadatas struct {
		Metadatas []APIResponseGetVDCMetadata `json:"metadataEntry" fakesize:"1"`
	}

	APIResponseGetVDCMetadata struct {
		Name  string                         `json:"key"`
		Value APIResponseGetVDCMetadataValue `json:"typedValue"`
	}

	APIResponseGetVDCMetadataValue struct {
		Value string `json:"value"`
	}

	// * CreateVDC
	APIRequestCreateVDC struct {
		VDC APIRequestCreateVDCVDC `json:"vdc"`
	}
	APIRequestCreateVDCVDC struct {
		Name                string                        `json:"name" validator:"required"`
		Description         string                        `json:"description,omitempty"`
		ServiceClass        string                        `json:"vdcServiceClass" validator:"required,oneof=ECO STD HP VOIP"`
		DisponibilityClass  string                        `json:"vdcDisponibilityClass" validator:"required,oneof=ONE-ROOM DUAL-ROOM HA-DUAL-ROOM"`
		BillingModel        string                        `json:"vdcBillingModel" validator:"required,oneof=PAYG DRAAS RESERVED"`
		StorageBillingModel string                        `json:"vdcStorageBillingModel"`
		VCPUInMhz           int                           `json:"vcpuInMhz2"`
		CPUAllocated        int                           `json:"cpuAllocated"`
		MemoryAllocated     int                           `json:"memoryAllocated"`
		StorageProfiles     []APIRequestVDCStorageProfile `json:"vdcStorageProfiles"`
	}

	APIRequestVDCStorageProfile struct {
		Class       string `json:"class"`
		Limit       int    `json:"limit"`
		Used        int    `json:"used,omitempty"`
		Default     bool   `json:"default"`
		Description string `json:"description,omitempty"`
	}

	// * UpdateVDC
	APIRequestUpdateVDC struct {
		VDC APIRequestUpdateVDCVDC `json:"vdc"`
	}

	APIRequestUpdateVDCVDC struct {
		Name            string                        `json:"name"`
		Description     string                        `json:"description,omitempty"`
		CPUAllocated    int                           `json:"cpuAllocated,omitempty"`
		MemoryAllocated int                           `json:"memoryAllocated,omitempty"`
		StorageProfiles []APIRequestVDCStorageProfile `json:"vdcStorageProfiles,omitempty"`
	}
)

func (r *APIResponseGetVDC) ToModel() types.ModelGetVDC {
	m := types.ModelGetVDC{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		ComputeCapacity: types.ModelGetVDCComputeCapacity{
			CPU: types.ModelGetVDCComputeCapacityCPU{
				Limit: func() int {
					mhz := r.ComputeCapacity.CPU.Allocated
					if mhz == 0 {
						mhz = r.ComputeCapacity.CPU.Limit
					}
					return mhz / r.VCPUInMhz
				}(),
				Used: r.ComputeCapacity.CPU.Used / r.VCPUInMhz,
				FrequencyLimit: func() int {
					if r.ComputeCapacity.CPU.Allocated != 0 {
						return r.ComputeCapacity.CPU.Allocated
					}
					return r.ComputeCapacity.CPU.Limit
				}(),
				FrequencyUsed: r.ComputeCapacity.CPU.Used,
				VCPUFrequency: r.VCPUInMhz,
			},
			Memory: types.ModelGetVDCComputeCapacityMemory{
				Limit: r.ComputeCapacity.Memory.Limit,
				Used:  r.ComputeCapacity.Memory.Used,
			},
		},
	}

	for _, network := range r.Networks.Networks {
		m.Networks = append(m.Networks, types.ModelVDCNetworkRef{
			ID:   network.ID,
			Name: network.Name,
		})
	}

	for _, profile := range r.StorageProfiles.StorageProfiles {
		m.StorageProfiles = append(m.StorageProfiles, types.ModelGetVDCStorageProfile{
			ID:   profile.ID,
			Name: profile.Name,
		})
	}

	return m
}

func (r *APIResponseListVDCRecord) ToModel() types.ModelListVDCDetails {
	return types.ModelListVDCDetails{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,

		NumberOfVMS:             r.NumberOfVMS,
		NumberOfRunningVMS:      r.NumberOfRunningVMS,
		NumberOfVAPPS:           r.NumberOfVAPPS,
		NumberOfStorageProfiles: r.NumberOfStorageProfiles,
		NumberOfDisks:           r.NumberOfDisks,
	}
}

func (r *APIResponseListVDC) ToModel() *types.ModelListVDC {
	model := &types.ModelListVDC{
		VDCS: make([]types.ModelListVDCDetails, 0),
	}

	for _, vdc := range r.Records {
		model.VDCS = append(model.VDCS, vdc.ToModel())
	}

	return model
}
