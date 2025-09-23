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
	ModelListVDC struct {
		VDCS []ModelListVDCDetails `documentation:"List of VDCs"`
	}
	ModelListVDCDetails struct {
		ID          string `documentation:"ID of the VDC"`
		Name        string `documentation:"Name of the VDC"`
		Description string `documentation:"Description of the VDC"`

		NumberOfVMS             int `documentation:"Number of VMs in the VDC"`
		NumberOfRunningVMS      int `documentation:"Number of running VMs in the VDC"`
		NumberOfVAPPS           int `documentation:"Number of deployed VApps in the VDC"`
		NumberOfStorageProfiles int `documentation:"Number of storage profiles in the VDC"`
		NumberOfDisks           int `documentation:"Number of disks in the VDC"`
	}

	ModelGetVDC struct {
		ID          string `documentation:"ID of the VDC"`
		Name        string `documentation:"Name of the VDC"`
		Description string `documentation:"Description of the VDC"`

		ComputeCapacity ModelGetVDCComputeCapacity  `documentation:"Compute capacity of the VDC"`
		Networks        []ModelGetVDCNetwork        `documentation:"Available networks in the VDC"`
		StorageProfiles []ModelGetVDCStorageProfile `documentation:"Storage profiles available in the VDC"`

		Properties ModelGetVDCProperties `documentation:"Properties of the VDC"`

		NumberOfVMS             int `documentation:"Number of VMs in the VDC"`
		NumberOfRunningVMS      int `documentation:"Number of running VMs in the VDC"`
		NumberOfVAPPS           int `documentation:"Number of deployed VApps in the VDC"`
		NumberOfStorageProfiles int `documentation:"Number of storage profiles in the VDC"`
		NumberOfDisks           int `documentation:"Number of disks in the VDC"`
	}

	ModelGetVDCStorageProfile struct {
		ID   string `documentation:"ID of the storage profile"`
		Name string `documentation:"Class name of the storage profile"`
	}

	ModelGetVDCComputeCapacity struct {
		CPU    ModelGetVDCComputeCapacityCPU    `documentation:"CPU capacity details"`
		Memory ModelGetVDCComputeCapacityMemory `documentation:"Memory capacity details"`
	}

	ModelGetVDCComputeCapacityCPU struct {
		Limit          int `documentation:"Maximum number of vCPUs that can be allocated in the VDC (CPU quota)."`
		Used           int `documentation:"Current number of vCPUs allocated in the VDC"`
		FrequencyLimit int `documentation:"Frequency of all VCPUs in MHz"`
		FrequencyUsed  int `documentation:"Used frequency of all VCPUs in MHz"`
		VCPUFrequency  int `documentation:"Frequency of a single VCPU in MHz"`
	}

	ModelGetVDCComputeCapacityMemory struct {
		Limit int `documentation:"Limit of memory in GB"`
		Used  int `documentation:"Used memory in GB"`
	}

	ModelGetVDCNetwork struct {
		ID   string `documentation:"ID of the network"`
		Name string `documentation:"Name of the network"`
	}

	ModelGetVDCProperties struct {
		ServiceClass        string `documentation:"Service class of the VDC"`
		DisponibilityClass  string `documentation:"Disponibility class of the VDC"`
		BillingModel        string `documentation:"Billing model of the VDC compute"`
		StorageBillingModel string `documentation:"Billing model of the VDC storage"`
	}
)

type (
	ParamsListVDC struct {
		// ID is the unique identifier of the VDC to filter by.
		ID string `documentation:"ID of the VDC to filter by"`
		// Name is the name of the VDC to filter by.
		Name string `documentation:"Name of the VDC to filter by"`
	}

	// ParamsGetVDC is the parameters for the GetVDC command.
	ParamsGetVDC struct {
		// ID is the unique identifier of the VDC to get.
		ID string

		// Name is the name of the VDC to get.
		Name string
	}

	// ParamsCreateVDC is the parameters for the CreateVDC command.
	ParamsCreateVDC struct {
		Name        string
		Description string

		ServiceClass        string
		DisponibilityClass  string
		BillingModel        string
		StorageBillingModel string

		Vcpu   int
		Memory int

		StorageProfiles []ParamsCreateVDCStorageProfile
	}

	// ParamsCreateVDCStorageProfile defines the parameters for creating a storage profile in a VDC.
	ParamsCreateVDCStorageProfile struct {
		Class   string
		Limit   int
		Default bool
	}

	// ParamsUpdateVDCStorageProfile defines the parameters for updating a storage profile in a VDC.
	ParamsUpdateVDCStorageProfile struct {
		Class   string
		Limit   int
		Default *bool
	}

	// ParamsDeleteVDCStorageProfile defines the parameters for deleting a storage profile from a VDC.
	ParamsDeleteVDCStorageProfile struct {
		Class string
	}

	ParamsUpdateVDC struct {
		// ID and Name cannot be updated. Will be used for retrieving the VDC.
		ID   string
		Name string

		// Editable fields
		Description *string
		Vcpu        *int
		Memory      *int
	}

	ParamsDeleteVDC struct {
		ID   string
		Name string
	}
)
