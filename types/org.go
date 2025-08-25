/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

// ParamsUpdateOrganization defines the parameters for updating an organization
type ParamsUpdateOrganization struct {
	FullName            string
	Description         *string
	CustomerMail        string
	InternetBillingMode string
}

type ModelGetOrganization struct {
	// ID of organization in urn format
	// Example: urn:vcloud:org:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string `documentation: "URN of the organization in the format urn:vcloud:org:<UUID>"`

	// Name of the organization
	Name string `documentation: "Name of the organization"`

	// Display name of the organization
	DisplayName string `documentation: "Display name of the organization"`

	// Description of the organization
	Description string `documentation: "Description of the organization"`

	// Indicates if the organization is enabled
	IsEnabled bool `documentation: "Indicates if the organization is enabled"`

	Resources ModelGetOrganizationResources `documentation: "Resources usage details of the organization"`

	// Contact email of the organization
	CustomerMail string `documentation: "Contact email of the organization"`

	// Internet billing mode of the organization
	InternetBillingMode string `documentation: "Internet billing mode of the organization"`
}

type ModelGetOrganizationResources struct {
	// Number of Org VDCs
	Vdc int `documentation: "Number of Org VDCs"`

	// Number of Catalog media(s)
	Catalog int `documentation: "Number of Catalog media(s)"`

	// Number of vApps
	Vapp int `documentation: "Number of vApps"`

	// Number of VM(s) in power on state
	RunningVM int `documentation: "Number of VM(s) in power on state"`

	// Number of users in the organization
	User int `documentation: "Number of users in the organization"`

	// Number of standalone disks in the organization
	Disk int `documentation: "Number of standalone disks in the organization"`
}

type ModelUpdateOrganization struct {
	FullName            string `json:"fullName"`
	Description         string `json:"description"`
	CustomerMail        string `json:"customerMail"`
	InternetBillingMode string `json:"internetBillingMode"`
}
