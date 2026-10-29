/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func jsonResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cloud-Avenue-Mock", "true")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TestListOrgRoles(t *testing.T) {
	tests := []struct {
		name               string
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
		expectedLen        int
	}{
		{
			name:        "List Org Roles Success",
			expectedErr: false,
			expectedLen: 2,
		},
		{
			name:               "List Org Roles Error 500",
			mockResponseStatus: 500,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			ms.CleanResponse(endpoints.ListRoles())
			ms.SetResponseFunc(endpoints.ListRoles(), func(w http.ResponseWriter, r *http.Request) {
				if tt.mockResponseStatus != 0 {
					w.WriteHeader(tt.mockResponseStatus)
					return
				}
				assert.Equal(t, "application/json;version=39.1", r.Header.Get("Accept"))
				// Raw payload pins the CloudAPI pageable wire format (Page + values).
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Cloud-Avenue-Mock", "true")
				fmt.Fprint(w, `{"values":[{"id":"role-1","name":"Role 1","description":"Description 1"},{"id":"role-2","name":"Role 2","description":"Description 2"}],"resultTotal":2,"pageCount":1,"page":1,"pageSize":128}`)
			})

			result, err := client.ListOrgRoles(t.Context())
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Len(t, result, tt.expectedLen)
			}
		})
	}
}

func TestListOrgRolesPagination(t *testing.T) {
	client, ms := newClient(t)

	var requests []string
	ms.CleanResponse(endpoints.ListRoles())
	ms.SetResponseFunc(endpoints.ListRoles(), func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		requests = append(requests, page)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cloud-Avenue-Mock", "true")
		if page == "1" {
			fmt.Fprint(w, `{"values":[{"id":"role-1","name":"Role 1"}],"resultTotal":2,"pageCount":2,"page":1,"pageSize":1}`)
			return
		}
		fmt.Fprint(w, `{"values":[{"id":"role-2","name":"Role 2"}],"resultTotal":2,"pageCount":2,"page":2,"pageSize":1}`)
	})

	result, err := client.ListOrgRoles(t.Context())
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, []string{"1", "2"}, requests)
}

func TestGetOrgRole(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsGetOrgRole
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Get Org Role Success",
			params: types.ParamsGetOrgRole{
				ID: "role-1",
			},
			expectedErr: false,
		},
		{
			name: "Get Org Role Not Found",
			params: types.ParamsGetOrgRole{
				ID: "nonexistent",
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Get Org Role Validation Error",
			params: types.ParamsGetOrgRole{
				ID: "",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			ms.CleanResponse(endpoints.GetRole())
			ms.SetResponseFunc(endpoints.GetRole(), func(w http.ResponseWriter, r *http.Request) {
				if tt.mockResponseStatus != 0 {
					w.WriteHeader(tt.mockResponseStatus)
					return
				}
				resp := itypes.APIResponseRole{
					ID:          "role-1",
					Name:        "Role 1",
					Description: "Description 1",
				}
				jsonResponse(w, resp)
			})

			result, err := client.GetOrgRole(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, "role-1", result.ID)
			}
		})
	}
}

func TestCreateOrgRole(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.CreateRole())
	ms.SetResponseFunc(endpoints.CreateRole(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json;version=39.1", r.Header.Get("Accept"))
		resp := itypes.APIResponseRole{
			ID:          "role-new",
			Name:        "New Role",
			Description: "New Description",
		}
		jsonResponse(w, resp)
	})

	result, err := client.CreateOrgRole(t.Context(), types.ParamsCreateOrgRole{
		Name:        "New Role",
		Description: "New Description",
	})
	assert.Nil(t, err)
	assert.Equal(t, "role-new", result.ID)
}

func TestUpdateOrgRole(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.UpdateRole())
	ms.SetResponseFunc(endpoints.UpdateRole(), func(w http.ResponseWriter, r *http.Request) {
		resp := itypes.APIResponseRole{
			ID:          "role-1",
			Name:        "Updated Role",
			Description: "Updated Description",
		}
		jsonResponse(w, resp)
	})

	result, err := client.UpdateOrgRole(t.Context(), types.ParamsUpdateOrgRole{
		ID:          "role-1",
		Name:        "Updated Role",
		Description: "Updated Description",
	})
	assert.Nil(t, err)
	assert.Equal(t, "Updated Role", result.Name)
}

