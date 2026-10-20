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
	APIResponseListVDCGroup struct {
		Values []APIResponseListVDCGroupDetails `json:"values" fakesize:"3"`
	}

	APIResponseListVDCGroupDetails struct {
		ID          string                                `json:"id" fake:"{urn:vdcGroup}"`
		OrgID       string                                `json:"orgId" fake:"{urn:org}"`
		Name        string                                `json:"name" fake:"mockvdcgroup-{word}"`
		Description string                                `json:"description" fake:"{sentence}"`
		Vdcs        []APIResponseVDCGroupParticipatingVDC `json:"participatingOrgVdcs" fakesize:"2"`
	}

	APIResponseVDCGroupParticipatingVDC struct {
		VDC                  APIResponseVDCGroupParticipatingVDCRef  `json:"vdcRef"`
		Site                 APIResponseVDCGroupParticipatingSiteRef `json:"siteRef"`
		FaultDomainTag       string                                  `json:"faultDomainTag" fake:"AZ01"`
		NetworkProviderScope string                                  `json:"networkProviderScope" fake:"AZ01"`
	}
	APIResponseVDCGroupParticipatingSiteRef struct {
		ID   string `json:"id,omitempty" fake:"{urn:vdc}"`
		Name string `json:"name,omitempty" fake:"mockvdc-{word}"`
	}

	APIResponseVDCGroupParticipatingVDCRef struct {
		ID   string `json:"id,omitempty" fake:"{urn:vdc}"`
		Name string `json:"name,omitempty" fake:"mockvdc-{word}"`
	}

	// * Create
	APIRequestCreateVDCGroup struct {
		OrgID               string                                `json:"orgId" fake:"{org}"`
		Name                string                                `json:"name" fake:"mockvdcgroup-{word}"`
		Description         string                                `json:"description,omitempty" fake:"{sentence}"`
		Vdcs                []APIResponseVDCGroupParticipatingVDC `json:"participatingOrgVdcs" fakesize:"2"`
		NetworkProviderType string                                `json:"networkProviderType" fake:"NSX_T"`
		Type                string                                `json:"type" fake:"LOCAL"`
	}

	// * Update
	APIRequestUpdateVDCGroup struct {
		ID                  string                                `json:"id" fake:"{urn:vdcGroup}"`
		OrgID               string                                `json:"orgId" fake:"{org}"`
		Name                string                                `json:"name" fake:"mockvdcgroup-{word}"`
		Description         string                                `json:"description,omitempty" fake:"{sentence}"`
		Vdcs                []APIResponseVDCGroupParticipatingVDC `json:"participatingOrgVdcs" fakesize:"2"`
		NetworkProviderType string                                `json:"networkProviderType" fake:"NSX_T"`
		Type                string                                `json:"type" fake:"LOCAL"`
	}
)

func (r *APIResponseListVDCGroup) ToModel() *types.ModelListVDCGroup {
	model := &types.ModelListVDCGroup{
		VDCGroups: make([]types.ModelGetVDCGroup, 0),
	}

	for _, vdcGroup := range r.Values {
		model.VDCGroups = append(model.VDCGroups, vdcGroup.ToModel())
	}

	return model
}

func (r *APIResponseListVDCGroupDetails) ToModel() types.ModelGetVDCGroup {
	detail := types.ModelGetVDCGroup{
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

func (r *APIResponseVDCGroupParticipatingVDC) ToModel() types.ModelGetVDCGroupVDC {
	return types.ModelGetVDCGroupVDC{
		ID:   r.VDC.ID,
		Name: r.VDC.Name,
	}
}

func (r *APIRequestCreateVDCGroup) ToModel() types.ModelGetVDCGroup {
	model := types.ModelGetVDCGroup{
		Name:        r.Name,
		Description: r.Description,
	}

	for _, vdc := range r.Vdcs {
		model.Vdcs = append(model.Vdcs, vdc.ToModel())
	}

	model.NumberOfVdcs = len(model.Vdcs)

	return model
}

func (r *APIRequestUpdateVDCGroup) ToModel() types.ModelGetVDCGroup {
	model := types.ModelGetVDCGroup{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, vdc := range r.Vdcs {
		model.Vdcs = append(model.Vdcs, vdc.ToModel())
	}

	model.NumberOfVdcs = len(model.Vdcs)

	return model
}
