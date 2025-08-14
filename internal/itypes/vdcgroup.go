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
	ApiResponseListVdcGroup struct {
		Values []ApiResponseListVdcGroupDetails `json:"values" fakesize:"3"`
	}

	ApiResponseListVdcGroupDetails struct {
		ID          string                                `json:"id" fake:"{urn:vdcGroup}"`
		OrgID       string                                `json:"orgId" fake:"{urn:org}"`
		Name        string                                `json:"name" fake:"mockvdcgroup-{word}"`
		Description string                                `json:"description" fake:"{sentence}"`
		Vdcs        []ApiResponseVdcGroupParticipatingVdc `json:"participatingOrgVdcs" fakesize:"2"`
	}

	ApiResponseVdcGroupParticipatingVdc struct {
		Vdc                  ApiResponseVdcGroupParticipatingVdcRef  `json:"vdcRef"`
		Site                 ApiResponseVdcGroupParticipatingSiteRef `json:"siteRef"`
		FaultDomainTag       string                                  `json:"faultDomainTag" fake:"AZ01"`
		NetworkProviderScope string                                  `json:"networkProviderScope" fake:"AZ01"`
	}
	ApiResponseVdcGroupParticipatingSiteRef struct {
		ID   string `json:"id,omitempty" fake:"{urn:vdc}"`
		Name string `json:"name,omitempty" fake:"mockvdc-{word}"`
	}

	ApiResponseVdcGroupParticipatingVdcRef struct {
		ID   string `json:"id,omitempty" fake:"{urn:vdc}"`
		Name string `json:"name,omitempty" fake:"mockvdc-{word}"`
	}

	// * Create
	ApiRequestCreateVdcGroup struct {
		OrgID               string                                `json:"orgId" fake:"{org}"`
		Name                string                                `json:"name" fake:"mockvdcgroup-{word}"`
		Description         string                                `json:"description,omitempty" fake:"{sentence}"`
		Vdcs                []ApiResponseVdcGroupParticipatingVdc `json:"participatingOrgVdcs" fakesize:"2"`
		NetworkProviderType string                                `json:"networkProviderType" fake:"NSX_T"`
		Type                string                                `json:"type" fake:"LOCAL"`
	}

	// * Update
	ApiRequestUpdateVdcGroup struct {
		Id                  string                                `json:"id" fake:"{urn:vdcGroup}"`
		OrgID               string                                `json:"orgId" fake:"{org}"`
		Name                string                                `json:"name" fake:"mockvdcgroup-{word}"`
		Description         string                                `json:"description,omitempty" fake:"{sentence}"`
		Vdcs                []ApiResponseVdcGroupParticipatingVdc `json:"participatingOrgVdcs" fakesize:"2"`
		NetworkProviderType string                                `json:"networkProviderType" fake:"NSX_T"`
		Type                string                                `json:"type" fake:"LOCAL"`
	}
)

func (r *ApiResponseListVdcGroup) ToModel() *types.ModelListVdcGroup {
	model := &types.ModelListVdcGroup{
		VdcGroups: make([]types.ModelGetVdcGroup, 0),
	}

	for _, vdcGroup := range r.Values {
		model.VdcGroups = append(model.VdcGroups, vdcGroup.ToModel())
	}

	return model
}

func (r *ApiResponseListVdcGroupDetails) ToModel() types.ModelGetVdcGroup {
	detail := types.ModelGetVdcGroup{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, vdc := range r.Vdcs {
		detail.Vdcs = append(detail.Vdcs, vdc.ToModel())
	}

	detail.NumberOfVdcs = len(detail.Vdcs)

	return detail
}

func (r *ApiResponseVdcGroupParticipatingVdc) ToModel() types.ModelGetVdcGroupVdc {
	return types.ModelGetVdcGroupVdc{
		ID:   r.Vdc.ID,
		Name: r.Vdc.Name,
	}
}
