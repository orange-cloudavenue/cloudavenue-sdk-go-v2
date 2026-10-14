/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

// APIRequestOAuthClient is the console OAuth client registration request.
type APIRequestOAuthClient struct {
	ClientName      string  `json:"client_name" validate:"required"`
	SoftwareID      *string `json:"software_id,omitempty"`
	SoftwareVersion *string `json:"software_version,omitempty"`
	ClientURI       *string `json:"client_uri,omitempty"`
	Scope           *string `json:"scope,omitempty"`
}

// APIResponseOAuthClient is the console OAuth client registration response.
type APIResponseOAuthClient struct {
	ClientName              string   `json:"client_name"`
	ClientID                string   `json:"client_id"`
	GrantTypes              []string `json:"grant_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	ClientURI               *string  `json:"client_uri"`
	SoftwareID              *string  `json:"software_id"`
	SoftwareVersion         *string  `json:"software_version"`
	Scope                   *string  `json:"scope"`
}

// APIRequestOAuthToken is the form request for exchanging a tenant OAuth
// assertion or device code for an access token.
type APIRequestOAuthToken struct {
	ClientID   string
	GrantType  string
	Assertion  string
	DeviceCode string
}

// APIResponseOAuthToken is the tenant OAuth token exchange response.
type APIResponseOAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
