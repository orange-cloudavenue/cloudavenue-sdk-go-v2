/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "encoding/xml"

// Reference represents a VMware vCD reference (e.g., Role).
type Reference struct {
	Href string `xml:"href,attr"`
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// User represents a VMware vCD AdminOrg User XML element.
type User struct {
	XMLName         xml.Name  `xml:"User"`
	Name            string    `xml:"name,attr"`
	ID              string    `xml:"id,attr"`
	Role            Reference `xml:"Role"`
	Description     string    `xml:"Description,omitempty"`
	FullName        string    `xml:"FullName,omitempty"`
	EmailAddress    string    `xml:"EmailAddress,omitempty"`
	Telephone       string    `xml:"Telephone,omitempty"`
	IsEnabled       bool      `xml:"IsEnabled,omitempty"`
	DeployedVMQuota int       `xml:"DeployedVMQuota,omitempty"`
	StoredVMQuota   int       `xml:"StoredVMQuota,omitempty"`
	Password        string    `xml:"Password,omitempty"`
	ProviderType    string    `xml:"ProviderType,omitempty"`
}

// Users represents the wrapper for a list of users in VMware vCD AdminOrg API.
type Users struct {
	XMLName xml.Name `xml:"Users"`
	Users   []User   `xml:"User"`
}

// NewPassword represents the request body for changing a user's password.
type NewPassword struct {
	XMLName  xml.Name `xml:"NewPassword"`
	Password string   `xml:",chardata"`
}
