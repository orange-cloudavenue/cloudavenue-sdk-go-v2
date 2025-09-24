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

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/utils"
)

func TestListVdcGroup(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsListVdcGroup
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
			params: types.ParamsListVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			expectedErr: false,
		},
		{
			name: "List Vdc Groups by Name",
			params: types.ParamsListVdcGroup{
				Name: generator.MustGenerate("{word}"),
			},
			expectedErr: false,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsListVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.ListVdcGroup().CleanMockResponse()
				endpoints.ListVdcGroup().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			resp, err := client.ListVdcGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.VdcGroups, "Vdc Groups should not be empty")
			for _, vdcGroup := range resp.VdcGroups {
				assert.NotEmpty(t, vdcGroup.ID, "Vdc Group ID should not be empty")
				assert.NotEmpty(t, vdcGroup.Name, "Vdc Group Name should not be empty")
			}
		})
	}
}

func TestGetVdcGroup(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsGetVdcGroup
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Get Vdc Group by ID",
			params: types.ParamsGetVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			expectedErr: false,
		},
		{
			name: "Simulate VDCGroup not found",
			params: types.ParamsGetVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockResponse:       &itypes.ApiResponseListVdcGroup{},
			mockResponseStatus: 200,
			expectedErr:        true,
		},
		{
			name: "Simulate multi VDCGroup found",
			params: types.ParamsGetVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockResponse: func() *itypes.ApiResponseListVdcGroup {
				resp := &itypes.ApiResponseListVdcGroup{}
				generator.MustStruct(resp)
				return resp
			}(),
			mockResponseStatus: 200,
			expectedErr:        true,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsGetVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.ListVdcGroup().CleanMockResponse()
				endpoints.ListVdcGroup().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			resp, err := client.GetVdcGroup(t.Context(), tt.params)
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

func TestCreateVdcGroup(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsCreateVdcGroup
		mockResponseStatus int
		mockResponse       any

		mockListVdcGroupResponse       any
		mockListVdcGroupResponseStatus int

		mockListVdcResponse       any
		mockListVdcResponseStatus int

		expectedErr bool
	}{
		{
			name: "Create Vdc Group",
			params: types.ParamsCreateVdcGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						ID:   generator.MustGenerate("{urn:vdc}"),
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVdcGroupResponse: &itypes.ApiResponseListVdcGroup{
				Values: []itypes.ApiResponseListVdcGroupDetails{},
			},
			mockListVdcGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Create Vdc Group without VDC ID",
			params: types.ParamsCreateVdcGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVdcGroupResponse: &itypes.ApiResponseListVdcGroup{
				Values: []itypes.ApiResponseListVdcGroupDetails{},
			},
			mockListVdcGroupResponseStatus: 200,
			expectedErr:                    false,
		},
		{
			name: "Error List VDCGroup",
			params: types.ParamsCreateVdcGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVdcGroupResponseStatus: 404,
			expectedErr:                    true,
		},
		{
			name: "Error List VDC",
			params: types.ParamsCreateVdcGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVdcGroupResponse: &itypes.ApiResponseListVdcGroup{
				Values: []itypes.ApiResponseListVdcGroupDetails{},
			},
			mockListVdcGroupResponseStatus: 200,
			mockListVdcResponseStatus:      404,
			expectedErr:                    true,
		},
		{
			name: "VDCGroup already exists",
			params: types.ParamsCreateVdcGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "Error 400 Bad Request",
			params: types.ParamsCreateVdcGroup{
				Name:        generator.MustGenerate("{word}"),
				Description: generator.MustGenerate("{sentence}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVdcGroupResponse: &itypes.ApiResponseListVdcGroup{
				Values: []itypes.ApiResponseListVdcGroupDetails{},
			},
			mockListVdcGroupResponseStatus: 200,
			mockResponseStatus:             400,
			expectedErr:                    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.CreateVdcGroup().CleanMockResponse()
				endpoints.CreateVdcGroup().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListVdcGroupResponseStatus != 0 {
				endpoints.ListVdcGroup().CleanMockResponse()
				endpoints.ListVdcGroup().SetMockResponse(tt.mockListVdcGroupResponse, &tt.mockListVdcGroupResponseStatus)
			}

			if tt.mockListVdcResponseStatus != 0 {
				endpoints.ListVdc().CleanMockResponse()
				endpoints.ListVdc().SetMockResponse(tt.mockListVdcResponse, &tt.mockListVdcResponseStatus)
			}

			client := newClient(t)

			resp, err := client.CreateVdcGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
		})
	}
}

