/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

import "fmt"

// ModelToken represents a CloudAvenue IAM token.
type ModelToken struct {
	// ID of the token in URN format
	ID string `documentation:"URN of the token in the format urn:vcloud:token:<UUID>"`

	// Name of the token
	Name string `documentation:"Name of the token"`

	// Description of the token
	Description string `documentation:"Description of the token"`

	// Indicates if the token is enabled
	Enabled bool `documentation:"Indicates if the token is enabled"`

	// Name of the token's role
	RoleName string `documentation:"Name of the token's role"`

	// Href of the token's role
	RoleHref string `documentation:"Href of the token's role"`
}

// ParamsGetToken defines parameters for getting a token.
type ParamsGetToken struct {
	ID string
}

// Validate checks ParamsGetToken structural constraints.
func (p ParamsGetToken) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsCreateToken defines parameters for creating a token.
type ParamsCreateToken struct {
	Name        string
	RoleID      string
	RoleName    string
	Description string
	IsEnabled   bool
}

// Validate checks ParamsCreateToken structural constraints.
func (p ParamsCreateToken) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.RoleName == "" {
		return fmt.Errorf("role_name is required")
	}
	return nil
}

// ParamsUpdateToken defines parameters for updating a token.
type ParamsUpdateToken struct {
	ID          string
	Name        string
	RoleID      string
	RoleName    string
	Description string
	IsEnabled   bool
}

// Validate checks ParamsUpdateToken structural constraints.
func (p ParamsUpdateToken) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ParamsDeleteToken defines parameters for deleting a token.
type ParamsDeleteToken struct {
	ID string
}

// Validate checks ParamsDeleteToken structural constraints.
func (p ParamsDeleteToken) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}
