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
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

func newClient(t *testing.T) (*Client, *mock.Server) {
	t.Helper()

	mC, ms, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating iam client")
	return eC, ms
}

func xmlResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("X-Cloud-Avenue-Mock", "true")
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TestListUsers(t *testing.T) {
	tests := []struct {
		name               string
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
		expectedLen        int
	}{
		{
			name:        "List Users Success",
			expectedErr: false,
			expectedLen: 2,
		},
		{
			name:               "List Users Error 500",
			mockResponseStatus: 500,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListUsers())
				ms.SetResponseFunc(endpoints.ListUsers(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.ListUsers())
				ms.SetResponseFunc(endpoints.ListUsers(), func(w http.ResponseWriter, r *http.Request) {
					users := itypes.Users{
						Users: []itypes.User{
							{Name: "user1", FullName: "User One", EmailAddress: "user1@example.com", IsEnabled: true, Role: itypes.Reference{Name: "Organization Administrator"}},
							{Name: "user2", FullName: "User Two", EmailAddress: "user2@example.com", IsEnabled: false, Role: itypes.Reference{Name: "vApp User"}},
						},
					}
					xmlResponse(w, users)
				})
			}

			result, err := client.ListUsers(t.Context())
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.expectedLen)
			assert.Equal(t, "user1", result[0].Name)
			assert.Equal(t, "User One", result[0].FullName)
			assert.Equal(t, "Organization Administrator", result[0].RoleName)
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsGetUser
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Get User by ID Success",
			params: ParamsGetUser{
				ID: generator.MustGenerate("{urn:user}"),
			},
			expectedErr: false,
		},
		{
			name: "Get User by Name Success",
			params: ParamsGetUser{
				Name: "user1",
			},
			expectedErr: false,
		},
		{
			name: "Get User Error 404",
			params: ParamsGetUser{
				ID: generator.MustGenerate("{urn:user}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name:        "Get User Missing ID and Name",
			params:      ParamsGetUser{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.GetUser())
				ms.SetResponseFunc(endpoints.GetUser(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.GetUser())
				ms.SetResponseFunc(endpoints.GetUser(), func(w http.ResponseWriter, r *http.Request) {
					user := itypes.User{
						Name:         "user1",
						FullName:     "User One",
						EmailAddress: "user1@example.com",
						IsEnabled:    true,
						Role:         itypes.Reference{Name: "Organization Administrator", Href: "https://example.com/role/123"},
					}
					xmlResponse(w, user)
				})
			}

			result, err := client.GetUser(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "user1", result.Name)
			assert.Equal(t, "Organization Administrator", result.RoleName)
		})
	}
}

func TestCreateLocalUser(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsCreateLocalUser
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Create Local User Success",
			params: ParamsCreateLocalUser{
				Name:            "newuser",
				Password:        "secret123",
				RoleName:        "Organization Administrator",
				FullName:        "New User",
				EmailAddress:    "newuser@example.com",
				IsEnabled:       true,
				DeployedVMQuota: 10,
				StoredVMQuota:   10,
			},
			expectedErr: false,
		},
		{
			name: "Create Local User Missing Name",
			params: ParamsCreateLocalUser{
				Password: "secret123",
				RoleName: "Organization Administrator",
			},
			expectedErr: true,
		},
		{
			name: "Create Local User Error 409",
			params: ParamsCreateLocalUser{
				Name:     "newuser",
				Password: "secret123",
				RoleName: "Organization Administrator",
			},
			mockResponseStatus: 409,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.CreateUser())
				ms.SetResponseFunc(endpoints.CreateUser(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.CreateUser())
				ms.SetResponseFunc(endpoints.CreateUser(), func(w http.ResponseWriter, r *http.Request) {
					user := itypes.User{
						Name:         tt.params.Name,
						FullName:     tt.params.FullName,
						EmailAddress: tt.params.EmailAddress,
						IsEnabled:    tt.params.IsEnabled,
						Role:         itypes.Reference{Name: tt.params.RoleName},
					}
					xmlResponse(w, user)
				})
			}

			result, err := client.CreateLocalUser(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.params.Name, result.Name)
			assert.Equal(t, tt.params.RoleName, result.RoleName)
		})
	}
}