func TestDeleteVdcGroup(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsDeleteVdcGroup

		mockResponse       any
		mockResponseStatus int

		mockListVdcGroupResponse       any
		mockListVdcGroupResponseStatus int

		expectedErr bool
	}{
		{
			name: "Delete Vdc Group",
			params: types.ParamsDeleteVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
			},
			expectedErr: false,
		},
		{
			name: "Delete with VDC Group Name",
			params: types.ParamsDeleteVdcGroup{
				Name: generator.MustGenerate("{word}"),
			},
			expectedErr: false,
		},
		{
			name: "Delete Vdc Group with Force",
			params: types.ParamsDeleteVdcGroup{
				ID:    generator.MustGenerate("{urn:vdcGroup}"),
				Force: true,
			},
			expectedErr: false,
		},
		{
			name: "Failed to retrieve Vdc Group",
			params: types.ParamsDeleteVdcGroup{
				Name: generator.MustGenerate("{word}"),
			},
			mockListVdcGroupResponseStatus: 404,
			expectedErr:                    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.DeleteVdcGroup().CleanMockResponse()
				endpoints.DeleteVdcGroup().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListVdcGroupResponseStatus != 0 {
				endpoints.ListVdcGroup().CleanMockResponse()
				endpoints.ListVdcGroup().SetMockResponse(tt.mockListVdcGroupResponse, &tt.mockListVdcGroupResponseStatus)
			}

			client := newClient(t)

			err := client.DeleteVdcGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestAddVdcToVdcGroup(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsAddVdcToVdcGroup

		mockResponse       any
		mockResponseStatus int

		mockListVdcGroupResponse       any
		mockListVdcGroupResponseStatus int

		mockListVdcResponse       any
		mockListVdcResponseStatus int

		expectedErr bool
	}{
		{
			name: "Add VDC to VDC Group",
			params: types.ParamsAddVdcToVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "Add VDC to VDC Group with VDCGroup and VDC by name",
			params: types.ParamsAddVdcToVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			expectedErr: false,
		},

		{
			name: "Failed to retrieve VDC Group",
			params: types.ParamsAddVdcToVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVdcGroupResponseStatus: 404,
			expectedErr:                    true,
		},
		{
			name: "Failed to list VDC",
			params: types.ParamsAddVdcToVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVdcResponseStatus: 404,
			expectedErr:               true,
		},
		{
			name: "Failed VdcGroup doesn't exist",
			params: types.ParamsAddVdcToVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVdcGroupResponse: itypes.ApiResponseListVdcGroup{
				Values: []itypes.ApiResponseListVdcGroupDetails{},
			},
			mockListVdcGroupResponseStatus: 200,
			expectedErr:                    true,
		},
		{
			name: "Failed to add vdc to vdc group",
			params: types.ParamsAddVdcToVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
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
			if tt.mockResponseStatus != 0 {
				endpoints.UpdateVdcGroup().CleanMockResponse()
				endpoints.UpdateVdcGroup().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListVdcGroupResponseStatus != 0 {
				endpoints.ListVdcGroup().CleanMockResponse()
				endpoints.ListVdcGroup().SetMockResponse(tt.mockListVdcGroupResponse, &tt.mockListVdcGroupResponseStatus)
			}

			if tt.mockListVdcResponseStatus != 0 {
				endpoints.ListVdc().CleanMockResponse()
				endpoints.ListVdc().SetMockResponse(tt.mockListVdcResponse, &tt.mockListVdcResponseStatus)
			}

			client := newClient(t)

			err := client.AddVdcToVdcGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestRemoveVdcToVdcGroup(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsRemoveVdcFromVdcGroup

		mockResponse       any
		mockResponseStatus int

		mockListVdcGroupResponse       any
		mockListVdcGroupResponseStatus int

		expectedErr bool
	}{
		{
			name: "Remove VDC from VDC Group",
			params: types.ParamsRemoveVdcFromVdcGroup{
				ID: generator.MustGenerate("{urn:vdcGroup}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "Remove VDC from VDC Group with VDCGroup and VDC by name",
			params: types.ParamsRemoveVdcFromVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						Name: "my-vdc",
					},
				},
			},
			expectedErr: false,
			mockListVdcGroupResponse: itypes.ApiResponseListVdcGroup{
				Values: []itypes.ApiResponseListVdcGroupDetails{
					{
						ID:   generator.MustGenerate("{urn:vdcGroup}"),
						Name: "my-vdc-group",
						Vdcs: []itypes.ApiResponseVdcGroupParticipatingVdc{
							{
								Vdc: itypes.ApiResponseVdcGroupParticipatingVdcRef{
									ID:   generator.MustGenerate("{urn:vdc}"),
									Name: "my-vdc",
								},
							},
						},
					},
				},
			},
			mockListVdcGroupResponseStatus: 200,
		},

		{
			name: "Failed to retrieve VDC Group",
			params: types.ParamsRemoveVdcFromVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVdcGroupResponseStatus: 404,
			expectedErr:                    true,
		},
		{
			name: "Failed VdcGroup doesn't exist",
			params: types.ParamsRemoveVdcFromVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
					{
						ID: generator.MustGenerate("{urn:vdc}"),
					},
				},
			},
			mockListVdcGroupResponse: itypes.ApiResponseListVdcGroup{
				Values: []itypes.ApiResponseListVdcGroupDetails{},
			},
			mockListVdcGroupResponseStatus: 200,
			expectedErr:                    true,
		},
		{
			name: "Failed to remove vdc from vdc group",
			params: types.ParamsRemoveVdcFromVdcGroup{
				Name: generator.MustGenerate("{word}"),
				Vdcs: []types.ParamsCreateVdcGroupVdc{
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
			if tt.mockResponseStatus != 0 {
				endpoints.UpdateVdcGroup().CleanMockResponse()
				endpoints.UpdateVdcGroup().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListVdcGroupResponseStatus != 0 {
				endpoints.ListVdcGroup().CleanMockResponse()
				endpoints.ListVdcGroup().SetMockResponse(tt.mockListVdcGroupResponse, &tt.mockListVdcGroupResponseStatus)
			}

			client := newClient(t)

			err := client.RemoveVdcFromVdcGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestUpdateVdcGroup(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsUpdateVdcGroup

		mockResponse       any
		mockResponseStatus int

		mockListVdcGroupResponse       any
		mockListVdcGroupResponseStatus int

		expectedErr bool
	}{
		{
			name: "Update VDC Group",
			params: types.ParamsUpdateVdcGroup{
				ID:          generator.MustGenerate("{urn:vdcGroup}"),
				Description: utils.ToPTR("My updated VDC Group"),
			},
			expectedErr: false,
		},
		{
			name: "Update VDC Group with VDCGroup name",
			params: types.ParamsUpdateVdcGroup{
				Name:        "my-updated-vdc-group",
				Description: utils.ToPTR("My updated VDC Group"),
			},
			expectedErr: false,
		},
		{
			name: "Failed to update VDC Group",
			params: types.ParamsUpdateVdcGroup{
				ID:          generator.MustGenerate("{urn:vdcGroup}"),
				Name:        "my-updated-vdc-group",
				Description: utils.ToPTR("My updated VDC Group"),
			},
			mockResponseStatus: 400,
			expectedErr:        true,
		},
		{
			name: "Failed to list VDC Group",
			params: types.ParamsUpdateVdcGroup{
				ID:          generator.MustGenerate("{urn:vdcGroup}"),
				Name:        "my-updated-vdc-group",
				Description: utils.ToPTR("My updated VDC Group"),
			},
			mockListVdcGroupResponseStatus: 404,
			expectedErr:                    true,
		},
		{
			name: "List VDCGroup are empty",
			params: types.ParamsUpdateVdcGroup{
				ID:          generator.MustGenerate("{urn:vdcGroup}"),
				Name:        "my-updated-vdc-group",
				Description: utils.ToPTR("My updated VDC Group"),
			},
			mockListVdcGroupResponse:       itypes.ApiResponseListVdcGroup{Values: []itypes.ApiResponseListVdcGroupDetails{}},
			mockListVdcGroupResponseStatus: 200,
			expectedErr:                    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.UpdateVdcGroup().CleanMockResponse()
				endpoints.UpdateVdcGroup().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListVdcGroupResponseStatus != 0 {
				endpoints.ListVdcGroup().CleanMockResponse()
				endpoints.ListVdcGroup().SetMockResponse(tt.mockListVdcGroupResponse, &tt.mockListVdcGroupResponseStatus)
			}

			client := newClient(t)

			vdc, err := client.UpdateVdcGroup(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, vdc)
		})
	}
}
