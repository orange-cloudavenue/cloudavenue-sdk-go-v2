/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestListVApp(t *testing.T) {
	tests := []struct {
		name               string
		vdcID              string
		mockResponse       any
		mockResponseStatus int

		expectedErr bool
	}{
		{
			name:  "List VApps with valid VDC ID",
			vdcID: "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012",
			mockResponse: &itypes.APIResponseListVApp{
				Records: []itypes.APIResponseListVAppRecord{
					{
						ID:   "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
						Name: "test-vapp-1",
					},
				},
			},
			mockResponseStatus: 200,
			expectedErr:        false,
		},
		{
			name:        "List VApps with empty VDC ID",
			vdcID:       "",
			expectedErr: true,
		},
		{
			name:  "Error 401 Unauthorized",
			vdcID: "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012",
			mockResponse: &itypes.APIResponseListVApp{
				Records: []itypes.APIResponseListVAppRecord{},
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVApp())
				ms.SetResponse(endpoints.ListVApp(), tt.mockResponse, &tt.mockResponseStatus)
			}

			resp, err := client.ListVApp(t.Context(), tt.vdcID)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Expected response to be not nil")
		})
	}
}

func TestGetVApp(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsGetVApp
		mockResponse       any
		mockResponseStatus int

		mockListResponse       any
		mockListResponseStatus int

		expectedErr bool
	}{
		{
			name: "Get VApp by ID",
			params: types.ParamsGetVApp{
				ID: "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
			},
			mockResponse: &itypes.APIResponseGetVApp{
				ID:   "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
				Name: "test-vapp-1",
			},
			mockResponseStatus: 200,
			expectedErr:        false,
		},
		{
			name: "Get VApp by Name",
			params: types.ParamsGetVApp{
				Name: "test-vapp-1",
			},
			mockListResponse: &itypes.APIResponseListVApp{
				Records: []itypes.APIResponseListVAppRecord{
					{
						ID:   "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
						Name: "test-vapp-1",
					},
				},
			},
			mockListResponseStatus: 200,
			expectedErr:            false,
		},
		{
			name: "Error 401 Unauthorized",
			params: types.ParamsGetVApp{
				ID: "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
		{
			name: "VApp not found by name",
			params: types.ParamsGetVApp{
				Name: "test-vapp-1",
			},
			mockListResponse:       &itypes.APIResponseListVApp{Records: []itypes.APIResponseListVAppRecord{}},
			mockListResponseStatus: 200,
			expectedErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.GetVApp())
				ms.SetResponse(endpoints.GetVApp(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListResponse != nil || tt.mockListResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVApp())
				ms.SetResponse(endpoints.ListVApp(), tt.mockListResponse, &tt.mockListResponseStatus)
			}

			resp, err := client.GetVApp(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Expected response to be not nil")
		})
	}
}

func TestVAppCreationEndpoints(t *testing.T) {
	assert.Equal(t, "/api/vdc/{vdc-id}/action/composeVApp", endpoints.ComposeVApp().PathTemplate)
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-ComposeVApp.html", endpoints.ComposeVApp().DocumentationURL)
	assert.Equal(t, "/api/vdc/{vdc-id}/action/instantiateVAppTemplate", endpoints.InstantiateVAppTemplate().PathTemplate)
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-InstantiateVAppTemplate.html", endpoints.InstantiateVAppTemplate().DocumentationURL)
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-NetworkConfigSection-vApp.html", endpoints.RemoveAllNetworks().DocumentationURL)
	assert.Equal(t, "/api/vapp/{vapp-id}/networkConfigSection/", endpoints.GetVAppNetworkConfigSection().PathTemplate)
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-NetworkConfigSection-vApp.html", endpoints.GetVAppNetworkConfigSection().DocumentationURL)
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-NetworkConfigSection-vApp.html", endpoints.UpdateVAppNetworkConfigSection().DocumentationURL)
	assert.Equal(t, "/api/vapp/{vapp-id}/productSections", endpoints.GetVAppProductSections().PathTemplate)
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/GET-ProductSections-vApp.html", endpoints.GetVAppProductSections().DocumentationURL)
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-ProductSections-vApp.html", endpoints.UpdateVAppProductSections().DocumentationURL)
}

