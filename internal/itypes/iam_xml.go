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
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// User represents a legacy VMware vCD AdminOrg User XML element retained for
// XML fixtures and compatibility types; IAM user operations use CloudAPI.
type User struct {
	XMLName         xml.Name  `xml:"User"`
	Xmlns           string    `xml:"xmlns,attr,omitempty"`
	Href            string    `xml:"href,attr,omitempty"`
	Type            string    `xml:"type,attr,omitempty"`
	Name            string    `xml:"name,attr"`
	ID              string    `xml:"id,attr"`
	Role            Reference `xml:"Role"`
	Description     *string   `xml:"Description,omitempty"`
	FullName        string    `xml:"FullName,omitempty"`
	EmailAddress    string    `xml:"EmailAddress,omitempty"`
	Telephone       string    `xml:"Telephone,omitempty"`
	IsEnabled       bool      `xml:"IsEnabled,omitempty"`
	IsLocked        bool      `xml:"IsLocked,omitempty"`
	IM              string    `xml:"IM,omitempty"`
	NameInSource    string    `xml:"NameInSource,omitempty"`
	IsExternal      bool      `xml:"IsExternal,omitempty"`
	IsGroupRole     bool      `xml:"IsGroupRole,omitempty"`
	DeployedVMQuota *int      `xml:"DeployedVmQuota,omitempty"`
	StoredVMQuota   *int      `xml:"StoredVmQuota,omitempty"`
	Password        string    `xml:"Password,omitempty"`
	ProviderType    string    `xml:"ProviderType,omitempty"`
}

const (
	UserProviderSAML  = "SAML"
	UserProviderLocal = "INTEGRATED"
)

// UserRequest represents a legacy VMware vCD AdminOrg User XML request body.
type UserRequest struct {
	XMLName         xml.Name  `xml:"User"`
	Xmlns           string    `xml:"xmlns,attr,omitempty"`
	Name            string    `xml:"name,attr"`
	ID              string    `xml:"id,attr"`
	Role            Reference `xml:"Role"`
	Description     *string   `xml:"Description,omitempty"`
	FullName        string    `xml:"FullName,omitempty"`
	EmailAddress    string    `xml:"EmailAddress,omitempty"`
	Telephone       string    `xml:"Telephone,omitempty"`
	IsEnabled       *bool     `xml:"IsEnabled,omitempty"`
	IsLocked        *bool     `xml:"IsLocked,omitempty"`
	IM              string    `xml:"IM,omitempty"`
	NameInSource    string    `xml:"NameInSource,omitempty"`
	IsExternal      *bool     `xml:"IsExternal,omitempty"`
	IsGroupRole     *bool     `xml:"IsGroupRole,omitempty"`
	DeployedVMQuota *int      `xml:"DeployedVmQuota,omitempty"`
	StoredVMQuota   *int      `xml:"StoredVmQuota,omitempty"`
	Password        string    `xml:"Password,omitempty"`
	ProviderType    string    `xml:"ProviderType,omitempty"`
}

// Users represents a legacy VMware vCD AdminOrg XML user-list response.
type Users struct {
	XMLName xml.Name `xml:"Users"`
	Users   []User   `xml:"User"`
}

// NewPassword represents the request body for changing a user's password.
type NewPassword struct {
	XMLName  xml.Name `xml:"NewPassword"`
	Password string   `xml:",chardata"`
}
