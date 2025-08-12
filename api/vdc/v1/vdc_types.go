/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"regexp"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/common-go/utils"
)

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
		VCPU      int `documentation:"Number of VCPUs"`
		Used      int `documentation:"Used VCPUs"`
		Frequency int `documentation:"Frequency of the VCPUs in MHz"`
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
		ID string `documentation:"ID of the VDC to get"`

		// Name is the name of the VDC to get.
		Name string `documentation:"Name of the VDC to get"`
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

	ParamsCreateVDCStorageProfile struct {
		Class   string
		Limit   int
		Default bool
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

type (
	// * List
	apiResponseListVDC struct {
		Records []apiResponseListVDCRecord `json:"record" fakesize:"2"`
	}

	apiResponseListVDCRecord struct {
		HREF                    string `json:"href" fake:"{href_uuid}"`
		ID                      string `json:"id"`
		Name                    string `json:"name" fake:"mockvdc-{word}"`
		Description             string `json:"description" fake:"{sentence}"`
		NumberOfVMS             int    `json:"numberOfVms" fake:"{number:0,10}"`
		NumberOfRunningVMS      int    `json:"numberOfRunningVms" fake:"{number:0,10}"`
		NumberOfVAPPS           int    `json:"numberOfDeployedVApps" fake:"{number:0,10}"`
		NumberOfStorageProfiles int    `json:"numberOfStorageProfiles" fake:"{number:1,5}"`
		NumberOfDisks           int    `json:"numberOfDisks" fake:"{number:0,10}"`
	}

	// * Get
	apiResponseGetVDC struct {
		ID          string `json:"id" fake:"{urn:vdc}"`
		Name        string `json:"name" fake:"mockvdc-{word}"`
		Description string `json:"description" fake:"{sentence}"`

		IsEnabled bool `json:"isEnabled"`

		ComputeCapacity apiResponseGetVDCComputeCapacity `json:"computeCapacity"`
		Networks        apiResponseGetVDCNetworks        `json:"availableNetworks"`
		StorageProfiles apiResponseGetVDCStorageProfiles `json:"vdcStorageProfiles"`

		VCPUInMhz int `json:"vcpuInMhz2" fake:"2200"`
	}

	apiResponseGetVDCStorageProfiles struct {
		StorageProfiles []apiResponseGetVDCStorageProfile `json:"vdcStorageProfile" fakesize:"1"`
	}

	apiResponseGetVDCStorageProfile struct {
		ID   string `json:"id" fake:"{urn:vdcstorageProfile}"`
		Name string `json:"name" fake:"platinum3k_r1"`
	}

	apiResponseGetVDCNetworks struct {
		Networks []apiResponseGetVDCNetwork `json:"network" fakesize:"1"`
	}

	apiResponseGetVDCNetwork struct {
		ID   string `json:"id" fake:"{urn:network}"`
		Name string `json:"name" fake:"mocknetwork-{word}"`
	}

	apiResponseGetVDCComputeCapacity struct {
		CPU    apiResponseGetVDCComputeCapacityDetails `json:"cpu"`
		Memory apiResponseGetVDCComputeCapacityDetails `json:"memory"`
	}

	apiResponseGetVDCComputeCapacityDetails struct {
		Units    string `json:"units"`
		Limit    int    `json:"limit"`
		Reserved int    `json:"reserved"`
		Used     int    `json:"used"`
	}

	// * GetVDCMetadata

	apiResponseGetVDCMetadatas struct {
		Metadatas []apiResponseGetVDCMetadata `json:"metadataEntry" fakesize:"1"`
	}

	apiResponseGetVDCMetadata struct {
		Name  string                         `json:"key"`
		Value apiResponseGetVDCMetadataValue `json:"typedValue"`
	}

	apiResponseGetVDCMetadataValue struct {
		Value string `json:"value"`
	}

	// * CreateVDC
	apiRequestCreateVDC struct {
		VDC apiRequestCreateVDCVDC `json:"vdc"`
	}
	apiRequestCreateVDCVDC struct {
		Name                string                        `json:"name" validator:"required"`
		Description         string                        `json:"description,omitempty"`
		ServiceClass        string                        `json:"vdcServiceClass" validator:"required,oneof=ECO STD HP VOIP"`
		DisponibilityClass  string                        `json:"vdcDisponibilityClass" validator:"required,oneof=ONE-ROOM DUAL-ROOM HA-DUAL-ROOM"`
		BillingModel        string                        `json:"vdcBillingModel" validator:"required,oneof=PAYG DRAAS RESERVED"`
		StorageBillingModel string                        `json:"vdcStorageBillingModel"`
		VCPUInMhz           int                           `json:"vcpuInMhz2"`
		CPUAllocated        int                           `json:"cpuAllocated"`
		MemoryAllocated     int                           `json:"memoryAllocated"`
		StorageProfiles     []apiRequestVDCStorageProfile `json:"vdcStorageProfiles"`
	}

	apiRequestVDCStorageProfile struct {
		Class       string `json:"class"`
		Limit       int    `json:"limit"`
		Used        int    `json:"used,omitempty"`
		Default     bool   `json:"default"`
		Description string `json:"description,omitempty"`
	}

	// * UpdateVDC
	apiRequestUpdateVDC struct {
		VDC apiRequestUpdateVDCVDC `json:"vdc"`
	}

	apiRequestUpdateVDCVDC struct {
		Name            string                        `json:"name"`
		Description     string                        `json:"description,omitempty"`
		CPUAllocated    int                           `json:"cpuAllocated,omitempty"`
		MemoryAllocated int                           `json:"memoryAllocated,omitempty"`
		StorageProfiles []apiRequestVDCStorageProfile `json:"vdcStorageProfiles,omitempty"`
	}
)

var vdcRules = commands.NewRules([]commands.ConditionalRule{
	// * ----------- disponibility_class ----------- *
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console1,
			consoles.Console2,
		},
		When: commands.And(
			commands.NewCondition("service_class", "ECO"),
		).Build(),
		Target: "disponibility_class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"ONE-ROOM", "DUAL-ROOM"},
			Description: "Disponibility class allowed for Service Class ECO",
		},
	},
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console4,
			consoles.Console5,
		},
		Target: "disponibility_class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"ONE-ROOM"},
			Description: "Disponibility class allowed for Service Class ECO",
		},
	},
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console1,
			consoles.Console2,
		},
		When: commands.Or(
			commands.NewCondition("service_class", "STD"),
			commands.NewCondition("service_class", "HP"),
			commands.NewCondition("service_class", "VOIP"),
		).Build(),
		Target: "disponibility_class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"ONE-ROOM", "DUAL-ROOM", "HA-DUAL-ROOM"},
			Description: "Disponibility class allowed for Service Class STD, HP, VOIP",
		},
	},

	// * ----------- billing_model ----------- *
	{
		When: commands.Or(
			commands.NewCondition("service_class", "ECO"),
			commands.NewCondition("service_class", "STD"),
		).Build(),
		Target: "billing_model",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"PAYG", "DRAAS", "RESERVED"},
			Description: "Billing model allowed for Service Class ECO, STD",
		},
	},
	{
		When: commands.Or(
			commands.NewCondition("service_class", "HP"),
		).Build(),
		Target: "billing_model",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"PAYG", "RESERVED"},
			Description: "Billing model allowed for Service Class HP",
		},
	},
	{
		When: commands.Or(
			commands.NewCondition("service_class", "VOIP"),
		).Build(),
		Target: "billing_model",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"RESERVED"},
			Description: "Billing model allowed for Service Class VOIP",
		},
	},

	// * ----------- storage_billing_model ----------- *
	{
		Target: "storage_billing_model",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"PAYG", "RESERVED"},
			Description: "Storage billing model allowed for Service Class ECO",
		},
	},

	// * ----------- vcpu -----------

	{
		When: commands.Or(
			commands.NewCondition("billing_model", "PAYG"),
			commands.NewCondition("billing_model", "DRAAS"),
		).Build(),
		Target: "vcpu",
		Rule: commands.RuleValues{
			Editable:    true,
			Min:         utils.ToPTR(5),
			Max:         utils.ToPTR(200),
			Description: "VCPU allowed for Service Class ECO with PAYG or DRAAS billing model",
		},
	},
	{
		When: commands.Or(
			commands.NewCondition("billing_model", "RESERVED"),
		).Build(),
		Target: "vcpu",
		Rule: commands.RuleValues{
			Editable:    true,
			Min:         utils.ToPTR(2),
			Max:         utils.ToPTR(1136),
			Description: "VCPU allowed for Service Class ECO with RESERVED billing model",
		},
	},

	// * ----------- memory ----------- *

	{
		Target: "memory",
		Rule: commands.RuleValues{
			Editable:    true,
			Unit:        "GB",
			Min:         utils.ToPTR(1),
			Max:         utils.ToPTR(5120),
			Description: "Memory allowed for Service Class ECO",
		},
	},

	// * ----------- storage_profiles ----------- *
	{
		When: commands.Or(
			commands.NewCondition("disponibility_class", "ONE-ROOM"),
		).Build(),
		Target: "storage_profiles.{index}.class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"silver", "gold", "platinum3k", "platinum7k", regexp.MustCompile("^(silver|gold|platinum[3|7]{1}k)_(ocb[0-9]{1,7})$")},
			Description: "Storage profile class allowed for Disponibility Class ONE-ROOM",
		},
	},
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console1,
			consoles.Console2,
		},
		When: commands.Or(
			commands.NewCondition("disponibility_class", "DUAL-ROOM"),
		).Build(),
		Target: "storage_profiles.{index}.class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"silver_r1", "silver_r2", "gold_r1", "gold_r2", "platinum3k_r1", "platinum3k_r2", "platinum7k_r1", "platinum7k_r2", regexp.MustCompile("^(silver|gold|platinum[3|7]{1}k)_(ocb[0-9]{1,7})_(r1|r2)$")},
			Description: "Storage profile class allowed for Disponibility Class DUAL-ROOM",
		},
	},
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console1,
			consoles.Console2,
		},
		When: commands.Or(
			commands.NewCondition("disponibility_class", "HA-DUAL-ROOM"),
		).Build(),
		Target: "storage_profiles.{index}.class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"gold_hm", "platinum3k_hm", "platinum7k_hm", regexp.MustCompile("^(gold|platinum[3|7]{1}k)_(ocb[0-9]{1,7})_(hm)$")},
			Description: "Storage profile class allowed for Disponibility Class HA-DUAL-ROOM",
		},
	},
},
)