func TestGetVAppNetworkConfig(t *testing.T) {
	client, ms := newClient(t)
	vappID := "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012"

	ms.CleanResponse(endpoints.GetVAppNetworkConfigSection())
	ms.SetResponseFunc(endpoints.GetVAppNetworkConfigSection(), func(w http.ResponseWriter, r *http.Request) {
		connected := true
		w.Header().Set("Content-Type", "application/xml")
		assert.NoError(t, xml.NewEncoder(w).Encode(&itypes.APIResponseVAppNetworkConfigSection{
			XMLNS: "http://www.vmware.com/vcloud/v1.5",
			OVF:   "http://schemas.dmtf.org/ovf/envelope/1",
			Info:  "Runtime network config",
			NetworkConfig: []itypes.APIResponseVAppNetworkConfig{{
				NetworkID:   "urn:vcloud:vappNetwork:isolated-net",
				NetworkHREF: "https://example.com/api/vAppNetwork/isolated-net",
				Description: "network description",
				NetworkType: "application/vnd.vmware.vcloud.vAppNetwork+xml",
				NetworkName: "isolated-net",
				Configuration: itypes.APIResponseVAppRuntimeConfiguration{
					FenceMode:   "isolated",
					IsConnected: &connected,
					IPScopes: &itypes.APIResponseVAppRuntimeIPScopes{
						IPScope: []itypes.APIResponseVAppRuntimeIPScope{{
							Gateway:   "192.168.1.1",
							Netmask:   "255.255.255.0",
							DNS1:      "8.8.8.8",
							DNSSuffix: "example.org",
							IPRanges:  itypes.APIResponseVAppRuntimeIPRanges{IPRange: []itypes.APIResponseVAppRuntimeIPRange{{StartAddress: "192.168.1.10", EndAddress: "192.168.1.20"}}},
						}},
					},
				},
			}},
		}))
	})

	resp, err := client.GetVAppNetworkConfig(t.Context(), types.ParamsGetVAppNetworkConfig{ID: vappID})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.NetworkConfig, 1) {
		assert.Equal(t, "urn:vcloud:vappNetwork:isolated-net", resp.NetworkConfig[0].ID)
		assert.Equal(t, "https://example.com/api/vAppNetwork/isolated-net", resp.NetworkConfig[0].HREF)
		assert.Equal(t, "application/vnd.vmware.vcloud.vAppNetwork+xml", resp.NetworkConfig[0].Type)
		assert.Equal(t, "network description", resp.NetworkConfig[0].Description)
		assert.Equal(t, "isolated-net", resp.NetworkConfig[0].NetworkName)
		assert.Equal(t, "isolated", resp.NetworkConfig[0].Configuration.FenceMode)
	}

	ms.CleanResponse(endpoints.GetVAppNetworkConfigSection())
}

