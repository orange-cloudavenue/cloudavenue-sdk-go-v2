/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestCreateNetworkContextProfile(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	orgID := generator.MustGenerate("{urn:org}")
	profileID := generator.MustGenerate("{uuid}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.ApiResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.ListNetworkContextProfile())
	ms.SetResponseFunc(endpoints.ListNetworkContextProfile(), func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Query().Get("filter"), "vdcGroupId=="+vdcGroupID)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cloud-Avenue-Mock", "true")
		assert.NoError(t, json.NewEncoder(w).Encode(&itypes.ApiResponseListNetworkContextProfile{
			Values: []itypes.ApiResponseNetworkContextProfile{{
				ID:              profileID,
				Name:            "ncp-1",
				Description:     "desc",
				Scope:           types.NetworkContextProfileScopeTenant,
				ContextEntityId: vdcGroupID,
				OrgRef:          &itypes.ApiObjectReference{ID: orgID},
				Attributes: []itypes.ApiNetworkContextProfileAttribute{{
					Type:   types.NetworkContextProfileAttributeTypeAppID,
					Values: []string{"HTTP"},
				}},
			}},
		}))
	})

	resp, err := client.CreateNetworkContextProfile(t.Context(), types.ParamsCreateEdgeGatewayNetworkContextProfile{
		EdgeGatewayID: edgeGatewayID,
		Name:          "ncp-1",
		Description:   "desc",
		Attributes: []types.ParamsNetworkContextProfileAttribute{{
			Type:   types.NetworkContextProfileAttributeTypeAppID,
			Values: []string{"HTTP"},
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, profileID, resp.ID)
	assert.Equal(t, orgID, resp.OrgID)
	assert.Equal(t, types.NetworkContextProfileScopeTenant, resp.Scope)
	assert.Len(t, resp.Attributes, 1)
	assert.Equal(t, types.NetworkContextProfileAttributeTypeAppID, resp.Attributes[0].Type)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.ListNetworkContextProfile())
}

func TestListNetworkContextProfileUsesVdcScope(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcID := generator.MustGenerate("{urn:vdc}")
	orgID := generator.MustGenerate("{urn:org}")
	profileID := generator.MustGenerate("{uuid}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.ApiResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.ApiObjectReference{ID: vdcID, Name: "vdc-1"},
	}, nil)

	ms.CleanResponse(endpoints.ListNetworkContextProfile())
	ms.SetResponseFunc(endpoints.ListNetworkContextProfile(), func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Query().Get("filter"), "orgVdcId=="+vdcID)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cloud-Avenue-Mock", "true")
		assert.NoError(t, json.NewEncoder(w).Encode(&itypes.ApiResponseListNetworkContextProfile{
			Values: []itypes.ApiResponseNetworkContextProfile{{
				ID:     profileID,
				Name:   "ncp-1",
				Scope:  types.NetworkContextProfileScopeTenant,
				OrgRef: &itypes.ApiObjectReference{ID: orgID},
				Attributes: []itypes.ApiNetworkContextProfileAttribute{{
					Type:   types.NetworkContextProfileAttributeTypeDomainName,
					Values: []string{"example.com"},
				}},
			}},
		}))
	})

	resp, err := client.ListNetworkContextProfile(t.Context(), types.ParamsListEdgeGatewayNetworkContextProfile{EdgeGatewayID: edgeGatewayID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.NetworkContextProfiles, 1)
	assert.Equal(t, profileID, resp.NetworkContextProfiles[0].ID)
	assert.Equal(t, orgID, resp.NetworkContextProfiles[0].OrgID)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.ListNetworkContextProfile())
}

