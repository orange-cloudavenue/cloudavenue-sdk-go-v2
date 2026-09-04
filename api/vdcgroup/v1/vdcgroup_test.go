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

func TestListVDCGroup(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsListVDCGroup
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name:        "List Vdc Groups no Param",
			expectedErr: false,
		},
		{
			name: "List Vdc Groups by ID",
			params: types.ParamsListVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			expectedErr: false,
		},
		{
			name: "List Vdc Groups by Name",
			params: types.ParamsListVDCGroup{
				Name: generator.MustGenerate("{word}"),
			},
			expectedErr: false,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsListVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), tt.mockResponse, &tt.mockResponseStatus)
			}

			resp, err := client.ListVDCGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.VDCGroups, "Vdc Groups should not be empty")
			for _, vdcGroup := range resp.VDCGroups {
				assert.NotEmpty(t, vdcGroup.ID, "Vdc Group ID should not be empty")
				assert.NotEmpty(t, vdcGroup.Name, "Vdc Group Name should not be empty")
			}
		})
	}
}

func TestGetVDCGroup(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsGetVDCGroup
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Get Vdc Group by ID",
			params: types.ParamsGetVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			expectedErr: false,
		},
		{
			name: "Simulate VDCGroup not found",
			params: types.ParamsGetVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockResponse:       &itypes.APIResponseListVDCGroup{},
			mockResponseStatus: 200,
			expectedErr:        true,
		},
		{
			name: "Simulate multi VDCGroup found",
			params: types.ParamsGetVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockResponse: func() *itypes.APIResponseListVDCGroup {
				resp := &itypes.APIResponseListVDCGroup{}
				generator.MustStruct(resp)
				return resp
			}(),
			mockResponseStatus: 200,
			expectedErr:        true,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsGetVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.name == "Get Vdc Group by ID" {
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
					Values: []itypes.APIResponseListVDCGroupDetails{{
						ID:   tt.params.ID,
						Name: generator.MustGenerate("{word}"),
					}},
				}, nil)
			} else if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), tt.mockResponse, &tt.mockResponseStatus)
			}

			resp, err := client.GetVDCGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.ID, "Vdc Group ID should not be empty")
			assert.NotEmpty(t, resp.Name, "Vdc Group Name should not be empty")
		})
	}
}