func TestRuntimeVAppNetworkMutations(t *testing.T) {
	tests := []struct {
		name        string
		call        func(*Client) (*types.ModelVAppRuntimeNetworkConfigSection, error)
		expectPUT   func(t *testing.T, body *itypes.APIResponseVAppNetworkConfigSection)
		expectedErr string
	}{
		{
			name: "create isolated",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				return client.CreateVAppNetwork(t.Context(), types.ParamsCreateVAppNetwork{VAppID: testVAppID, Network: isolatedRuntimeNetwork("isolated-net")})
			},
			expectPUT: func(t *testing.T, body *itypes.APIResponseVAppNetworkConfigSection) {
				if assert.Len(t, body.NetworkConfig, 3) {
					assert.Equal(t, "urn:vcloud:vappNetwork:existing-net", body.NetworkConfig[0].NetworkID)
					assert.Equal(t, "urn:vcloud:vappNetwork:sibling-net", body.NetworkConfig[1].NetworkID)
					assert.Equal(t, "isolated-net", body.NetworkConfig[2].NetworkName)
					assert.Equal(t, "isolated", body.NetworkConfig[2].Configuration.FenceMode)
				}
			},
		},
		{
			name: "add org network",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				network := isolatedRuntimeNetwork("org-net")
				network.Configuration.FenceMode = "bridged"
				network.Configuration.ParentNetwork = &types.VAppReference{HREF: "https://example.com/api/network/1234", ID: "urn:vcloud:network:1234", Name: "org-net"}
				network.Configuration.IPScopes = nil
				return client.AddOrgNetworkToVApp(t.Context(), types.ParamsAddOrgNetworkToVApp{VAppID: testVAppID, Network: network})
			},
			expectPUT: func(t *testing.T, body *itypes.APIResponseVAppNetworkConfigSection) {
				if assert.Len(t, body.NetworkConfig, 3) {
					assert.Equal(t, "org-net", body.NetworkConfig[2].NetworkName)
					if assert.NotNil(t, body.NetworkConfig[2].Configuration.ParentNetwork) {
						assert.Equal(t, "urn:vcloud:network:1234", body.NetworkConfig[2].Configuration.ParentNetwork.ID)
					}
					assert.Nil(t, body.NetworkConfig[2].Configuration.IPScopes)
				}
			},
		},
		{
			name: "create with parent network fails validation",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				network := isolatedRuntimeNetwork("invalid-create")
				network.Configuration.ParentNetwork = &types.VAppReference{ID: "urn:vcloud:network:1234"}
				return client.CreateVAppNetwork(t.Context(), types.ParamsCreateVAppNetwork{VAppID: testVAppID, Network: network})
			},
			expectedErr: "parent network is not allowed for create",
		},
		{
			name: "create non isolated fails validation",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				network := isolatedRuntimeNetwork("invalid-fence")
				network.Configuration.FenceMode = "bridged"
				return client.CreateVAppNetwork(t.Context(), types.ParamsCreateVAppNetwork{VAppID: testVAppID, Network: network})
			},
			expectedErr: "create requires isolated fence mode",
		},
		{
			name: "update existing",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				network := isolatedRuntimeNetwork("existing-net")
				network.Configuration.IPScopes[0].Gateway = "192.168.2.1"
				return client.UpdateVAppNetwork(t.Context(), types.ParamsUpdateVAppNetwork{VAppID: testVAppID, NetworkName: "existing-net", Network: network})
			},
			expectPUT: func(t *testing.T, body *itypes.APIResponseVAppNetworkConfigSection) {
				assert.Equal(t, "192.168.2.1", body.NetworkConfig[0].Configuration.IPScopes.IPScope[0].Gateway)
				assert.Equal(t, "urn:vcloud:vappNetwork:existing-net", body.NetworkConfig[0].NetworkID)
				assert.Equal(t, "https://example.com/api/vAppNetwork/existing-net", body.NetworkConfig[0].NetworkHREF)
				assert.NotNil(t, body.NetworkConfig[0].Configuration.Features)
				assert.NotNil(t, body.NetworkConfig[0].Configuration.RouterInfo)
			},
		},
		{
			name: "remove existing",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				return client.RemoveVAppNetwork(t.Context(), types.ParamsRemoveVAppNetwork{VAppID: testVAppID, NetworkName: "existing-net"})
			},
			expectPUT: func(t *testing.T, body *itypes.APIResponseVAppNetworkConfigSection) {
				assert.Len(t, body.NetworkConfig, 1)
				assert.Equal(t, "sibling-net", body.NetworkConfig[0].NetworkName)
				assert.Equal(t, "urn:vcloud:vappNetwork:sibling-net", body.NetworkConfig[0].NetworkID)
			},
		},
		{
			name: "create duplicate fails",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				return client.CreateVAppNetwork(t.Context(), types.ParamsCreateVAppNetwork{VAppID: testVAppID, Network: isolatedRuntimeNetwork("existing-net")})
			},
			expectedErr: "already exists",
		},
		{
			name: "update missing fails",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				return client.UpdateVAppNetwork(t.Context(), types.ParamsUpdateVAppNetwork{VAppID: testVAppID, NetworkName: "missing-net", Network: isolatedRuntimeNetwork("missing-net")})
			},
			expectedErr: "does not exist",
		},
		{
			name: "update rename mismatch fails",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				return client.UpdateVAppNetwork(t.Context(), types.ParamsUpdateVAppNetwork{VAppID: testVAppID, NetworkName: "existing-net", Network: isolatedRuntimeNetwork("renamed-net")})
			},
			expectedErr: "network name mismatch",
		},
		{
			name: "remove missing fails",
			call: func(client *Client) (*types.ModelVAppRuntimeNetworkConfigSection, error) {
				return client.RemoveVAppNetwork(t.Context(), types.ParamsRemoveVAppNetwork{VAppID: testVAppID, NetworkName: "missing-net"})
			},
			expectedErr: "does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)
			getCalls := 0

			ms.CleanResponse(endpoints.GetVAppNetworkConfigSection())
			ms.SetResponseFunc(endpoints.GetVAppNetworkConfigSection(), func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				section := existingRuntimeSection()
				if getCalls > 0 {
					section = postMutationRuntimeSection(tt.name)
				}
				getCalls++
				assert.NoError(t, xml.NewEncoder(w).Encode(section))
			})

			ms.CleanResponse(endpoints.UpdateVAppNetworkConfigSection())
			putCalled := false
			ms.SetResponseFunc(endpoints.UpdateVAppNetworkConfigSection(), func(w http.ResponseWriter, r *http.Request) {
				putCalled = true
				var body itypes.APIResponseVAppNetworkConfigSection
				assert.Equal(t, "application/xml;version=39.1", r.Header.Get("Content-Type"))
				assert.NoError(t, xml.NewDecoder(r.Body).Decode(&body))
				if tt.expectPUT != nil {
					tt.expectPUT(t, &body)
				}
				w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
				w.WriteHeader(http.StatusAccepted)
			})

			resp, err := tt.call(client)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.False(t, putCalled)
				if tt.name == "update rename mismatch fails" || tt.name == "create with parent network fails validation" || tt.name == "create non isolated fails validation" {
					assert.Equal(t, 0, getCalls)
				} else {
					assert.Equal(t, 1, getCalls)
				}
				return
			}
			assert.NoError(t, err)
			assert.True(t, putCalled)
			assert.Equal(t, 2, getCalls)
			assert.NotNil(t, resp)
		})
	}
}

