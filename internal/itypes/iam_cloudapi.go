/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "encoding/xml"

// APIUser is a VMware CloudAPI user representation.
type APIUser struct {
	ID              string               `json:"id,omitempty"`
	Username        string               `json:"username,omitempty"`
	Description     *string              `json:"description,omitempty"`
	FullName        string               `json:"fullName,omitempty"`
	Email           string               `json:"email,omitempty"`
	Phone           string               `json:"phone,omitempty"`
	Enabled         *bool                `json:"enabled,omitempty"`
	Locked          *bool                `json:"locked,omitempty"`
	External        bool                 `json:"external,omitempty"`
	ProviderType    string               `json:"providerType,omitempty"`
	Password        string               `json:"password,omitempty"`
	RoleEntityRefs  []APIObjectReference `json:"roleEntityRefs,omitempty"`
	DeployedVMQuota *int                 `json:"deployedVmQuota,omitempty"`
	StoredVMQuota   *int                 `json:"storedVmQuota,omitempty"`
}

// APIResponseListUsers is a VMware CloudAPI paged user response.
type APIResponseListUsers struct {
	Users []APIUser `json:"values"`
}

// UnmarshalXML keeps legacy XML fixtures readable while the endpoint itself
// uses CloudAPI JSON in production.
func (r *APIResponseListUsers) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var legacy Users
	if err := d.DecodeElement(&legacy, &start); err != nil {
		return err
	}
	r.Users = make([]APIUser, 0, len(legacy.Users))
	for _, user := range legacy.Users {
		r.Users = append(r.Users, APIUser{
			ID: user.ID, Username: user.Name, FullName: user.FullName,
			Email: user.EmailAddress, Phone: user.Telephone,
			Description: user.Description, Enabled: &user.IsEnabled,
			ProviderType:    user.ProviderType,
			RoleEntityRefs:  []APIObjectReference{{ID: user.Role.ID, Name: user.Role.Name}},
			DeployedVMQuota: user.DeployedVMQuota, StoredVMQuota: user.StoredVMQuota,
		})
	}
	return nil
}

// UnmarshalXML keeps legacy XML user responses compatible with CloudAPI models.
func (u *APIUser) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var user User
	if err := d.DecodeElement(&user, &start); err != nil {
		return err
	}
	u.ID, u.Username, u.FullName = user.ID, user.Name, user.FullName
	u.Email, u.Phone, u.Description = user.EmailAddress, user.Telephone, user.Description
	u.Enabled, u.Locked = &user.IsEnabled, &user.IsLocked
	u.ProviderType = user.ProviderType
	u.DeployedVMQuota, u.StoredVMQuota = user.DeployedVMQuota, user.StoredVMQuota
	u.RoleEntityRefs = []APIObjectReference{{ID: user.Role.ID, Name: user.Role.Name}}
	return nil
}

// APIRequestPasswordChange is a CloudAPI password change request.
type APIRequestPasswordChange struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