func TestCreateVDCGroup(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsCreateVDCGroup
		mockResponseStatus int
		mockResponse       any

		mockListVDCGroupResponse       any
		mockListVDCGroupResponseStatus int

		mockListVDCResponse       any
		mockListVDCResponseStatus int

		expectedErr bool
	}{
		{
			name: "Create Vdc Group",
			params: types.ParamsCreateVDCGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				VDCs: []types.ParamsCreateVDCGroupVDC{
					{
						ID:   generator.MustGenerate("{urn:vdc}"),
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCGroupResponse: &itypes.APIResponseListVDCGroup{
				Values: []itypes.APIResponseListVDCGroupDetails{},
			},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Create Vdc Group without VDC ID",
			params: types.ParamsCreateVDCGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				VDCs: []types.ParamsCreateVDCGroupVDC{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCGroupResponse: &itypes.APIResponseListVDCGroup{
				Values: []itypes.APIResponseListVDCGroupDetails{},
			},
			mockListVDCGroupResponseStatus: 200,
			mockListVDCResponse: &itypes.APIResponseListVDC{Records: []itypes.APIResponseListVDCRecord{{
				HREF: "https://example.invalid/api/vdc/23d9d9a0-0efb-4591-9869-3dd33aa2247c",
				ID:   "urn:vdc:23d9d9a0-0efb-4591-9869-3dd33aa2247c",
				Name: "resolve-me",
			}}},
			mockListVDCResponseStatus: 200,
			expectedErr:               false,
		},
		{
			name: "Error List VDCGroup",
			params: types.ParamsCreateVDCGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				VDCs: []types.ParamsCreateVDCGroupVDC{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCGroupResponseStatus: 404,
			expectedErr:                    true,
		},
		{
			name: "Error List VDC",
			params: types.ParamsCreateVDCGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				VDCs: []types.ParamsCreateVDCGroupVDC{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCGroupResponse: &itypes.APIResponseListVDCGroup{
				Values: []itypes.APIResponseListVDCGroupDetails{},
			},
			mockListVDCGroupResponseStatus: 200,
			mockListVDCResponseStatus:      404,
			expectedErr:                    true,
		},
		{
			name: "VDCGroup already exists",
			params: types.ParamsCreateVDCGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				VDCs: []types.ParamsCreateVDCGroupVDC{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCGroupResponse: &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
				Name: "existing-vdc-group",
			}}},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    true,
		},
		{
			name: "Error 400 Bad Request",
			params: types.ParamsCreateVDCGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				VDCs: []types.ParamsCreateVDCGroupVDC{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCGroupResponse: &itypes.APIResponseListVDCGroup{
				Values: []itypes.APIResponseListVDCGroupDetails{},
			},
			mockListVDCGroupResponseStatus: 200,
			mockResponseStatus:             400,
			expectedErr:                    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.CreateVDCGroup())
				ms.SetResponse(endpoints.CreateVDCGroup(), tt.mockResponse, &tt.mockResponseStatus)
			} else {
				ms.CleanResponse(endpoints.CreateVDCGroup())
			}

			if tt.mockListVDCGroupResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), tt.mockListVDCGroupResponse, &tt.mockListVDCGroupResponseStatus)
			}

			if tt.mockListVDCResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVDC())
				ms.SetResponse(endpoints.ListVDC(), tt.mockListVDCResponse, &tt.mockListVDCResponseStatus)
				ms.CleanResponse(endpoints.QueryEdgeGateway())
				ms.SetResponse(endpoints.QueryEdgeGateway(), tt.mockListVDCResponse, &tt.mockListVDCResponseStatus)
			}

			if tt.name == "Add VDC to VDC Group with VDCGroup and VDC by name" {
				tt.params.Name = "current-vdc-group"
				tt.params.VDCs[0].Name = "resolve-me"
			}

			if tt.name == "Add VDC to VDC Group with VDCGroup and VDC by name" {
				tt.params.Name = "current-vdc-group"
				tt.params.VDCs[0].Name = "resolve-me"
			}

			if tt.name == "Create Vdc Group without VDC ID" {
				tt.params.VDCs[0].Name = "resolve-me"
			}

			if tt.name == "VDCGroup already exists" {
				tt.params.Name = "existing-vdc-group"
			}

			resp, err := client.CreateVDCGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.Equal(t, tt.params.Name, resp.Name)
		})
	}
}

