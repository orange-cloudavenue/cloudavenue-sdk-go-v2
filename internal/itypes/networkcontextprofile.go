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

type (
	// * List
	ApiResponseListNetworkContextProfile struct {
		Values []ApiResponseNetworkContextProfile `json:"values" fakesize:"3"`
	}

	// * Get / Create / Update / shared model
	ApiResponseNetworkContextProfile struct {
		ID              string                              `json:"id,omitempty" fake:"{uuid}"`
		Name            string                              `json:"name,omitempty" fake:"mocknetworkcontextprofile-{word}"`
		Description     string                              `json:"description,omitempty" fake:"{sentence}"`
		Scope           string                              `json:"scope,omitempty"`
		OrgRef          *ApiObjectReference                 `json:"orgRef,omitempty"`
		ContextEntityId string                              `json:"contextEntityId,omitempty"`
		Attributes      []ApiNetworkContextProfileAttribute `json:"attributes,omitempty"`
	}

	ApiNetworkContextProfileAttribute struct {
		Type          string                                 `json:"type"`
		Values        []string                               `json:"values,omitempty"`
		SubAttributes []ApiNetworkContextProfileSubAttribute `json:"subAttributes,omitempty"`
	}

	ApiNetworkContextProfileSubAttribute struct {
		Type   string   `json:"type"`
		Values []string `json:"values,omitempty"`
	}

	// * Create / Update request (same shape as ApiResponseNetworkContextProfile)
	ApiRequestNetworkContextProfile struct {
		ID              string                              `json:"id,omitempty"`
		Name            string                              `json:"name,omitempty"`
		Description     string                              `json:"description,omitempty"`
		Scope           string                              `json:"scope,omitempty"`
		OrgRef          *ApiObjectReference                 `json:"orgRef,omitempty"`
		ContextEntityId string                              `json:"contextEntityId,omitempty"`
		Attributes      []ApiNetworkContextProfileAttribute `json:"attributes,omitempty"`
	}

	// * Attributes static reference catalog (live, server-computed)
	ApiNetworkContextProfileAttributesResponse struct {
		Attributes []ApiNetworkContextProfileAttribute `json:"attributes,omitempty"`
	}
)

func (r *ApiResponseListNetworkContextProfile) ToModel() *types.ModelListNetworkContextProfile {
	model := &types.ModelListNetworkContextProfile{
		NetworkContextProfiles: make([]types.ModelGetNetworkContextProfile, 0),
	}

	for _, profile := range r.Values {
		model.NetworkContextProfiles = append(model.NetworkContextProfiles, profile.ToModel())
	}

	return model
}

func (r *ApiResponseNetworkContextProfile) ToModel() types.ModelGetNetworkContextProfile {
	m := types.ModelGetNetworkContextProfile{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Scope:       r.Scope,
	}

	if r.OrgRef != nil {
		m.OrgID = r.OrgRef.ID
	}

	for _, attr := range r.Attributes {
		modelAttr := types.ModelNetworkContextProfileAttribute{
			Type:   attr.Type,
			Values: attr.Values,
		}

		for _, sub := range attr.SubAttributes {
			modelAttr.SubAttributes = append(modelAttr.SubAttributes, types.ModelNetworkContextProfileSubAttribute{
				Type:   sub.Type,
				Values: sub.Values,
			})
		}

		m.Attributes = append(m.Attributes, modelAttr)
	}

	return m
}

func (r *ApiRequestNetworkContextProfile) ToModel() types.ModelGetNetworkContextProfile {
	m := types.ModelGetNetworkContextProfile{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Scope:       r.Scope,
	}

	if r.OrgRef != nil {
		m.OrgID = r.OrgRef.ID
	}

	for _, attr := range r.Attributes {
		modelAttr := types.ModelNetworkContextProfileAttribute{
			Type:   attr.Type,
			Values: attr.Values,
		}

		for _, sub := range attr.SubAttributes {
			modelAttr.SubAttributes = append(modelAttr.SubAttributes, types.ModelNetworkContextProfileSubAttribute{
				Type:   sub.Type,
				Values: sub.Values,
			})
		}

		m.Attributes = append(m.Attributes, modelAttr)
	}

	return m
}
