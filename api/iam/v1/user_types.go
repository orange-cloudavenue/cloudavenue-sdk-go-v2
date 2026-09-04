/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"fmt"
)

// ParamsCreateLocalUser defines parameters for creating a local user.
type ParamsCreateLocalUser struct {
	Name            string
	Password        string
	RoleName        string
	FullName        string
	EmailAddress    string
	Telephone       string
	Description     string
	IsEnabled       bool
	DeployedVMQuota int
	StoredVMQuota   int
}

// Validate checks ParamsCreateLocalUser structural constraints.
func (p ParamsCreateLocalUser) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Password == "" {
		return fmt.Errorf("password is required")
	}
	if p.RoleName == "" {
		return fmt.Errorf("role_name is required")
	}
	return nil
}

// ParamsCreateSAMLUser defines parameters for creating a SAML user.
type ParamsCreateSAMLUser struct {
	Name         string
	RoleName     string
	FullName     string
	EmailAddress string
	Telephone    string
	Description  string
	IsEnabled    bool
}

// Validate checks ParamsCreateSAMLUser structural constraints.
func (p ParamsCreateSAMLUser) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.RoleName == "" {
		return fmt.Errorf("role_name is required")
	}
	return nil
}

// ParamsGetUser defines parameters for getting a user.
type ParamsGetUser struct {
	ID   string
	Name string
}

// Validate checks ParamsGetUser structural constraints.
func (p ParamsGetUser) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	return nil
}

// ParamsUpdateUser defines parameters for updating a user.
type ParamsUpdateUser struct {
	ID              string
	Name            string
	Password        string
	RoleName        string
	FullName        string
	EmailAddress    string
	Telephone       string
	Description     string
	IsEnabled       *bool
	DeployedVMQuota int
	StoredVMQuota   int
}

// Validate checks ParamsUpdateUser structural constraints.
func (p ParamsUpdateUser) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	return nil
}

// ParamsDeleteUser defines parameters for deleting a user.
type ParamsDeleteUser struct {
	ID            string
	TakeOwnership bool
}

// Validate checks ParamsDeleteUser structural constraints.
func (p ParamsDeleteUser) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsEnableUser defines parameters for enabling a user.
type ParamsEnableUser struct {
	ID   string
	Name string
}

// Validate checks ParamsEnableUser structural constraints.
func (p ParamsEnableUser) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	return nil
}

// ParamsDisableUser defines parameters for disabling a user.
type ParamsDisableUser struct {
	ID   string
	Name string
}

// Validate checks ParamsDisableUser structural constraints.
func (p ParamsDisableUser) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	return nil
}

// ParamsUnlockUser defines parameters for unlocking a user.
type ParamsUnlockUser struct {
	ID   string
	Name string
}

// Validate checks ParamsUnlockUser structural constraints.
func (p ParamsUnlockUser) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	return nil
}

// ParamsChangePassword defines parameters for changing a user's password.
type ParamsChangePassword struct {
	ID       string
	Name     string
	Password string
}

// Validate checks ParamsChangePassword structural constraints.
func (p ParamsChangePassword) Validate() error {
	if p.ID == "" && p.Name == "" {
		return fmt.Errorf("id or name is required")
	}
	if p.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// ModelUser represents a VMware vCD user.
type ModelUser struct {
	// ID of the user in URN format
	ID string `documentation:"URN of the user in the format urn:vcloud:user:<UUID>"`

	// Name of the user
	Name string `documentation:"Name of the user"`

	// Full name of the user
	FullName string `documentation:"Full name of the user"`

	// Email address of the user
	EmailAddress string `documentation:"Email address of the user"`

	// Telephone number of the user
	Telephone string `documentation:"Telephone number of the user"`

	// Description of the user
	Description string `documentation:"Description of the user"`

	// Indicates if the user is enabled
	IsEnabled bool `documentation:"Indicates if the user is enabled"`

	// Deployed VM quota for the user
	DeployedVMQuota int `documentation:"Deployed VM quota for the user"`

	// Stored VM quota for the user
	StoredVMQuota int `documentation:"Stored VM quota for the user"`

	// Provider type of the user (e.g., INTEGRATED, SAML)
	ProviderType string `documentation:"Provider type of the user"`

	// Name of the user's role
	RoleName string `documentation:"Name of the user's role"`

	// Href of the user's role
	RoleHref string `documentation:"Href of the user's role"`
}
