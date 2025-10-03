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
	ID string `documentation:"Uniq organization Identifier in the format urn:vcloud:org:<UUID>"`

	// Name of the organization
	Name string `documentation:"Name of the organization"`

	// Full name of the organization
	FullName string `documentation:"Display name of the organization as shown in VMware Cloud Director."`

	// Description of the organization
	Description string `documentation:"Human-readable description of the organization."`

	// Indicates if the organization is enabled
	Enabled bool `documentation:"Indicates whether the organization is enabled. When false, access and resource operations are suspended."`

	// Contact email of the organization
	Email string `documentation:"Primary contact email for the organization."`

	// Internet billing mode of the organization
	InternetBillingMode string `documentation:"Internet bandwidth billing method for the organization (for example, PAYG or TRAFFIC_VOLUME). Choose the model that best matches your traffic profile."`
}
