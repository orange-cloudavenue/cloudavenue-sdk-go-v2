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

func TestCreateNetworkContextProfile(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	orgID := generator.MustGenerate("{urn:org}")
	profileID := generator.MustGenerate("{uuid}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.ListNetworkContextProfile())
	ms.SetResponse(endpoints.ListNetworkContextProfile(), &itypes.APIResponseListNetworkContextProfile{
		Values: []itypes.APIResponseNetworkContextProfile{{
			ID:              profileID,
			Name:            "ncp-1",
			Description:     "desc",
			Scope:           types.NetworkContextProfileScopeTenant,
			ContextEntityID: vdcGroupID,
			OrgRef:          &itypes.APIObjectReference{ID: orgID},
			Attributes: []itypes.APINetworkContextProfileAttribute{{
				Type:   types.NetworkContextProfileAttributeTypeAppID,
				Values: []string{"HTTP"},
				SubAttributes: []itypes.APINetworkContextProfileSubAttribute{{
					Type:   types.NetworkContextProfileSubAttributeTypeTLSVersion,
					Values: []string{"TLS_V13"},
				}},
			}},
		}},
	}, nil)

	resp, err := client.CreateNetworkContextProfile(t.Context(), types.ParamsCreateNetworkContextProfile{
		Name:         "ncp-1",
		Description:  "desc",
		VDCGroupName: vdcGroupName,
		Attributes: []types.ParamsNetworkContextProfileAttribute{{
			Type:   types.NetworkContextProfileAttributeTypeAppID,
			Values: []string{"HTTP"},
			SubAttributes: []types.ParamsNetworkContextProfileSubAttribute{{
				Type:   types.NetworkContextProfileSubAttributeTypeTLSVersion,
				Values: []string{"TLS_V13"},
			}},
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, profileID, resp.ID)
	assert.Equal(t, "ncp-1", resp.Name)
	assert.Equal(t, types.NetworkContextProfileScopeTenant, resp.Scope)
	assert.Equal(t, orgID, resp.OrgID)
	assert.Len(t, resp.Attributes, 1)
	assert.Equal(t, types.NetworkContextProfileAttributeTypeAppID, resp.Attributes[0].Type)
	assert.Equal(t, []string{"HTTP"}, resp.Attributes[0].Values)
	assert.Len(t, resp.Attributes[0].SubAttributes, 1)
	assert.Equal(t, types.NetworkContextProfileSubAttributeTypeTLSVersion, resp.Attributes[0].SubAttributes[0].Type)
	assert.Equal(t, []string{"TLS_V13"}, resp.Attributes[0].SubAttributes[0].Values)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.ListNetworkContextProfile())
}

func TestListNetworkContextProfile(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	orgID := generator.MustGenerate("{urn:org}")
	profileID := generator.MustGenerate("{uuid}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.ListNetworkContextProfile())
	ms.SetResponse(endpoints.ListNetworkContextProfile(), &itypes.APIResponseListNetworkContextProfile{
		Values: []itypes.APIResponseNetworkContextProfile{{
			ID:              profileID,
			Name:            "ncp-1",
			Scope:           types.NetworkContextProfileScopeTenant,
			ContextEntityID: vdcGroupID,
			OrgRef:          &itypes.APIObjectReference{ID: orgID},
			Attributes: []itypes.APINetworkContextProfileAttribute{{
				Type:   types.NetworkContextProfileAttributeTypeAppID,
				Values: []string{"HTTP"},
			}},
		}},
	}, nil)

	resp, err := client.ListNetworkContextProfile(t.Context(), types.ParamsListNetworkContextProfile{VDCGroupName: vdcGroupName})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.NetworkContextProfiles, 1)
	assert.Equal(t, profileID, resp.NetworkContextProfiles[0].ID)
	assert.Equal(t, orgID, resp.NetworkContextProfiles[0].OrgID)
	assert.Equal(t, types.NetworkContextProfileAttributeTypeAppID, resp.NetworkContextProfiles[0].Attributes[0].Type)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.ListNetworkContextProfile())
}

func TestGetNetworkContextProfile(t *testing.T) {
	profileID := generator.MustGenerate("{uuid}")
	orgID := generator.MustGenerate("{urn:org}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
	ms.SetResponse(endpoints.GetNetworkContextProfile(), &itypes.APIResponseNetworkContextProfile{
		ID:     profileID,
		Name:   "ncp-1",
		Scope:  types.NetworkContextProfileScopeTenant,
		OrgRef: &itypes.APIObjectReference{ID: orgID},
		Attributes: []itypes.APINetworkContextProfileAttribute{{
			Type:   types.NetworkContextProfileAttributeTypeDomainName,
			Values: []string{"example.com"},
		}},
	}, nil)

	resp, err := client.GetNetworkContextProfile(t.Context(), types.ParamsGetNetworkContextProfile{ID: profileID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, profileID, resp.ID)
	assert.Equal(t, orgID, resp.OrgID)
	assert.Equal(t, types.NetworkContextProfileAttributeTypeDomainName, resp.Attributes[0].Type)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
}

func TestUpdateNetworkContextProfile(t *testing.T) {
	profileID := generator.MustGenerate("{uuid}")
	orgID := generator.MustGenerate("{urn:org}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
	ms.SetResponse(endpoints.GetNetworkContextProfile(), &itypes.APIResponseNetworkContextProfile{
		ID:              profileID,
		Name:            "ncp-1",
		Description:     "old-desc",
		Scope:           types.NetworkContextProfileScopeTenant,
		ContextEntityID: vdcGroupID,
		OrgRef:          &itypes.APIObjectReference{ID: orgID},
		Attributes: []itypes.APINetworkContextProfileAttribute{{
			Type:   types.NetworkContextProfileAttributeTypeAppID,
			Values: []string{"HTTP"},
		}},
	}, nil)
	ms.CleanResponse(endpoints.UpdateNetworkContextProfile())

	resp, err := client.UpdateNetworkContextProfile(t.Context(), types.ParamsUpdateNetworkContextProfile{
		ID:          profileID,
		Description: "new-desc",
		Attributes: []types.ParamsNetworkContextProfileAttribute{{
			Type:   types.NetworkContextProfileAttributeTypeDomainName,
			Values: []string{"example.com"},
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, profileID, resp.ID)
	assert.Equal(t, "ncp-1", resp.Name)
	assert.Equal(t, "new-desc", resp.Description)
	assert.Equal(t, types.NetworkContextProfileScopeTenant, resp.Scope)
	assert.Equal(t, orgID, resp.OrgID)
	assert.Len(t, resp.Attributes, 1)
	assert.Equal(t, types.NetworkContextProfileAttributeTypeDomainName, resp.Attributes[0].Type)
	assert.Equal(t, []string{"example.com"}, resp.Attributes[0].Values)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
	ms.CleanResponse(endpoints.UpdateNetworkContextProfile())
}

func TestDeleteNetworkContextProfile(t *testing.T) {
	profileID := generator.MustGenerate("{uuid}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
	ms.SetResponse(endpoints.GetNetworkContextProfile(), &itypes.APIResponseNetworkContextProfile{
		ID:              profileID,
		Name:            "ncp-1",
		ContextEntityID: vdcGroupID,
		Scope:           types.NetworkContextProfileScopeTenant,
	}, nil)
	ms.CleanResponse(endpoints.DeleteNetworkContextProfile())

	err := client.DeleteNetworkContextProfile(t.Context(), types.ParamsDeleteNetworkContextProfile{ID: profileID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
	ms.CleanResponse(endpoints.DeleteNetworkContextProfile())
}

func TestGetNetworkContextProfileAttributes(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.GetNetworkContextProfileAttributes())
	ms.SetResponse(endpoints.GetNetworkContextProfileAttributes(), &itypes.APINetworkContextProfileAttributesResponse{
		Attributes: []itypes.APINetworkContextProfileAttribute{
			{
				Type:   types.NetworkContextProfileAttributeTypeAppID,
				Values: []string{"HTTP", "DNS"},
			},
			{
				Type:   types.NetworkContextProfileAttributeTypeDomainName,
				Values: []string{"example.com", "orange.com"},
			},
		},
	}, nil)

	resp, err := client.GetNetworkContextProfileAttributes(t.Context(), types.ParamsGetNetworkContextProfileAttributes{
		VDCGroupName: vdcGroupName,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, []string{"HTTP", "DNS"}, resp.AppIDValues)
	assert.Equal(t, []string{"example.com", "orange.com"}, resp.DomainNameValues)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.GetNetworkContextProfileAttributes())
}