func TestRuntimeNetworkFalseFidelity(t *testing.T) {
	connected := false
	inherited := false
	section := &itypes.APIResponseVAppNetworkConfigSection{
		XMLNS: "http://www.vmware.com/vcloud/v1.5",
		OVF:   "http://schemas.dmtf.org/ovf/envelope/1",
		Info:  "Runtime network config",
		NetworkConfig: []itypes.APIResponseVAppNetworkConfig{{
			NetworkName: "existing-net",
			Configuration: itypes.APIResponseVAppRuntimeConfiguration{
				FenceMode:   "isolated",
				IsConnected: &connected,
				IPScopes: &itypes.APIResponseVAppRuntimeIPScopes{IPScope: []itypes.APIResponseVAppRuntimeIPScope{{
					IsInherited: &inherited,
					Gateway:     "192.168.1.1",
					Netmask:     "255.255.255.0",
				}}},
			},
		}},
	}

	model := section.ToModel()
	assert.False(t, model.NetworkConfig[0].Configuration.IsConnected)
	roundTrip := itypes.ToAPIResponseVAppNetworkConfigSection(model)
	if assert.NotNil(t, roundTrip.NetworkConfig[0].Configuration.IsConnected) {
		assert.False(t, *roundTrip.NetworkConfig[0].Configuration.IsConnected)
	}
	if assert.NotNil(t, roundTrip.NetworkConfig[0].Configuration.IPScopes) && assert.Len(t, roundTrip.NetworkConfig[0].Configuration.IPScopes.IPScope, 1) {
		if assert.NotNil(t, roundTrip.NetworkConfig[0].Configuration.IPScopes.IPScope[0].IsInherited) {
			assert.False(t, *roundTrip.NetworkConfig[0].Configuration.IPScopes.IPScope[0].IsInherited)
		}
	}
}

func TestRuntimeVAppNetworkValidation(t *testing.T) {
	client, _ := newClient(t)
	_, err := client.CreateVAppNetwork(t.Context(), types.ParamsCreateVAppNetwork{})
	assert.Error(t, err)
	_, err = client.AddOrgNetworkToVApp(t.Context(), types.ParamsAddOrgNetworkToVApp{VAppID: testVAppID, Network: isolatedRuntimeNetwork("org-net")})
	assert.Error(t, err)
	_, err = client.UpdateVAppNetwork(t.Context(), types.ParamsUpdateVAppNetwork{VAppID: testVAppID, Network: isolatedRuntimeNetwork("")})
	assert.Error(t, err)
	invalidNetwork := isolatedRuntimeNetwork("missing-fence")
	invalidNetwork.Configuration.FenceMode = ""
	_, err = client.UpdateVAppNetwork(t.Context(), types.ParamsUpdateVAppNetwork{VAppID: testVAppID, NetworkName: "missing-fence", Network: invalidNetwork})
	assert.Error(t, err)
	_, err = client.RemoveVAppNetwork(t.Context(), types.ParamsRemoveVAppNetwork{VAppID: testVAppID})
	assert.Error(t, err)
}

func TestVAppGuestProperties(t *testing.T) {
	client, ms := newClient(t)
	vappID := testVAppID
	getCalls := 0

	guestProperties := &itypes.APIResponseGuestProperties{ProductSectionList: []itypes.APIResponseProductSection{{Class: "vami", Instance: "default", Property: []itypes.APIResponseProductProperty{{Key: "hostname", Value: "vm-01"}, {Key: "role", Value: "app"}}}}}

	ms.CleanResponse(endpoints.GetVAppProductSections())
	ms.SetResponseFunc(endpoints.GetVAppProductSections(), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(guestProperties))
		getCalls++
	})

	resp, err := client.GetVAppGuestProperties(t.Context(), types.ParamsGetVAppGuestProperties{ID: vappID})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Sections, 1) {
		assert.Equal(t, "vami", resp.Sections[0].Class)
		assert.Equal(t, "default", resp.Sections[0].Instance)
		assert.Equal(t, "vm-01", resp.Sections[0].Properties["hostname"])
	}

	ms.CleanResponse(endpoints.UpdateVAppProductSections())
	ms.SetResponseFunc(endpoints.UpdateVAppProductSections(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.APIResponseGuestProperties
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		if assert.Len(t, body.ProductSectionList, 1) {
			assert.Equal(t, "vami", body.ProductSectionList[0].Class)
			assert.Len(t, body.ProductSectionList[0].Property, 2)
			assert.Equal(t, []itypes.APIResponseProductProperty{{Key: "hostname", Value: "vm-02"}, {Key: "role", Value: "db"}}, body.ProductSectionList[0].Property)
		}
		guestProperties = &body
		w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
		w.WriteHeader(http.StatusAccepted)
	})

	updated, err := client.UpdateVAppGuestProperties(t.Context(), types.ParamsUpdateVAppGuestProperties{ID: vappID, GuestProperties: types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{{Class: "vami", Instance: "default", Properties: map[string]string{"hostname": "vm-02", "role": "db"}}}}})
	assert.NoError(t, err)
	if assert.NotNil(t, updated) && assert.Len(t, updated.Sections, 1) {
		assert.Equal(t, "vm-02", updated.Sections[0].Properties["hostname"])
		assert.Equal(t, "db", updated.Sections[0].Properties["role"])
	}
	assert.Equal(t, 2, getCalls)
}

