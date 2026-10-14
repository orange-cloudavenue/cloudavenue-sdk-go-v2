/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

// * Request / Response API for Tokens

type (
	// APIRequestToken represents the request body for token operations.
	APIRequestToken struct {
		ID          string             `json:"id,omitempty"`
		Name        string             `json:"name,omitempty"`
		Description string             `json:"description,omitempty"`
		Role        APIObjectReference `json:"role,omitempty"`
		Enabled     bool               `json:"enabled,omitempty"`
	}

	// APIResponseToken represents a token response.
	APIResponseToken struct {
		ID          string             `json:"id,omitempty"`
		Name        string             `json:"name,omitempty"`
		Description string             `json:"description,omitempty"`
		Enabled     bool               `json:"enabled,omitempty"`
		Role        APIObjectReference `json:"role,omitempty"`
	}

	// APIResponseListTokens represents the wrapper for a list of tokens.
	APIResponseListTokens struct {
		Tokens []APIResponseToken `json:"tokens,omitempty"`
	}
)

// ToModel converts the APIResponseToken to ModelToken.
func (api *APIResponseToken) ToModel() *types.ModelToken {
	if api == nil {
		return nil
	}

	return &types.ModelToken{
		ID:          api.ID,
		Name:        api.Name,
		Description: api.Description,
		Enabled:     api.Enabled,
		RoleName:    api.Role.Name,
		RoleHref:    api.Role.ID,
	}
}

// ToModel converts each token in the list to ModelToken.
func (api *APIResponseListTokens) ToModel() []*types.ModelToken {
	if api == nil {
		return nil
	}

	result := make([]*types.ModelToken, len(api.Tokens))
	for i, t := range api.Tokens {
		result[i] = t.ToModel()
	}

	return result
}