func (r *apiResponseGetVDC) toModel() ModelGetVDC {
	m := ModelGetVDC{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		ComputeCapacity: ModelGetVDCComputeCapacity{
			CPU: ModelGetVDCComputeCapacityCPU{
				VCPU: func() int {
					mhz := r.ComputeCapacity.CPU.Reserved
					if mhz == 0 {
						mhz = r.ComputeCapacity.CPU.Limit
					}
					return mhz / r.VCPUInMhz
				}(),
				Used:      r.ComputeCapacity.Memory.Used / r.VCPUInMhz,
				Frequency: r.VCPUInMhz,
			},
			Memory: ModelGetVDCComputeCapacityMemory{
				Limit: func() int {
					if r.ComputeCapacity.Memory.Limit == 0 {
						return r.ComputeCapacity.Memory.Used
					}
					return r.ComputeCapacity.Memory.Limit
				}(),
				Used: r.ComputeCapacity.Memory.Used,
			},
		},
	}

	for _, network := range r.Networks.Networks {
		m.Networks = append(m.Networks, ModelGetVDCNetwork{
			ID:   network.ID,
			Name: network.Name,
		})
	}

	for _, profile := range r.StorageProfiles.StorageProfiles {
		m.StorageProfiles = append(m.StorageProfiles, ModelGetVDCStorageProfile{
			ID:   profile.ID,
			Name: profile.Name,
		})
	}

	return m
}

func (r *apiResponseListVDCRecord) toModel() ModelListVDCDetails {
	return ModelListVDCDetails{
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

func (r *apiResponseListVDC) toModel() *ModelListVDC {
	model := &ModelListVDC{
		VDCS: make([]ModelListVDCDetails, 0),
	}

	for _, vdc := range r.Records {
		model.VDCS = append(model.VDCS, vdc.toModel())
	}

	return model
}

func vcpuToMhz(vcpu, frequency int) int {
	if vcpu <= 0 || frequency <= 0 {
		return 0
	}
	return vcpu * frequency
}

func serviceClassToCPUInMhz(serviceClass string) int {
	switch serviceClass {
	case "VOIP":
		return 3000
	default:
		return 2200
	}
}