func TestVAppGuestPropertiesPayloadDeterministicOrder(t *testing.T) {
	first := toAPIResponseGuestProperties(types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{{Class: "class", Instance: "b", Properties: map[string]string{"gamma": "3", "beta": "4"}}, {Class: "class", Instance: "a", Properties: map[string]string{"zeta": "1", "alpha": "2"}}}})
	second := toAPIResponseGuestProperties(types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{{Class: "class", Instance: "a", Properties: map[string]string{"alpha": "2", "zeta": "1"}}, {Class: "class", Instance: "b", Properties: map[string]string{"beta": "4", "gamma": "3"}}}})

	assert.Equal(t, second, first)
	if assert.Len(t, first.ProductSectionList, 2) {
		assert.Equal(t, "a", first.ProductSectionList[0].Instance)
		assert.Equal(t, []itypes.APIResponseProductProperty{{Key: "alpha", Value: "2"}, {Key: "zeta", Value: "1"}}, first.ProductSectionList[0].Property)
		assert.Equal(t, "b", first.ProductSectionList[1].Instance)
		assert.Equal(t, []itypes.APIResponseProductProperty{{Key: "beta", Value: "4"}, {Key: "gamma", Value: "3"}}, first.ProductSectionList[1].Property)
	}
}

func TestVAppGuestPropertiesAllowEmptySection(t *testing.T) {
	err := (types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{{Class: "vami", Instance: "default", Properties: map[string]string{}}}}).Validate()
	assert.NoError(t, err)

	err = (types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{{Class: "vami", Instance: "default"}}}).Validate()
	assert.Error(t, err)
}

func TestUndeployVAppModes(t *testing.T) {
	tests := []struct {
		name           string
		mode           types.VAppUndeployMode
		expectAction   string
		expectPowerOff *bool
	}{
		{name: "default power off", mode: "", expectAction: "", expectPowerOff: new(bool)},
		{name: "explicit power off", mode: types.VAppUndeployModePowerOff, expectAction: "", expectPowerOff: new(bool)},
		{name: "suspend", mode: types.VAppUndeployModeSuspend, expectAction: "suspend"},
		{name: "force", mode: types.VAppUndeployModeForce, expectAction: "force", expectPowerOff: new(bool)},
	}

	*tests[0].expectPowerOff = true
	*tests[1].expectPowerOff = true
	*tests[3].expectPowerOff = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)
			ms.CleanResponse(endpoints.UndeployVApp())
			ms.SetResponseFunc(endpoints.UndeployVApp(), func(w http.ResponseWriter, r *http.Request) {
				var body itypes.APIRequestUndeployVApp
				assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))
				assert.NoError(t, xml.NewDecoder(r.Body).Decode(&body))
				assert.Equal(t, "http://www.vmware.com/vcloud/v1.5", body.XMLNS)
				assert.Equal(t, tt.expectAction, body.UndeployAction)
				assert.Equal(t, tt.expectPowerOff, body.UndeployPowerOff)
				w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
				w.WriteHeader(http.StatusAccepted)
			})

			err := client.UndeployVApp(t.Context(), types.ParamsUndeployVApp{ID: testVAppID, Mode: tt.mode})
			assert.NoError(t, err)
		})
	}

	client, _ := newClient(t)
	assert.Error(t, client.UndeployVApp(t.Context(), types.ParamsUndeployVApp{ID: testVAppID, Mode: "invalid"}))
}

const testVAppID = "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012"

func isolatedRuntimeNetwork(name string) types.ModelVAppRuntimeNetworkConfig {
	prefix := 24
	return types.ModelVAppRuntimeNetworkConfig{
		NetworkName: name,
		Configuration: types.ModelVAppRuntimeNetworkConfiguration{
			FenceMode:   "isolated",
			IsConnected: true,
			IPScopes: []types.ModelVAppRuntimeIPScope{{
				Gateway:      "192.168.1.1",
				Netmask:      "255.255.255.0",
				DNS1:         "8.8.8.8",
				DNSSuffix:    "example.org",
				PrefixLength: &prefix,
				IPRanges: []types.ModelVAppRuntimeIPRange{{
					StartAddress: "192.168.1.10",
					EndAddress:   "192.168.1.20",
				}},
			}},
		},
	}
}

func existingRuntimeSection() *itypes.APIResponseVAppNetworkConfigSection {
	return &itypes.APIResponseVAppNetworkConfigSection{
		XMLNS: "http://www.vmware.com/vcloud/v1.5",
		OVF:   "http://schemas.dmtf.org/ovf/envelope/1",
		Info:  "Runtime network config",
		NetworkConfig: []itypes.APIResponseVAppNetworkConfig{
			richRuntimeNetwork("existing-net", "192.168.1.1"),
			richRuntimeNetwork("sibling-net", "10.0.0.1"),
		},
	}
}

func postMutationRuntimeSection(testName string) *itypes.APIResponseVAppNetworkConfigSection {
	switch testName {
	case "create isolated":
		section := existingRuntimeSection()
		section.NetworkConfig = append(section.NetworkConfig, toRawRuntimeNetwork(isolatedRuntimeNetwork("isolated-net")))
		return section
	case "add org network":
		network := isolatedRuntimeNetwork("org-net")
		network.Configuration.FenceMode = "bridged"
		network.Configuration.ParentNetwork = &types.VAppReference{HREF: "https://example.com/api/network/1234", ID: "urn:vcloud:network:1234", Name: "org-net"}
		network.Configuration.IPScopes = nil
		section := existingRuntimeSection()
		section.NetworkConfig = append(section.NetworkConfig, toRawRuntimeNetwork(network))
		return section
	case "update existing":
		section := existingRuntimeSection()
		section.NetworkConfig[0].Configuration.IPScopes.IPScope[0].Gateway = "192.168.2.1"
		return section
	case "remove existing":
		section := existingRuntimeSection()
		section.NetworkConfig = section.NetworkConfig[1:]
		return section
	default:
		return existingRuntimeSection()
	}
}

