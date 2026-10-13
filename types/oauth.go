/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

import "fmt"

// ParamsCreateOAuthClient defines parameters for registering an OAuth client.
type ParamsCreateOAuthClient struct {
	TenantID        string
	ClientName      string
	SoftwareID      *string
	SoftwareVersion *string
	ClientURI       *string
	Scope           *string
}

// Validate checks required OAuth client registration parameters.
func (p ParamsCreateOAuthClient) Validate() error {
	if p.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if p.ClientName == "" {
		return fmt.Errorf("client_name is required")
	}
	return nil
}

// ModelOAuthClient represents a registered console OAuth client.
type ModelOAuthClient struct {
	ClientName              string
	ClientID                string
	GrantTypes              []string
	TokenEndpointAuthMethod string
	ClientURI               *string
	SoftwareID              *string
	SoftwareVersion         *string
	Scope                   *string
}

// ParamsCreateOAuthToken defines parameters for exchanging a tenant OAuth
// assertion or device code for an access token.
type ParamsCreateOAuthToken struct {
	TenantID   string
	ClientID   string
	GrantType  string
	Assertion  string
	DeviceCode string
}

// Validate checks required OAuth token exchange parameters.
func (p ParamsCreateOAuthToken) Validate() error {
	if p.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if p.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if p.GrantType == "" {
		return fmt.Errorf("grant_type is required")
	}
	if (p.Assertion == "") == (p.DeviceCode == "") {
		return fmt.Errorf("exactly one of assertion or device_code is required")
	}
	return nil
}

// ModelOAuthToken represents a tenant OAuth token exchange response.
type ModelOAuthToken struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	RefreshToken string
}
