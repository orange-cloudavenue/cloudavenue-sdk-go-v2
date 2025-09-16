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
	Email               string
	InternetBillingMode string
}

type ModelGetOrganization struct {
	// ID of organization in urn format
	// Example: urn:vcloud:org:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	ID string `documentation:"URN of the organization in the format urn:vcloud:org:<UUID>"`

	// Name of the organization
	Name string `documentation:"Name of the organization"`

	// Full name of the organization
	FullName string `documentation:"Long name of the organization"`

	// Description of the organization
	Description string `documentation:"Description of the organization"`

	// Indicates if the organization is enabled
	Enabled bool `documentation:"Indicates if the organization is enabled"`

	// Contact email of the organization
	Email string `documentation:"Contact email of the organization"`

	// Internet billing mode of the organization
	InternetBillingMode string `documentation:"Internet billing mode of the organization"`

	// Represent a details count of resources in the organization
	Resources ModelGetOrganizationResources `documentation:"Resources usage details of the organization"`
}

type ModelGetOrganizationResources struct {
	// Number of Org VDCs
	Vdc int `documentation:"Number of VDC(s) in your organization"`

	// Number of Catalog media(s)
	Catalog int `documentation:"Number of Catalog media(s)"`

	// Number of vApps
	Vapp int `documentation:"Number of vApp(s)"`

	// Number of VM(s) in power on state
	VMRunning int `documentation:"Number of VM(s) in state power on"`

	// Number of users in the organization
	User int `documentation:"Number of user(s) in the organization"`

	// Number of standalone disks in the organization
	Disk int `documentation:"Number of standalone disk(s) configured in the organization"`
}