func toRawRuntimeNetwork(model types.ModelVAppRuntimeNetworkConfig) itypes.APIResponseVAppNetworkConfig {
	section := itypes.ToAPIResponseVAppNetworkConfigSection(&types.ModelVAppRuntimeNetworkConfigSection{
		Info:          "Runtime network config",
		NetworkConfig: []types.ModelVAppRuntimeNetworkConfig{model},
	})
	return section.NetworkConfig[0]
}

func richRuntimeNetwork(name, gateway string) itypes.APIResponseVAppNetworkConfig {
	connected := true
	return itypes.APIResponseVAppNetworkConfig{
		NetworkID:   "urn:vcloud:vappNetwork:" + name,
		NetworkHREF: "https://example.com/api/vAppNetwork/" + name,
		Description: name + " description",
		NetworkType: "application/vnd.vmware.vcloud.vAppNetwork+xml",
		NetworkName: name,
		Configuration: itypes.APIResponseVAppRuntimeConfiguration{
			FenceMode:   "isolated",
			IsConnected: &connected,
			Features: &itypes.APIResponseVAppRuntimeFeatures{
				FirewallService: &itypes.APIResponseVAppRuntimeInnerXML{InnerXML: "<IsEnabled>true</IsEnabled>"},
			},
			RouterInfo: &itypes.APIResponseVAppRuntimeRouterInfo{ExternalIP: "203.0.113.10", ConnectionType: "none"},
			IPScopes: &itypes.APIResponseVAppRuntimeIPScopes{IPScope: []itypes.APIResponseVAppRuntimeIPScope{{
				Gateway:  gateway,
				Netmask:  "255.255.255.0",
				IPRanges: itypes.APIResponseVAppRuntimeIPRanges{IPRange: []itypes.APIResponseVAppRuntimeIPRange{{StartAddress: gateway[:len(gateway)-1] + "10", EndAddress: gateway[:len(gateway)-1] + "20"}}},
			}}},
		},
	}
}

func TestVAppLifecycleEndpoints(t *testing.T) {
	tests := []struct {
		name string
		path string
		doc  string
		get  func() *cav.Endpoint
	}{
		{"deploy", "/api/vapp/{vapp-id}/action/deploy", "POST-DeployVApp.html", endpoints.DeployVApp},
		{"power off", "/api/vapp/{vapp-id}/power/action/powerOff", "POST-PowerOffVApp.html", endpoints.PowerOffVApp},
		{"power on", "/api/vapp/{vapp-id}/power/action/powerOn", "POST-PowerOnVApp.html", endpoints.PowerOnVApp},
		{"reboot", "/api/vapp/{vapp-id}/power/action/reboot", "POST-RebootVApp.html", endpoints.RebootVApp},
		{"reset", "/api/vapp/{vapp-id}/power/action/reset", "POST-ResetVApp.html", endpoints.ResetVApp},
		{"shutdown", "/api/vapp/{vapp-id}/power/action/shutdown", "POST-ShutdownVApp.html", endpoints.ShutdownVApp},
		{"suspend", "/api/vapp/{vapp-id}/power/action/suspend", "POST-SuspendVApp.html", endpoints.SuspendVApp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := tt.get()
			assert.Equal(t, tt.path, ep.PathTemplate)
			assert.Equal(t, "POST", string(ep.Method))
			assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/"+tt.doc, ep.DocumentationURL)
		})
	}
}

func TestDiscardSuspendedStateVAppEndpoint(t *testing.T) {
	ep := endpoints.DiscardSuspendedStateVApp()
	assert.Equal(t, "/api/vapp/{vapp-id}/action/discardSuspendedState", ep.PathTemplate)
	assert.Equal(t, "POST", string(ep.Method))
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/POST-DiscardVAppState.html", ep.DocumentationURL)
}

func TestUpdateVAppEndpointDocumentation(t *testing.T) {
	ep := endpoints.UpdateVApp()
	assert.Equal(t, "/api/vapp/{vapp-id}", ep.PathTemplate)
	assert.Equal(t, "PUT", string(ep.Method))
	assert.Equal(t, "https://developer.broadcom.com/xapis/vmware-cloud-director-api/39.1/doc/operations/PUT-Vm.html", ep.DocumentationURL)
}