func TestGetNetworkContextProfileByName(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	profileID := generator.MustGenerate("{uuid}")
	orgID := generator.MustGenerate("{urn:org}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.ApiResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.ListNetworkContextProfile())
	ms.SetResponse(endpoints.ListNetworkContextProfile(), &itypes.ApiResponseListNetworkContextProfile{
		Values: []itypes.ApiResponseNetworkContextProfile{{
			ID:     profileID,
			Name:   "ncp-1",
			Scope:  types.NetworkContextProfileScopeTenant,
			OrgRef: &itypes.ApiObjectReference{ID: orgID},
			Attributes: []itypes.ApiNetworkContextProfileAttribute{{
				Type:   types.NetworkContextProfileAttributeTypeDomainName,
				Values: []string{"example.com"},
			}},
		}},
	}, nil)

	resp, err := client.GetNetworkContextProfile(t.Context(), types.ParamsGetEdgeGatewayNetworkContextProfile{
		Name:          "ncp-1",
		EdgeGatewayID: edgeGatewayID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, profileID, resp.ID)
	assert.Equal(t, orgID, resp.OrgID)
	assert.Equal(t, types.NetworkContextProfileAttributeTypeDomainName, resp.Attributes[0].Type)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.ListNetworkContextProfile())
}

func TestGetNetworkContextProfileAttributes(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.ApiResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.GetNetworkContextProfileAttributes())
	ms.SetResponseFunc(endpoints.GetNetworkContextProfileAttributes(), func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Query().Get("filter"), "vdcGroupId=="+vdcGroupID)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cloud-Avenue-Mock", "true")
		assert.NoError(t, json.NewEncoder(w).Encode(&itypes.ApiNetworkContextProfileAttributesResponse{
			Attributes: []itypes.ApiNetworkContextProfileAttribute{
				{Type: types.NetworkContextProfileAttributeTypeAppID, Values: []string{"HTTP", "DNS"}},
				{Type: types.NetworkContextProfileAttributeTypeDomainName, Values: []string{"example.com", "orange.com"}},
			},
		}))
	})

	resp, err := client.GetNetworkContextProfileAttributes(t.Context(), types.ParamsGetEdgeGatewayNetworkContextProfileAttributes{EdgeGatewayID: edgeGatewayID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, []string{"HTTP", "DNS"}, resp.AppIDValues)
	assert.Equal(t, []string{"example.com", "orange.com"}, resp.DomainNameValues)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.GetNetworkContextProfileAttributes())
}

func TestUpdateNetworkContextProfile(t *testing.T) {
	profileID := generator.MustGenerate("{uuid}")
	orgID := generator.MustGenerate("{urn:org}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
	ms.SetResponse(endpoints.GetNetworkContextProfile(), &itypes.ApiResponseNetworkContextProfile{
		ID:              profileID,
		Name:            "ncp-1",
		Description:     "old-desc",
		Scope:           types.NetworkContextProfileScopeTenant,
		ContextEntityId: vdcGroupID,
		OrgRef:          &itypes.ApiObjectReference{ID: orgID},
		Attributes: []itypes.ApiNetworkContextProfileAttribute{{
			Type:   types.NetworkContextProfileAttributeTypeAppID,
			Values: []string{"HTTP"},
		}},
	}, nil)
	ms.CleanResponse(endpoints.UpdateNetworkContextProfile())

	resp, err := client.UpdateNetworkContextProfile(t.Context(), types.ParamsUpdateEdgeGatewayNetworkContextProfile{
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

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
	ms.SetResponse(endpoints.GetNetworkContextProfile(), &itypes.ApiResponseNetworkContextProfile{
		ID:              profileID,
		Name:            "ncp-1",
		ContextEntityId: generator.MustGenerate("{urn:vdcGroup}"),
		Scope:           types.NetworkContextProfileScopeTenant,
	}, nil)
	ms.CleanResponse(endpoints.DeleteNetworkContextProfile())

	err := client.DeleteNetworkContextProfile(t.Context(), types.ParamsDeleteEdgeGatewayNetworkContextProfile{ID: profileID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetNetworkContextProfile())
	ms.CleanResponse(endpoints.DeleteNetworkContextProfile())
}