func TestDeleteOrgRole(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsDeleteOrgRole
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Delete Org Role Success",
			params: types.ParamsDeleteOrgRole{
				ID: "role-1",
			},
			expectedErr: false,
		},
		{
			name: "Delete Org Role Not Found",
			params: types.ParamsDeleteOrgRole{
				ID: "nonexistent",
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Delete Org Role Validation Error",
			params: types.ParamsDeleteOrgRole{
				ID: "",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			ms.CleanResponse(endpoints.DeleteRole())
			ms.SetResponseFunc(endpoints.DeleteRole(), func(w http.ResponseWriter, r *http.Request) {
				if tt.mockResponseStatus != 0 {
					w.WriteHeader(tt.mockResponseStatus)
				}
			})

			err := client.DeleteOrgRole(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestListOrgRoleRights(t *testing.T) {
	client, ms := newClient(t)

	var requests []string
	ms.CleanResponse(endpoints.ListRoleRights())
	ms.SetResponseFunc(endpoints.ListRoleRights(), func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		requests = append(requests, page)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cloud-Avenue-Mock", "true")
		if page == "1" {
			fmt.Fprint(w, `{"values":[{"id":"right-1","name":"Right 1"}],"resultTotal":2,"pageCount":2,"page":1,"pageSize":1}`)
			return
		}
		fmt.Fprint(w, `{"values":[{"id":"right-2","name":"Right 2"}],"resultTotal":2,"pageCount":2,"page":2,"pageSize":1}`)
	})

	result, err := client.ListOrgRoleRights(t.Context(), types.ParamsListOrgRoleRights{ID: "role-1"})
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, []string{"1", "2"}, requests)
}

func TestAddOrgRoleRights(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.AddRoleRights())
	ms.SetResponseFunc(endpoints.AddRoleRights(), func(w http.ResponseWriter, r *http.Request) {
		// Pin the request wire format: EntityReferences collection under "values".
		var body map[string][]itypes.APIEntityReference
		assert.Nil(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Contains(t, body, "values")

		resp := itypes.APIResponseRoleRights{
			Rights: []itypes.APIEntityReference{
				{ID: "right-1", Name: "Right 1"},
			},
		}
		jsonResponse(w, resp)
	})

	result, err := client.AddOrgRoleRights(t.Context(), types.ParamsAddOrgRoleRights{
		ID: "role-1",
		Rights: []types.ParamsEntityReference{
			{ID: "right-1", Name: "Right 1"},
		},
	})
	assert.Nil(t, err)
	assert.Len(t, result, 1)
}

func TestReplaceOrgRoleRights(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ReplaceRoleRights())
	ms.SetResponseFunc(endpoints.ReplaceRoleRights(), func(w http.ResponseWriter, r *http.Request) {
		resp := itypes.APIResponseRoleRights{
			Rights: []itypes.APIEntityReference{
				{ID: "right-2", Name: "Right 2"},
			},
		}
		jsonResponse(w, resp)
	})

	result, err := client.ReplaceOrgRoleRights(t.Context(), types.ParamsReplaceOrgRoleRights{
		ID: "role-1",
		Rights: []types.ParamsEntityReference{
			{ID: "right-2", Name: "Right 2"},
		},
	})
	assert.Nil(t, err)
	assert.Len(t, result, 1)
}

func TestRemoveAllOrgRoleRights(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ReplaceRoleRights())
	ms.SetResponseFunc(endpoints.ReplaceRoleRights(), func(w http.ResponseWriter, r *http.Request) {
		// Verify the PUT body carries an explicit empty values list.
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		assert.JSONEq(t, `{"values":[]}`, string(raw))

		resp := itypes.APIResponseRoleRights{
			Rights: []itypes.APIEntityReference{},
		}
		jsonResponse(w, resp)
	})

	err := client.RemoveAllOrgRoleRights(t.Context(), types.ParamsDeleteOrgRole{ID: "role-1"})
	assert.Nil(t, err)
}

