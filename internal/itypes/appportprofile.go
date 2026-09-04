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
	APIResponseListAppPortProfile struct {
		Values []APIResponseAppPortProfile `json:"values" fakesize:"3"`
	}

	// * Get / Create / Update / shared model
	APIResponseAppPortProfile struct {
		ID               string                  `json:"id,omitempty" fake:"{urn:applicationPortProfile}"`
		Name             string                  `json:"name,omitempty" fake:"mockappportprofile-{word}"`
		Description      string                  `json:"description,omitempty" fake:"{sentence}"`
		ApplicationPorts []APIAppPortProfilePort `json:"applicationPorts,omitempty"`
		OrgRef           *APIObjectReference     `json:"orgRef,omitempty"`
		ContextEntityID  string                  `json:"contextEntityId,omitempty"`
		Scope            string                  `json:"scope,omitempty"`
	}

	APIAppPortProfilePort struct {
		Protocol         string   `json:"protocol"`
		DestinationPorts []string `json:"destinationPorts,omitempty"`
	}

	// * Create / Update request (same shape as APIResponseAppPortProfile)
	APIRequestAppPortProfile struct {
		ID               string                  `json:"id,omitempty" fake:"{urn:applicationPortProfile}"`
		Name             string                  `json:"name,omitempty" fake:"mockappportprofile-{word}"`
		Description      string                  `json:"description,omitempty" fake:"{sentence}"`
		ApplicationPorts []APIAppPortProfilePort `json:"applicationPorts,omitempty"`
		OrgRef           *APIObjectReference     `json:"orgRef,omitempty"`
		ContextEntityID  string                  `json:"contextEntityId,omitempty"`
		Scope            string                  `json:"scope,omitempty"`
	}
)

func (r *APIResponseListAppPortProfile) ToModel() *types.ModelListAppPortProfile {
	model := &types.ModelListAppPortProfile{
		AppPortProfiles: make([]types.ModelGetAppPortProfile, 0),
	}

	for _, profile := range r.Values {
		model.AppPortProfiles = append(model.AppPortProfiles, profile.ToModel())
	}

	return model
}

func (r *APIResponseAppPortProfile) ToModel() types.ModelGetAppPortProfile {
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

func (r *APIRequestAppPortProfile) ToModel() types.ModelGetAppPortProfile {
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
