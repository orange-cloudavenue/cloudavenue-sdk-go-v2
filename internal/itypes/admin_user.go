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

// AdminUser represents the legacy admin XML user record. CloudAPI user
// responses omit the deployed/stored VM quotas, which only this
// representation exposes.
type AdminUser struct {
	// XMLName deliberately has no element constraint: the admin API serves
	// this record under User and AdminUser roots depending on API version.
	XMLName         xml.Name
	Href            string `xml:"href,attr,omitempty"`
	ID              string `xml:"id,attr,omitempty"`
	Name            string `xml:"name,attr,omitempty"`
	FullName        string `xml:"FullName,omitempty"`
	Description     string `xml:"Description,omitempty"`
	IsEnabled       bool   `xml:"IsEnabled,omitempty"`
	IsLocked        bool   `xml:"IsLocked,omitempty"`
	DeployedVMQuota int    `xml:"DeployedVmQuota,omitempty"`
	StoredVMQuota   int    `xml:"StoredVmQuota,omitempty"`
	ProviderType    string `xml:"ProviderType,omitempty"`
}