func TestDeleteVDCGroup(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsDeleteVDCGroup

		mockResponse       any
		mockResponseStatus int

		mockListVDCGroupResponse       any
		mockListVDCGroupResponseStatus int

		expectedErr bool
	}{
		{
			name: "Delete Vdc Group",
			params: types.ParamsDeleteVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Delete with VDC Group Name",
			params: types.ParamsDeleteVDCGroup{
				Name: generator.MustGenerate("{word}"),
			},
			mockListVDCGroupResponse: &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
				ID:   generator.MustGenerate("{urn:vdcGroup}"),
				Name: "delete-by-name",
			}}},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Delete Vdc Group with Force",
			params: types.ParamsDeleteVDCGroup{
				ID:    generator.MustGenerate("{urn:vdcGroup}"),
				Force: true,
			},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Failed to retrieve Vdc Group",
			params: types.ParamsDeleteVDCGroup{
				Name: generator.MustGenerate("{word}"),
			},
			mockListVDCGroupResponseStatus: 404,
			expectedErr:                    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.DeleteVDCGroup())
				ms.SetResponse(endpoints.DeleteVDCGroup(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListVDCGroupResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), tt.mockListVDCGroupResponse, &tt.mockListVDCGroupResponseStatus)
			}

			if tt.name == "Delete with VDC Group Name" {
				tt.params.Name = "delete-by-name"
			}

			err := client.DeleteVDCGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestAddVdcToVDCGroup(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsAddVDCToVDCGroup

		mockResponse       any
		mockResponseStatus int

		mockListVDCGroupResponse       any
		mockListVDCGroupResponseStatus int

		mockListVDCResponse       any
		mockListVDCResponseStatus int

		expectedErr bool
	}{
		{
			name: "Add VDC to VDC Group",
			params: types.ParamsAddVDCToVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Add VDC to VDC Group with VDCGroup and VDC by name",
			params: types.ParamsAddVDCToVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCGroupResponseStatus: 200,
			mockListVDCResponse: &itypes.APIResponseListVDC{Records: []itypes.APIResponseListVDCRecord{{
				HREF: "https://example.invalid/api/vdc/495fe7f3-3051-4408-b483-c71ce7bc2d26",
				ID:   "urn:vdc:495fe7f3-3051-4408-b483-c71ce7bc2d26",
				Name: "resolve-me",
			}}},
			mockListVDCResponseStatus: 200,
			expectedErr:               false,
		},

		{
			name: "Failed to retrieve VDC Group",
			params: types.ParamsAddVDCToVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVDCGroupResponseStatus: 404,
			expectedErr:                    true,
		},
		{
			name: "Failed to list VDC",
			params: types.ParamsAddVDCToVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCResponseStatus: 404,
			expectedErr:               true,
		},
		{
			name: "Failed VDCGroup doesn't exist",
			params: types.ParamsAddVDCToVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVDCGroupResponse: itypes.APIResponseListVDCGroup{
				Values: []itypes.APIResponseListVDCGroupDetails{},
			},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    true,
		},
		{
			name: "Failed to add vdc to vdc group",
			params: types.ParamsAddVDCToVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockResponseStatus: 400,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.name == "Add VDC to VDC Group" {
				tt.mockListVDCGroupResponseStatus = 200
				tt.mockListVDCGroupResponse = &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
					ID:          tt.params.ID,
					Name:        "current-vdc-group",
					Description: "current-description",
					Vdcs: []itypes.APIResponseVDCGroupParticipatingVDC{{
						VDC: itypes.APIResponseVDCGroupParticipatingVDCRef{
							ID:   generator.MustGenerate("{urn:vdc}"),
							Name: "existing-vdc",
						},
					}},
				}}}
			}

			if tt.name == "Add VDC to VDC Group with VDCGroup and VDC by name" {
				tt.params.Name = "current-vdc-group"
				tt.params.Vdcs[0].Name = "resolve-me"
				tt.mockListVDCGroupResponseStatus = 200
				tt.mockListVDCGroupResponse = &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
					ID:          generator.MustGenerate("{urn:vdcGroup}"),
					Name:        "current-vdc-group",
					Description: "current-description",
				}}}
				tt.mockListVDCResponseStatus = 200
				tt.mockListVDCResponse = &itypes.APIResponseListVDC{Records: []itypes.APIResponseListVDCRecord{{
					HREF: "https://example.invalid/api/vdc/495fe7f3-3051-4408-b483-c71ce7bc2d26",
					ID:   "urn:vdc:495fe7f3-3051-4408-b483-c71ce7bc2d26",
					Name: "resolve-me",
				}}}
			}

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.UpdateVDCGroup())
				ms.SetResponse(endpoints.UpdateVDCGroup(), tt.mockResponse, &tt.mockResponseStatus)
			} else {
				ms.CleanResponse(endpoints.UpdateVDCGroup())
			}

			if tt.mockListVDCGroupResponseStatus != 0 {
				if tt.mockListVDCGroupResponse == nil && !tt.expectedErr {
					vdcGroupName := "current-vdc-group"
					if tt.params.Name != "" {
						vdcGroupName = tt.params.Name
					}
					vdcGroupVdcs := []itypes.APIResponseVDCGroupParticipatingVDC{{
						VDC: itypes.APIResponseVDCGroupParticipatingVDCRef{
							ID:   generator.MustGenerate("{urn:vdc}"),
							Name: "existing-vdc",
						},
					}}
					if tt.name == "Add VDC to VDC Group with VDCGroup and VDC by name" {
						vdcGroupName = "current-vdc-group"
						vdcGroupVdcs = nil
					}
					tt.mockListVDCGroupResponse = &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
						ID: func() string {
							if tt.params.ID != "" {
								return tt.params.ID
							}
							return generator.MustGenerate("{urn:vdcGroup}")
						}(),
						Name:        vdcGroupName,
						Description: "current-description",
						Vdcs:        vdcGroupVdcs,
					}}}
				}
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), tt.mockListVDCGroupResponse, &tt.mockListVDCGroupResponseStatus)
			}

			if tt.mockListVDCResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVDC())
				ms.SetResponse(endpoints.ListVDC(), tt.mockListVDCResponse, &tt.mockListVDCResponseStatus)
				ms.CleanResponse(endpoints.QueryEdgeGateway())
				ms.SetResponse(endpoints.QueryEdgeGateway(), tt.mockListVDCResponse, &tt.mockListVDCResponseStatus)
			}

			err := client.AddVDCToVDCGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestRemoveVdcToVDCGroup(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsRemoveVDCFromVDCGroup

		mockResponse       any
		mockResponseStatus int

		mockListVDCGroupResponse       any
		mockListVDCGroupResponseStatus int

		expectedErr bool
	}{
		{
			name: "Remove VDC from VDC Group",
			params: types.ParamsRemoveVDCFromVDCGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "Remove VDC from VDC Group with VDCGroup and VDC by name",
			params: types.ParamsRemoveVDCFromVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						Name: "my-vdc",
					},
				},
			},
			expectedErr: false,
			mockListVDCGroupResponse: itypes.APIResponseListVDCGroup{
				Values: []itypes.APIResponseListVDCGroupDetails{
					{
						ID:   generator.MustGenerate("{urn:vdcGroup}"),
						Name: "my-vdc-group",
						Vdcs: []itypes.APIResponseVDCGroupParticipatingVDC{
							{
								VDC: itypes.APIResponseVDCGroupParticipatingVDCRef{
									ID:   generator.MustGenerate("{urn:vdc}"),
									Name: "my-vdc",
								},
							},
						},
					},
				},
			},
			mockListVDCGroupResponseStatus: 200,
		},

		{
			name: "Failed to retrieve VDC Group",
			params: types.ParamsRemoveVDCFromVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVDCGroupResponseStatus: 404,
			expectedErr:                    true,
		},
		{
			name: "Failed VDCGroup doesn't exist",
			params: types.ParamsRemoveVDCFromVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVDCGroupResponse: itypes.APIResponseListVDCGroup{
				Values: []itypes.APIResponseListVDCGroupDetails{},
			},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    true,
		},
		{
			name: "Failed to remove vdc from vdc group",
			params: types.ParamsRemoveVDCFromVDCGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVDCGroupVDC{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockResponseStatus: 400,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.name == "Remove VDC from VDC Group" {
				tt.mockListVDCGroupResponseStatus = 200
				tt.mockListVDCGroupResponse = &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
					ID:          tt.params.ID,
					Name:        "current-vdc-group",
					Description: "current-description",
					Vdcs: []itypes.APIResponseVDCGroupParticipatingVDC{
						{
							VDC: itypes.APIResponseVDCGroupParticipatingVDCRef{
								ID:   tt.params.Vdcs[0].ID,
								Name: "remove-me",
							},
						},
						{
							VDC: itypes.APIResponseVDCGroupParticipatingVDCRef{
								ID:   generator.MustGenerate("{urn:vdc}"),
								Name: "keep-vdc",
							},
						},
					},
				}}}
			}

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.UpdateVDCGroup())
				ms.SetResponse(endpoints.UpdateVDCGroup(), tt.mockResponse, &tt.mockResponseStatus)
			} else {
				ms.CleanResponse(endpoints.UpdateVDCGroup())
			}

			if tt.mockListVDCGroupResponseStatus != 0 {
				if tt.mockListVDCGroupResponse == nil && !tt.expectedErr {
					tt.mockListVDCGroupResponse = &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
						ID: func() string {
							if tt.params.ID != "" {
								return tt.params.ID
							}
							return generator.MustGenerate("{urn:vdcGroup}")
						}(),
						Name: func() string {
							if tt.params.Name != "" {
								return tt.params.Name
							}
							return "current-vdc-group"
						}(),
						Description: "current-description",
						Vdcs: []itypes.APIResponseVDCGroupParticipatingVDC{
							{
								VDC: itypes.APIResponseVDCGroupParticipatingVDCRef{
									ID:   generator.MustGenerate("{urn:vdc}"),
									Name: "my-vdc",
								},
							},
							{
								VDC: itypes.APIResponseVDCGroupParticipatingVDCRef{
									ID:   generator.MustGenerate("{urn:vdc}"),
									Name: "keep-vdc",
								},
							},
						},
					}}}
				}
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), tt.mockListVDCGroupResponse, &tt.mockListVDCGroupResponseStatus)
			}

			if tt.name == "Remove VDC from VDC Group with VDCGroup and VDC by name" {
				tt.params.Name = "my-vdc-group"
			}

			err := client.RemoveVDCFromVDCGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestUpdateVDCGroup(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsUpdateVDCGroup

		mockResponse       any
		mockResponseStatus int

		mockListVDCGroupResponse       any
		mockListVDCGroupResponseStatus int

		expectedErr bool
	}{
		{
			name: "Update VDC Group",
			params: types.ParamsUpdateVDCGroup{
				ID:          generator.MustGenerate("{urn:vdcGroup}"),
				Description: new("My updated VDC Group"),
			},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Update VDC Group with VDCGroup name",
			params: types.ParamsUpdateVDCGroup{
				Name:        "my-updated-vdc-group",
				Description: new("My updated VDC Group"),
			},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Failed to update VDC Group",
			params: types.ParamsUpdateVDCGroup{
				ID:          generator.MustGenerate("{urn:vdcGroup}"),
				Name:        "my-updated-vdc-group",
				Description: new("My updated VDC Group"),
			},
			mockResponseStatus: 400,
			expectedErr:        true,
		},
		{
			name: "Failed to list VDC Group",
			params: types.ParamsUpdateVDCGroup{
				ID:          generator.MustGenerate("{urn:vdcGroup}"),
				Name:        "my-updated-vdc-group",
				Description: new("My updated VDC Group"),
			},
			mockListVDCGroupResponseStatus: 404,
			expectedErr:                    true,
		},
		{
			name: "List VDCGroup are empty",
			params: types.ParamsUpdateVDCGroup{
				ID:          generator.MustGenerate("{urn:vdcGroup}"),
				Name:        "my-updated-vdc-group",
				Description: new("My updated VDC Group"),
			},
			mockListVDCGroupResponse:       itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{}},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.UpdateVDCGroup())
				ms.SetResponse(endpoints.UpdateVDCGroup(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListVDCGroupResponseStatus != 0 {
				if tt.mockListVDCGroupResponse == nil && !tt.expectedErr {
					tt.mockListVDCGroupResponse = &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
						ID: func() string {
							if tt.params.ID != "" {
								return tt.params.ID
							}
							return generator.MustGenerate("{urn:vdcGroup}")
						}(),
						Name: func() string {
							if tt.params.Name != "" {
								return tt.params.Name
							}
							return "current-vdc-group"
						}(),
						Description: "current-description",
					}}}
				}
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), tt.mockListVDCGroupResponse, &tt.mockListVDCGroupResponseStatus)
			}

			vdc, err := client.UpdateVDCGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, vdc)
			if tt.params.Description != nil {
				assert.Equal(t, *tt.params.Description, vdc.Description)
			}
		})
	}
}
