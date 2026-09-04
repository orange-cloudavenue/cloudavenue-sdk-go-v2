/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcgroup

import (
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestCreateAppPortProfile(t *testing.T) {
	createdID := generator.MustGenerate("{urn:applicationPortProfile}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	orgID := generator.MustGenerate("{urn:org}")
	ports := []types.ParamsAppPortProfilePort{{
		Protocol:         types.AppPortProfileProtocolTCP,
		DestinationPorts: []string{"443", "8443-8444"},
	}}

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{
			ID:   vdcGroupID,
			Name: vdcGroupName,
		}},
	}, nil)

	ms.CleanResponse(endpoints.CreateAppPortProfile())
	ms.SetResponse(endpoints.CreateAppPortProfile(), &itypes.ApiResponseAppPortProfile{
		ID:          createdID,
		Name:        "app-1",
		Description: "desc",
		Scope:       types.AppPortProfileScopeTenant,
		OrgRef:      &itypes.ApiObjectReference{ID: orgID},
		ApplicationPorts: []itypes.ApiAppPortProfilePort{{
			Protocol:         types.AppPortProfileProtocolTCP,
			DestinationPorts: []string{"443", "8443-8444"},
		}},
	}, nil)

	resp, err := client.CreateAppPortProfile(t.Context(), types.ParamsCreateAppPortProfile{
		Name:             "app-1",
		Description:      "desc",
		VdcGroupName:     vdcGroupName,
		ApplicationPorts: ports,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createdID, resp.ID)
	assert.Equal(t, "app-1", resp.Name)
	assert.Equal(t, types.AppPortProfileScopeTenant, resp.Scope)
	assert.Equal(t, orgID, resp.OrgID)
	assert.Len(t, resp.ApplicationPorts, 1)
	assert.Equal(t, types.AppPortProfileProtocolTCP, resp.ApplicationPorts[0].Protocol)
	assert.Equal(t, []string{"443", "8443-8444"}, resp.ApplicationPorts[0].DestinationPorts)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.CreateAppPortProfile())
}

func TestListAppPortProfile(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	profileID := generator.MustGenerate("{urn:applicationPortProfile}")
	orgID := generator.MustGenerate("{urn:org}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.ListAppPortProfile())
	ms.SetResponse(endpoints.ListAppPortProfile(), &itypes.ApiResponseListAppPortProfile{
		Values: []itypes.ApiResponseAppPortProfile{{
			ID:     profileID,
			Name:   "app-1",
			Scope:  types.AppPortProfileScopeTenant,
			OrgRef: &itypes.ApiObjectReference{ID: orgID},
			ApplicationPorts: []itypes.ApiAppPortProfilePort{{
				Protocol:         types.AppPortProfileProtocolTCP,
				DestinationPorts: []string{"443"},
			}},
		}},
	}, nil)

	resp, err := client.ListAppPortProfile(t.Context(), types.ParamsListAppPortProfile{VdcGroupName: vdcGroupName})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.AppPortProfiles, 1)
	assert.Equal(t, profileID, resp.AppPortProfiles[0].ID)
	assert.Equal(t, orgID, resp.AppPortProfiles[0].OrgID)
	assert.Equal(t, []string{"443"}, resp.AppPortProfiles[0].ApplicationPorts[0].DestinationPorts)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.ListAppPortProfile())
}

func TestGetAppPortProfile(t *testing.T) {
	profileID := generator.MustGenerate("{urn:applicationPortProfile}")
	orgID := generator.MustGenerate("{urn:org}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetAppPortProfile())
	ms.SetResponse(endpoints.GetAppPortProfile(), &itypes.ApiResponseAppPortProfile{
		ID:     profileID,
		Name:   "app-1",
		Scope:  types.AppPortProfileScopeTenant,
		OrgRef: &itypes.ApiObjectReference{ID: orgID},
		ApplicationPorts: []itypes.ApiAppPortProfilePort{{
			Protocol:         types.AppPortProfileProtocolUDP,
			DestinationPorts: []string{"53"},
		}},
	}, nil)

	resp, err := client.GetAppPortProfile(t.Context(), types.ParamsGetAppPortProfile{ID: profileID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, profileID, resp.ID)
	assert.Equal(t, orgID, resp.OrgID)
	assert.Equal(t, types.AppPortProfileProtocolUDP, resp.ApplicationPorts[0].Protocol)

	ms.CleanResponse(endpoints.GetAppPortProfile())
}

func TestUpdateAppPortProfile(t *testing.T) {
	profileID := generator.MustGenerate("{urn:applicationPortProfile}")
	orgID := generator.MustGenerate("{urn:org}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetAppPortProfile())
	ms.SetResponse(endpoints.GetAppPortProfile(), &itypes.ApiResponseAppPortProfile{
		ID:          profileID,
		Name:        "app-1",
		Description: "old-desc",
		Scope:       types.AppPortProfileScopeTenant,
		OrgRef:      &itypes.ApiObjectReference{ID: orgID},
		ApplicationPorts: []itypes.ApiAppPortProfilePort{{
			Protocol:         types.AppPortProfileProtocolTCP,
			DestinationPorts: []string{"443"},
		}},
	}, nil)
	ms.CleanResponse(endpoints.UpdateAppPortProfile())

	resp, err := client.UpdateAppPortProfile(t.Context(), types.ParamsUpdateAppPortProfile{
		ID:          profileID,
		Description: "new-desc",
		ApplicationPorts: []types.ParamsAppPortProfilePort{{
			Protocol:         types.AppPortProfileProtocolUDP,
			DestinationPorts: []string{"53"},
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, profileID, resp.ID)
	assert.Equal(t, "app-1", resp.Name)
	assert.Equal(t, "new-desc", resp.Description)
	assert.Equal(t, types.AppPortProfileScopeTenant, resp.Scope)
	assert.Equal(t, orgID, resp.OrgID)
	assert.Len(t, resp.ApplicationPorts, 1)
	assert.Equal(t, types.AppPortProfileProtocolUDP, resp.ApplicationPorts[0].Protocol)
	assert.Equal(t, []string{"53"}, resp.ApplicationPorts[0].DestinationPorts)

	ms.CleanResponse(endpoints.GetAppPortProfile())
	ms.CleanResponse(endpoints.UpdateAppPortProfile())
}

func TestDeleteAppPortProfile(t *testing.T) {
	profileID := generator.MustGenerate("{urn:applicationPortProfile}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetAppPortProfile())
	ms.SetResponse(endpoints.GetAppPortProfile(), &itypes.ApiResponseAppPortProfile{
		ID:              profileID,
		Name:            "app-1",
		ContextEntityId: vdcGroupID,
		Scope:           types.AppPortProfileScopeTenant,
	}, nil)
	ms.CleanResponse(endpoints.DeleteAppPortProfile())

	err := client.DeleteAppPortProfile(t.Context(), types.ParamsDeleteAppPortProfile{
		ID: profileID,
	})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetAppPortProfile())
	ms.CleanResponse(endpoints.DeleteAppPortProfile())
}