func TestVAppLifecycleActionsRequireID(t *testing.T) {
	client, _ := newClient(t)
	actions := []struct {
		name string
		call func() error
	}{
		{"deploy", func() error { return client.DeployVApp(t.Context(), "") }},
		{"power off", func() error { return client.PowerOffVApp(t.Context(), "") }},
		{"power on", func() error { return client.PowerOnVApp(t.Context(), "") }},
		{"reboot", func() error { return client.RebootVApp(t.Context(), "") }},
		{"reset", func() error { return client.ResetVApp(t.Context(), "") }},
		{"shutdown", func() error { return client.ShutdownVApp(t.Context(), "") }},
		{"suspend", func() error { return client.SuspendVApp(t.Context(), "") }},
		{"discard suspended state", func() error { return client.DiscardSuspendedStateVApp(t.Context(), "") }},
	}
	for _, action := range actions {
		t.Run(action.name, func(t *testing.T) {
			assert.Error(t, action.call())
		})
	}
}

func TestComposeVAppRequiresSource(t *testing.T) {
	client, _ := newClient(t)
	_, err := client.ComposeVApp(t.Context(), types.ParamsComposeVApp{
		Name:  "test-vapp",
		VDCID: "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012",
	})
	assert.Error(t, err)
}

func TestInstantiateVAppTemplateRequiresTemplate(t *testing.T) {
	client, _ := newClient(t)
	_, err := client.InstantiateVAppTemplate(t.Context(), types.ParamsInstantiateVAppTemplate{
		Name:  "test-vapp",
		VDCID: "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012",
	})
	assert.Error(t, err)
}

func TestUpdateVAppPreservesExistingDescription(t *testing.T) {
	client, ms := newClient(t)
	vappID := "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012"
	deploymentLease := 3600
	storageLease := 7200
	getCalls := 0

	ms.CleanResponse(endpoints.GetVApp())
	ms.SetResponseFunc(endpoints.GetVApp(), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := &itypes.APIResponseGetVApp{
			ID:          vappID,
			Name:        "test-vapp-1",
			Description: "current description",
			LeaseSettings: itypes.APIResponseLeaseSettings{
				DeploymentLeaseInSeconds: deploymentLease,
				StorageLeaseInSeconds:    storageLease,
			},
		}
		if getCalls > 0 {
			resp.LeaseSettings.DeploymentLeaseInSeconds = 1800
		}
		getCalls++
		assert.NoError(t, json.NewEncoder(w).Encode(resp))
	})

	ms.CleanResponse(endpoints.UpdateVApp())
	ms.SetResponseFunc(endpoints.UpdateVApp(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.APIRequestUpdateVApp
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		if assert.NotNil(t, body.Description) {
			assert.Equal(t, "current description", *body.Description)
		}
		if assert.NotNil(t, body.LeaseSettings) && assert.NotNil(t, body.LeaseSettings.DeploymentLeaseInSeconds) {
			assert.Equal(t, 1800, *body.LeaseSettings.DeploymentLeaseInSeconds)
		}
		if assert.NotNil(t, body.LeaseSettings) && assert.NotNil(t, body.LeaseSettings.StorageLeaseInSeconds) {
			assert.Equal(t, storageLease, *body.LeaseSettings.StorageLeaseInSeconds)
		}
		w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.UpdateVApp(t.Context(), types.ParamsUpdateVApp{ID: vappID, DeploymentLeaseInSeconds: new(1800)})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "current description", resp.Description)
	assert.NotNil(t, resp.DeploymentLeaseInSeconds)
	assert.Equal(t, 1800, *resp.DeploymentLeaseInSeconds)

	ms.CleanResponse(endpoints.GetVApp())
	ms.CleanResponse(endpoints.UpdateVApp())
}

func TestUpdateVAppCanClearDescription(t *testing.T) {
	client, ms := newClient(t)
	vappID := "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012"
	empty := ""
	getCalls := 0

	ms.CleanResponse(endpoints.GetVApp())
	ms.SetResponseFunc(endpoints.GetVApp(), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := &itypes.APIResponseGetVApp{
			ID:          vappID,
			Name:        "test-vapp-1",
			Description: "current description",
		}
		if getCalls > 0 {
			resp.Description = empty
		}
		getCalls++
		assert.NoError(t, json.NewEncoder(w).Encode(resp))
	})

	ms.CleanResponse(endpoints.UpdateVApp())
	ms.SetResponseFunc(endpoints.UpdateVApp(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.APIRequestUpdateVApp
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		if assert.NotNil(t, body.Description) {
			assert.Equal(t, empty, *body.Description)
		}
		w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.UpdateVApp(t.Context(), types.ParamsUpdateVApp{ID: vappID, Description: &empty})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, empty, resp.Description)

	ms.CleanResponse(endpoints.GetVApp())
	ms.CleanResponse(endpoints.UpdateVApp())
}

