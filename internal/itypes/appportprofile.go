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
	ApiResponseListAppPortProfile struct {
		Values []ApiResponseAppPortProfile `json:"values" fakesize:"3"`
	}

	// * Get / Create / Update / shared model
	ApiResponseAppPortProfile struct {
		ID               string                  `json:"id,omitempty" fake:"{urn:applicationPortProfile}"`
		Name             string                  `json:"name,omitempty" fake:"mockappportprofile-{word}"`
		Description      string                  `json:"description,omitempty" fake:"{sentence}"`
		ApplicationPorts []ApiAppPortProfilePort `json:"applicationPorts,omitempty"`
		OrgRef           *ApiObjectReference     `json:"orgRef,omitempty"`
		ContextEntityId  string                  `json:"contextEntityId,omitempty"`
		Scope            string                  `json:"scope,omitempty"`
	}

	ApiAppPortProfilePort struct {
		Protocol         string   `json:"protocol"`
		DestinationPorts []string `json:"destinationPorts,omitempty"`
	}

	// * Create / Update request (same shape as ApiResponseAppPortProfile)
	ApiRequestAppPortProfile struct {
		ID               string                  `json:"id,omitempty" fake:"{urn:applicationPortProfile}"`
		Name             string                  `json:"name,omitempty" fake:"mockappportprofile-{word}"`
		Description      string                  `json:"description,omitempty" fake:"{sentence}"`
		ApplicationPorts []ApiAppPortProfilePort `json:"applicationPorts,omitempty"`
		OrgRef           *ApiObjectReference     `json:"orgRef,omitempty"`
		ContextEntityId  string                  `json:"contextEntityId,omitempty"`
		Scope            string                  `json:"scope,omitempty"`
	}
)

func (r *ApiResponseListAppPortProfile) ToModel() *types.ModelListAppPortProfile {
	model := &types.ModelListAppPortProfile{
		AppPortProfiles: make([]types.ModelGetAppPortProfile, 0),
	}

	for _, profile := range r.Values {
		model.AppPortProfiles = append(model.AppPortProfiles, profile.ToModel())
	}

	return model
}

func (r *ApiResponseAppPortProfile) ToModel() types.ModelGetAppPortProfile {
	m := types.ModelGetAppPortProfile{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Scope:       r.Scope,
	}

	if r.OrgRef != nil {
		m.OrgID = r.OrgRef.ID
	}

	for _, port := range r.ApplicationPorts {
		m.ApplicationPorts = append(m.ApplicationPorts, types.ModelGetAppPortProfilePort{
			Protocol:         port.Protocol,
			DestinationPorts: port.DestinationPorts,
		})
	}

	return m
}

func (r *ApiRequestAppPortProfile) ToModel() types.ModelGetAppPortProfile {
	m := types.ModelGetAppPortProfile{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Scope:       r.Scope,
	}

	if r.OrgRef != nil {
		m.OrgID = r.OrgRef.ID
	}

	for _, port := range r.ApplicationPorts {
		m.ApplicationPorts = append(m.ApplicationPorts, types.ModelGetAppPortProfilePort{
			Protocol:         port.Protocol,
			DestinationPorts: port.DestinationPorts,
		})
	}

	return m
}