func TestCreateSAMLUser(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsCreateSAMLUser
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Create SAML User Success",
			params: ParamsCreateSAMLUser{
				Name:         "samluser",
				RoleName:     "Organization Administrator",
				FullName:     "SAML User",
				EmailAddress: "samluser@example.com",
				IsEnabled:    true,
			},
			expectedErr: false,
		},
		{
			name: "Create SAML User Missing Name",
			params: ParamsCreateSAMLUser{
				RoleName: "Organization Administrator",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.CreateUser())
				ms.SetResponseFunc(endpoints.CreateUser(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.CreateUser())
				ms.SetResponseFunc(endpoints.CreateUser(), func(w http.ResponseWriter, r *http.Request) {
					user := itypes.User{
						Name:         tt.params.Name,
						FullName:     tt.params.FullName,
						EmailAddress: tt.params.EmailAddress,
						IsEnabled:    tt.params.IsEnabled,
						ProviderType: "SAML",
						Role:         itypes.Reference{Name: tt.params.RoleName},
					}
					xmlResponse(w, user)
				})
			}

			result, err := client.CreateSAMLUser(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.params.Name, result.Name)
			assert.Equal(t, "SAML", result.ProviderType)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsUpdateUser
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Update User Success",
			params: ParamsUpdateUser{
				ID:        generator.MustGenerate("{urn:user}"),
				FullName:  "Updated Name",
				Telephone: "1234567890",
			},
			expectedErr: false,
		},
		{
			name: "Update User Missing ID",
			params: ParamsUpdateUser{
				FullName: "Updated Name",
			},
			expectedErr: true,
		},
		{
			name: "Update User Error 404",
			params: ParamsUpdateUser{
				ID: generator.MustGenerate("{urn:user}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.UpdateUser())
				ms.SetResponseFunc(endpoints.UpdateUser(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.UpdateUser())
				ms.SetResponseFunc(endpoints.UpdateUser(), func(w http.ResponseWriter, r *http.Request) {
					user := itypes.User{
						Name:      "user1",
						FullName:  tt.params.FullName,
						Telephone: tt.params.Telephone,
						Role:      itypes.Reference{Name: "Organization Administrator"},
					}
					xmlResponse(w, user)
				})
			}

			result, err := client.UpdateUser(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.params.FullName, result.FullName)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsDeleteUser
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Delete User Success",
			params: ParamsDeleteUser{
				ID:            generator.MustGenerate("{urn:user}"),
				TakeOwnership: false,
			},
			expectedErr: false,
		},
		{
			name: "Delete User Missing ID",
			params: ParamsDeleteUser{
				TakeOwnership: false,
			},
			expectedErr: true,
		},
		{
			name: "Delete User Error 404",
			params: ParamsDeleteUser{
				ID:            generator.MustGenerate("{urn:user}"),
				TakeOwnership: false,
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.DeleteUser())
				ms.SetResponseFunc(endpoints.DeleteUser(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
				})
			}

			err := client.DeleteUser(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEnableUser(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsEnableUser
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Enable User Success",
			params: ParamsEnableUser{
				ID: generator.MustGenerate("{urn:user}"),
			},
			expectedErr: false,
		},
		{
			name: "Enable User by Name",
			params: ParamsEnableUser{
				Name: "user1",
			},
			expectedErr: false,
		},
		{
			name: "Enable User Error 404",
			params: ParamsEnableUser{
				ID: generator.MustGenerate("{urn:user}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.EnableUser())
				ms.SetResponseFunc(endpoints.EnableUser(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.EnableUser())
				ms.SetResponseFunc(endpoints.EnableUser(), func(w http.ResponseWriter, r *http.Request) {
					user := itypes.User{
						Name:      "user1",
						IsEnabled: true,
						Role:      itypes.Reference{Name: "Organization Administrator"},
					}
					xmlResponse(w, user)
				})
			}

			result, err := client.EnableUser(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, result.IsEnabled)
		})
	}
}

func TestDisableUser(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsDisableUser
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Disable User Success",
			params: ParamsDisableUser{
				ID: generator.MustGenerate("{urn:user}"),
			},
			expectedErr: false,
		},
		{
			name: "Disable User by Name",
			params: ParamsDisableUser{
				Name: "user1",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.DisableUser())
				ms.SetResponseFunc(endpoints.DisableUser(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.DisableUser())
				ms.SetResponseFunc(endpoints.DisableUser(), func(w http.ResponseWriter, r *http.Request) {
					user := itypes.User{
						Name:      "user1",
						IsEnabled: false,
						Role:      itypes.Reference{Name: "Organization Administrator"},
					}
					xmlResponse(w, user)
				})
			}

			result, err := client.DisableUser(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsEnabled)
		})
	}
}

func TestUnlockUser(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsUnlockUser
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "Unlock User Success",
			params: ParamsUnlockUser{
				ID: generator.MustGenerate("{urn:user}"),
			},
			expectedErr: false,
		},
		{
			name: "Unlock User by Name",
			params: ParamsUnlockUser{
				Name: "user1",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.UnlockUser())
				ms.SetResponseFunc(endpoints.UnlockUser(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.UnlockUser())
				ms.SetResponseFunc(endpoints.UnlockUser(), func(w http.ResponseWriter, r *http.Request) {
					user := itypes.User{
						Name:      "user1",
						IsEnabled: true,
						Role:      itypes.Reference{Name: "Organization Administrator"},
					}
					xmlResponse(w, user)
				})
			}

			result, err := client.UnlockUser(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, result.IsEnabled)
		})
	}
}

func TestChangePassword(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsChangePassword
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Change Password Success",
			params: ParamsChangePassword{
				ID:       generator.MustGenerate("{urn:user}"),
				Password: "newsecret123",
			},
			expectedErr: false,
		},
		{
			name: "Change Password by Name",
			params: ParamsChangePassword{
				Name:     "user1",
				Password: "newsecret123",
			},
			expectedErr: false,
		},
		{
			name: "Change Password Missing Password",
			params: ParamsChangePassword{
				ID: generator.MustGenerate("{urn:user}"),
			},
			expectedErr: true,
		},
		{
			name: "Change Password Error 404",
			params: ParamsChangePassword{
				ID:       generator.MustGenerate("{urn:user}"),
				Password: "newsecret123",
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ChangePassword())
				ms.SetResponseFunc(endpoints.ChangePassword(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
				})
			}

			err := client.ChangePassword(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestNewClientNil(t *testing.T) {
	c, err := New(nil)
	assert.Nil(t, c, "Expected nil client when input is nil")
	assert.Error(t, err, "Expected error when input is nil")
}