func TestDeleteVApp(t *testing.T) {
	t.Run("undeploy ok", func(t *testing.T) {
		client, _ := newClient(t)
		assert.NoError(t, client.DeleteVApp(t.Context(), types.ParamsDeleteVApp{ID: testVAppID}))
	})

	t.Run("undeploy benign error continues and delete happens", func(t *testing.T) {
		client, ms := newClient(t)
		deleteCalls := 0

		ms.CleanResponse(endpoints.UndeployVApp())
		ms.SetResponseFunc(endpoints.UndeployVApp(), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"message": "The requested operation could not be executed since vApp test-vapp is not running.", "majorErrorCode": http.StatusBadRequest, "minorErrorCode": "BAD_REQUEST"}))
		})

		ms.CleanResponse(endpoints.DeleteVApp())
		ms.SetResponseFunc(endpoints.DeleteVApp(), func(w http.ResponseWriter, r *http.Request) {
			deleteCalls++
			w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
			w.WriteHeader(http.StatusAccepted)
		})

		err := client.DeleteVApp(t.Context(), types.ParamsDeleteVApp{ID: testVAppID})
		assert.NoError(t, err)
		assert.Equal(t, 1, deleteCalls)
	})

	t.Run("undeploy non-benign execute error stops and delete is not called", func(t *testing.T) {
		client, ms := newClient(t)
		deleteCalls := 0

		ms.CleanResponse(endpoints.UndeployVApp())
		ms.SetResponseFunc(endpoints.UndeployVApp(), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"message": "another undeploy failure", "majorErrorCode": http.StatusBadRequest, "minorErrorCode": "BAD_REQUEST"}))
		})

		ms.CleanResponse(endpoints.DeleteVApp())
		ms.SetResponseFunc(endpoints.DeleteVApp(), func(w http.ResponseWriter, r *http.Request) {
			deleteCalls++
			w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
			w.WriteHeader(http.StatusAccepted)
		})

		err := client.DeleteVApp(t.Context(), types.ParamsDeleteVApp{ID: testVAppID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "undeploy")
		assert.Equal(t, 0, deleteCalls)
	})

	t.Run("undeploy non-benign await error stops and delete is not called", func(t *testing.T) {
		client, ms := newClient(t)
		deleteCalls := 0
		jobCalls := 0

		ms.CleanResponse(endpoints.UndeployVApp())
		ms.SetResponseFunc(endpoints.UndeployVApp(), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
			w.WriteHeader(http.StatusAccepted)
		})

		ms.CleanResponse(cav.MustGetEndpoint("GetJobVmware"))
		ms.SetResponseFunc(cav.MustGetEndpoint("GetJobVmware"), func(w http.ResponseWriter, r *http.Request) {
			jobCalls++
			w.Header().Set("Content-Type", "application/json")
			assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"operation": "undeploy",
				"href":      "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d",
				"id":        "87ab1934-0146-4fb0-80bc-815fea03214d",
				"status":    "error",
				"error": map[string]any{
					"message":        "task failed for another reason",
					"majorErrorCode": http.StatusBadRequest,
					"minorErrorCode": "BAD_REQUEST",
				},
			}))
		})

		ms.CleanResponse(endpoints.DeleteVApp())
		ms.SetResponseFunc(endpoints.DeleteVApp(), func(w http.ResponseWriter, r *http.Request) {
			deleteCalls++
			w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
			w.WriteHeader(http.StatusAccepted)
		})

		err := client.DeleteVApp(t.Context(), types.ParamsDeleteVApp{ID: testVAppID})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "undeploy")
		assert.Equal(t, 0, deleteCalls)
		assert.GreaterOrEqual(t, jobCalls, 1)
	})

	t.Run("undeploy benign await error continues and delete happens", func(t *testing.T) {
		client, ms := newClient(t)
		deleteCalls := 0
		jobCalls := 0

		ms.CleanResponse(endpoints.UndeployVApp())
		ms.SetResponseFunc(endpoints.UndeployVApp(), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
			w.WriteHeader(http.StatusAccepted)
		})

		ms.CleanResponse(cav.MustGetEndpoint("GetJobVmware"))
		ms.SetResponseFunc(cav.MustGetEndpoint("GetJobVmware"), func(w http.ResponseWriter, r *http.Request) {
			jobCalls++
			w.Header().Set("Content-Type", "application/json")
			status := "success"
			payload := map[string]any{
				"operation": "undeploy",
				"href":      "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d",
				"id":        "87ab1934-0146-4fb0-80bc-815fea03214d",
				"status":    status,
			}
			if jobCalls == 2 {
				payload["status"] = "error"
				payload["error"] = map[string]any{
					"message":        "The requested operation could not be executed since vApp test-vapp is not running.",
					"majorErrorCode": http.StatusBadRequest,
					"minorErrorCode": "BAD_REQUEST",
				}
			}
			assert.NoError(t, json.NewEncoder(w).Encode(payload))
		})

		ms.CleanResponse(endpoints.DeleteVApp())
		ms.SetResponseFunc(endpoints.DeleteVApp(), func(w http.ResponseWriter, r *http.Request) {
			deleteCalls++
			w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
			w.WriteHeader(http.StatusAccepted)
		})

		err := client.DeleteVApp(t.Context(), types.ParamsDeleteVApp{ID: testVAppID})
		assert.NoError(t, err)
		assert.Equal(t, 1, deleteCalls)
	})
}

func Test_NewClient_ClientNil(t *testing.T) {
	c, err := New(nil)
	assert.Nil(t, c, "Expected nil client when input is nil")
	assert.Error(t, err, "Expected error when input is nil")
}
