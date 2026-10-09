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
	"bytes"
	"encoding/json"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

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

		ServiceClass        string `json:"vdcServiceClass,omitempty"`
		DisponibilityClass  string `json:"vdcDisponibilityClass,omitempty"`
		BillingModel        string `json:"vdcBillingModel,omitempty"`
		StorageBillingModel string `json:"vdcStorageBillingModel,omitempty"`
		CPUAllocated        int    `json:"cpuAllocated,omitempty"`
		MemoryAllocated     int    `json:"memoryAllocated,omitempty"`
	}

	APIResponseGetVDCStorageProfiles struct {
		StorageProfiles []APIResponseGetVDCStorageProfile `json:"vdcStorageProfile" fakesize:"1"`
	}

	APIResponseGetVDCStorageProfile struct {
		ID      string `json:"id" fake:"{urn:vdcstorageProfile}"`
		Name    string `json:"name" fake:"platinum3k_r1"`
		Class   string `json:"class,omitempty"`
		Limit   int    `json:"limit,omitempty"`
		Used    int    `json:"used,omitempty"`
		Default bool   `json:"default,omitempty"`
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

// UnmarshalJSON accepts both the legacy VMware query shape and Cerberus' array
// response. Keeping one response type preserves callers that resolve VDC IDs
// through the shared endpoint.
func (r *APIResponseListVDC) UnmarshalJSON(data []byte) error {
	if len(bytes.TrimSpace(data)) > 0 && bytes.TrimSpace(data)[0] == '[' {
		var values []struct {
			Name string `json:"vdc_name"` //nolint:tagliatelle // legacy API field
			ID   string `json:"vdc_uuid"` //nolint:tagliatelle // legacy API field
		}
		if err := json.Unmarshal(data, &values); err != nil {
			return err
		}
		r.Records = make([]APIResponseListVDCRecord, 0, len(values))
		for _, value := range values {
			r.Records = append(r.Records, APIResponseListVDCRecord{ID: value.ID, Name: value.Name})
		}
		return nil
	}
	type response APIResponseListVDC
	var legacy response
	if err := json.Unmarshal(data, &legacy); err != nil {
		return err
	}
	*r = APIResponseListVDC(legacy)
	return nil
}

// UnmarshalJSON unwraps Cerberus' {"vdc": {...}} detail response while
// retaining compatibility with the legacy direct VMware response.
func (r *APIResponseGetVDC) UnmarshalJSON(data []byte) error {
	var wrapped struct {
		VDC json.RawMessage `json:"vdc"`
	}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}
	if len(wrapped.VDC) > 0 && string(wrapped.VDC) != "null" {
		type response APIResponseGetVDC
		var detail response
		if err := json.Unmarshal(wrapped.VDC, &detail); err != nil {
			return err
		}
		*r = APIResponseGetVDC(detail)
		return nil
	}
	type response APIResponseGetVDC
	var direct response
	if err := json.Unmarshal(data, &direct); err != nil {
		return err
	}
	*r = APIResponseGetVDC(direct)
	return nil
}

func (r *APIResponseGetVDC) ToModel() types.ModelGetVDC {
	vCPUInMhz := r.VCPUInMhz
	if vCPUInMhz == 0 {
		vCPUInMhz = 2200
	}
	cpuAllocated := r.ComputeCapacity.CPU.Allocated
	if cpuAllocated == 0 {
		cpuAllocated = r.CPUAllocated
	}
	memoryAllocated := r.ComputeCapacity.Memory.Limit
	if memoryAllocated == 0 {
		memoryAllocated = r.MemoryAllocated
	}
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
					return mhz / vCPUInMhz
				}(),
				Used: r.ComputeCapacity.CPU.Used / vCPUInMhz,
				FrequencyLimit: func() int {
					if r.ComputeCapacity.CPU.Allocated != 0 {
						return r.ComputeCapacity.CPU.Allocated
					}
					return r.ComputeCapacity.CPU.Limit
				}(),
				FrequencyUsed: r.ComputeCapacity.CPU.Used,
				VCPUFrequency: vCPUInMhz,
			},
			Memory: types.ModelGetVDCComputeCapacityMemory{
				Limit: memoryAllocated,
				Used:  r.ComputeCapacity.Memory.Used,
			},
		},
		Properties: types.ModelGetVDCProperties{
			ServiceClass: r.ServiceClass, DisponibilityClass: r.DisponibilityClass,
			BillingModel: r.BillingModel, StorageBillingModel: r.StorageBillingModel,
		},
	}
	if m.ComputeCapacity.CPU.FrequencyLimit == 0 {
		m.ComputeCapacity.CPU.FrequencyLimit = cpuAllocated
	}

	for _, network := range r.Networks.Networks {
		m.Networks = append(m.Networks, types.ModelVDCNetworkRef{
			ID:   network.ID,
			Name: network.Name,
		})
	}

	for _, profile := range r.StorageProfiles.StorageProfiles {
		profileName := profile.Name
		if profileName == "" {
			profileName = profile.Class
		}
		m.StorageProfiles = append(m.StorageProfiles, types.ModelGetVDCStorageProfile{
			ID:      profile.ID,
			Name:    profileName,
			Class:   profile.Class,
			Limit:   profile.Limit,
			Default: profile.Default,
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