func TestListRights(t *testing.T) {
	tests := []struct {
		name         string
		params       types.ParamsListRights
		expectedLen  int
		checkImplied bool
	}{
		{
			name:        "List Rights Success",
			params:      types.ParamsListRights{},
			expectedLen: 2,
		},
		{
			name: "List Rights with Filter",
			params: types.ParamsListRights{
				Filter: "name==TestRight",
			},
			// Mock returns 2 items regardless of filter - filter is passed to API
			// but we just verify the call succeeds
			expectedLen: 2,
		},
		{
			name:         "List Rights with Implied Rights",
			params:       types.ParamsListRights{},
			expectedLen:  1,
			checkImplied: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			ms.CleanResponse(endpoints.ListRights())
			ms.SetResponseFunc(endpoints.ListRights(), func(w http.ResponseWriter, r *http.Request) {
				if tt.params.Filter != "" {
					assert.Equal(t, tt.params.Filter, r.URL.Query().Get("filter"))
				}
				resp := itypes.APIResponseListRights{
					Rights: []itypes.APIResponseRight{
						{ID: "right-1", Name: "Right 1"},
						{ID: "right-2", Name: "Right 2"},
					},
					PageCount: 1,
				}
				if tt.checkImplied {
					resp = itypes.APIResponseListRights{
						Rights: []itypes.APIResponseRight{
							{
								ID:   "right-1",
								Name: "Right 1",
								ImpliedRights: []itypes.APIEntityReference{
									{ID: "implied-right-1", Name: "Implied Right 1"},
								},
							},
						},
						PageCount: 1,
					}
				}
				jsonResponse(w, resp)
			})

			result, err := client.ListRights(t.Context(), tt.params)
			assert.Nil(t, err)
			assert.Len(t, result, tt.expectedLen)

			if tt.checkImplied && len(result) > 0 {
				assert.Len(t, result[0].ImpliedRights, 1, "Right should have implied rights from payload")
				assert.Equal(t, "implied-right-1", result[0].ImpliedRights[0].ID)
			}
		})
	}
}

func TestListRightsPagination(t *testing.T) {
	client, ms := newClient(t)

	var requests []string
	ms.CleanResponse(endpoints.ListRights())
	ms.SetResponseFunc(endpoints.ListRights(), func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		requests = append(requests, page)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cloud-Avenue-Mock", "true")
		if page == "1" {
			fmt.Fprint(w, `{"values":[{"id":"right-1","name":"Right 1"}],"resultTotal":2,"pageCount":2,"page":1,"pageSize":1}`)
			return
		}
		fmt.Fprint(w, `{"values":[{"id":"right-2","name":"Right 2"}],"resultTotal":2,"pageCount":2,"page":2,"pageSize":1}`)
	})

	result, err := client.ListRights(t.Context(), types.ParamsListRights{})
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, []string{"1", "2"}, requests)
}

func TestGetRight(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsGetRight
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Get Right Success",
			params: types.ParamsGetRight{
				ID: "right-1",
			},
			expectedErr: false,
		},
		{
			name: "Get Right Not Found",
			params: types.ParamsGetRight{
				ID: "nonexistent",
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Get Right Validation Error",
			params: types.ParamsGetRight{
				ID: "",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			ms.CleanResponse(endpoints.GetRight())
			ms.SetResponseFunc(endpoints.GetRight(), func(w http.ResponseWriter, r *http.Request) {
				if tt.mockResponseStatus != 0 {
					w.WriteHeader(tt.mockResponseStatus)
					return
				}
				resp := itypes.APIResponseRight{
					ID:          "right-1",
					Name:        "Right 1",
					Description: "Description 1",
				}
				jsonResponse(w, resp)
			})

			result, err := client.GetRight(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, "right-1", result.ID)
			}
		})
	}
}
