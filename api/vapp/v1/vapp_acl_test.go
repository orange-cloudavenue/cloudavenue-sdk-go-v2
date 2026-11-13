/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	testVAppACLUserURN         = "urn:vcloud:orguser:22222222-2222-2222-2222-222222222222"
	testVAppACLGroupURN        = "urn:vcloud:orggroup:33333333-3333-3333-3333-333333333333"
	testVAppACLExternalUserID  = "ext-user-1"
	testVAppACLExternalGroupID = "ext-group-1"
	testVAppACLAccess          = "Change"
)

func TestVAppACLEndpoints(t *testing.T) {
	assert.Equal(t, "/api/vapp/{vapp-id}/controlAccess", endpoints.GetVAppACL().PathTemplate)
	assert.Equal(t, "/api/vapp/{vapp-id}/action/controlAccess", endpoints.SetVAppACL().PathTemplate)
}

func TestGetVAppACL(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetVAppACL())

	ms.CleanResponse(endpoints.GetVAppACL())
	ms.SetResponseFunc(endpoints.GetVAppACL(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.ControlAccessParams{
			Xmlns:               vcloudXMLNS,
			IsSharedToEveryone:  true,
			EveryoneAccessLevel: new("ReadOnly"),
			AccessSettings: &itypes.AccessSettingList{
				AccessSetting: []*itypes.AccessSetting{
					{
						ExternalSubject: &itypes.ExternalSubject{IsUser: true, SubjectID: testVAppACLExternalUserID},
						AccessLevel:     testVAppACLAccess,
					},
					{
						ExternalSubject: &itypes.ExternalSubject{IsUser: false, SubjectID: testVAppACLExternalGroupID},
						AccessLevel:     "ReadOnly",
					},
				},
			},
		})
	})

	resp, err := client.GetVAppACL(t.Context(), types.ParamsGetVAppACL{ID: testVAppID})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.True(t, resp.IsSharedToEveryone)
		assert.Equal(t, "ReadOnly", *resp.EveryoneAccessLevel)
		if assert.Len(t, resp.SharedWith, 2) {
			assert.Equal(t, "user", resp.SharedWith[0].SubjectKind)
			assert.Equal(t, types.VAppACLSubjectSourceExternal, resp.SharedWith[0].SubjectSource)
			assert.Equal(t, testVAppACLExternalUserID, resp.SharedWith[0].SubjectID)
			assert.Equal(t, "group", resp.SharedWith[1].SubjectKind)
			assert.Equal(t, types.VAppACLSubjectSourceExternal, resp.SharedWith[1].SubjectSource)
			assert.Equal(t, testVAppACLExternalGroupID, resp.SharedWith[1].SubjectID)
		}
	}
}

func TestGetVAppACLLocalHREFNormalization(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetVAppACL())

	ms.CleanResponse(endpoints.GetVAppACL())
	ms.SetResponseFunc(endpoints.GetVAppACL(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.ControlAccessParams{
			Xmlns: vcloudXMLNS,
			AccessSettings: &itypes.AccessSettingList{
				AccessSetting: []*itypes.AccessSetting{
					{
						Subject: &itypes.LocalSubject{
							HREF: "https://api.cloudavenue.org/api/admin/user/22222222-2222-2222-2222-222222222222",
							Name: "user1",
							Type: "urn:vcloud:orguser",
						},
						AccessLevel: testVAppACLAccess,
					},
					{
						Subject: &itypes.LocalSubject{
							HREF: "https://api.cloudavenue.org/api/admin/group/33333333-3333-3333-3333-333333333333",
							Name: "group1",
							Type: "urn:vcloud:orggroup",
						},
						AccessLevel: "ReadOnly",
					},
				},
			},
		})
	})

	resp, err := client.GetVAppACL(t.Context(), types.ParamsGetVAppACL{ID: testVAppID})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		if assert.Len(t, resp.SharedWith, 2) {
			assert.Equal(t, "user", resp.SharedWith[0].SubjectKind)
			assert.Equal(t, types.VAppACLSubjectSourceLocal, resp.SharedWith[0].SubjectSource)
			assert.Equal(t, testVAppACLUserURN, resp.SharedWith[0].SubjectID)
			assert.Equal(t, "user1", resp.SharedWith[0].SubjectName)
			assert.Equal(t, "group", resp.SharedWith[1].SubjectKind)
			assert.Equal(t, types.VAppACLSubjectSourceLocal, resp.SharedWith[1].SubjectSource)
			assert.Equal(t, testVAppACLGroupURN, resp.SharedWith[1].SubjectID)
			assert.Equal(t, "group1", resp.SharedWith[1].SubjectName)
		}
	}
}

func TestSetVAppACLLocalRoundTrip(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.SetVAppACL())
	apiEndpoint := client.c.GetConsole().GetAPIVCDEndpoint()

	ms.CleanResponse(endpoints.SetVAppACL())
	ms.SetResponseFunc(endpoints.SetVAppACL(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.ControlAccessParams
		assert.NoError(t, xml.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, vcloudXMLNS, body.Xmlns)
		assert.False(t, body.IsSharedToEveryone)
		if assert.NotNil(t, body.AccessSettings) && assert.Len(t, body.AccessSettings.AccessSetting, 2) {
			assert.Nil(t, body.AccessSettings.AccessSetting[0].ExternalSubject)
			assert.Equal(t, apiEndpoint+"/api/admin/user/22222222-2222-2222-2222-222222222222", body.AccessSettings.AccessSetting[0].Subject.HREF)
			assert.Equal(t, urnOrgUserType, body.AccessSettings.AccessSetting[0].Subject.Type)
			assert.Nil(t, body.AccessSettings.AccessSetting[1].ExternalSubject)
			assert.Equal(t, apiEndpoint+"/api/admin/group/33333333-3333-3333-3333-333333333333", body.AccessSettings.AccessSetting[1].Subject.HREF)
			assert.Equal(t, urnOrgGroupType, body.AccessSettings.AccessSetting[1].Subject.Type)
		}

		xmlResponse(w, itypes.ControlAccessParams{
			Xmlns:              body.Xmlns,
			IsSharedToEveryone: body.IsSharedToEveryone,
			AccessSettings: &itypes.AccessSettingList{AccessSetting: []*itypes.AccessSetting{
				{
					Subject:     body.AccessSettings.AccessSetting[0].Subject,
					AccessLevel: body.AccessSettings.AccessSetting[0].AccessLevel,
				},
				{
					Subject:     body.AccessSettings.AccessSetting[1].Subject,
					AccessLevel: body.AccessSettings.AccessSetting[1].AccessLevel,
				},
			}},
		})
	})

	resp, err := client.SetVAppACL(t.Context(), types.ParamsSetVAppACL{
		ID: testVAppID,
		SharedWith: []types.ParamsSetVAppACLSubject{
			{SubjectKind: "user", SubjectSource: types.VAppACLSubjectSourceLocal, SubjectID: testVAppACLUserURN, AccessLevel: testVAppACLAccess},
			{SubjectKind: "group", SubjectSource: types.VAppACLSubjectSourceLocal, SubjectID: testVAppACLGroupURN, AccessLevel: "ReadOnly"},
		},
	})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.SharedWith, 2) {
		assert.Equal(t, "user", resp.SharedWith[0].SubjectKind)
		assert.Equal(t, types.VAppACLSubjectSourceLocal, resp.SharedWith[0].SubjectSource)
		assert.Equal(t, testVAppACLUserURN, resp.SharedWith[0].SubjectID)
		assert.Equal(t, "group", resp.SharedWith[1].SubjectKind)
		assert.Equal(t, types.VAppACLSubjectSourceLocal, resp.SharedWith[1].SubjectSource)
		assert.Equal(t, testVAppACLGroupURN, resp.SharedWith[1].SubjectID)
	}
}

func TestSetVAppACLExternalSerialization(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.SetVAppACL())

	ms.CleanResponse(endpoints.SetVAppACL())
	ms.SetResponseFunc(endpoints.SetVAppACL(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.ControlAccessParams
		assert.NoError(t, xml.NewDecoder(r.Body).Decode(&body))
		if assert.NotNil(t, body.AccessSettings) && assert.Len(t, body.AccessSettings.AccessSetting, 2) {
			assert.Nil(t, body.AccessSettings.AccessSetting[0].Subject)
			assert.NotNil(t, body.AccessSettings.AccessSetting[0].ExternalSubject)
			assert.True(t, body.AccessSettings.AccessSetting[0].ExternalSubject.IsUser)
			assert.Equal(t, testVAppACLExternalUserID, body.AccessSettings.AccessSetting[0].ExternalSubject.SubjectID)
			assert.Nil(t, body.AccessSettings.AccessSetting[1].Subject)
			assert.NotNil(t, body.AccessSettings.AccessSetting[1].ExternalSubject)
			assert.False(t, body.AccessSettings.AccessSetting[1].ExternalSubject.IsUser)
			assert.Equal(t, testVAppACLExternalGroupID, body.AccessSettings.AccessSetting[1].ExternalSubject.SubjectID)
		}

		xmlResponse(w, body)
	})

	resp, err := client.SetVAppACL(t.Context(), types.ParamsSetVAppACL{
		ID: testVAppID,
		SharedWith: []types.ParamsSetVAppACLSubject{
			{SubjectKind: "user", SubjectSource: types.VAppACLSubjectSourceExternal, SubjectID: testVAppACLExternalUserID, AccessLevel: testVAppACLAccess},
			{SubjectKind: "group", SubjectSource: types.VAppACLSubjectSourceExternal, SubjectID: testVAppACLExternalGroupID, AccessLevel: "ReadOnly"},
		},
	})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.SharedWith, 2) {
		assert.Equal(t, types.VAppACLSubjectSourceExternal, resp.SharedWith[0].SubjectSource)
		assert.Equal(t, testVAppACLExternalUserID, resp.SharedWith[0].SubjectID)
		assert.Equal(t, types.VAppACLSubjectSourceExternal, resp.SharedWith[1].SubjectSource)
		assert.Equal(t, testVAppACLExternalGroupID, resp.SharedWith[1].SubjectID)
	}
}

func TestGetVAppACLMixedSubjectPreservation(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetVAppACL())

	ms.CleanResponse(endpoints.GetVAppACL())
	ms.SetResponseFunc(endpoints.GetVAppACL(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.ControlAccessParams{
			Xmlns: vcloudXMLNS,
			AccessSettings: &itypes.AccessSettingList{AccessSetting: []*itypes.AccessSetting{
				{
					Subject: &itypes.LocalSubject{
						HREF: "https://api.cloudavenue.org/api/admin/user/22222222-2222-2222-2222-222222222222",
						Name: "user1",
						Type: urnOrgUserType,
					},
					AccessLevel: testVAppACLAccess,
				},
				{
					ExternalSubject: &itypes.ExternalSubject{IsUser: false, SubjectID: testVAppACLExternalGroupID},
					AccessLevel:     "ReadOnly",
				},
			}},
		})
	})

	resp, err := client.GetVAppACL(t.Context(), types.ParamsGetVAppACL{ID: testVAppID})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.SharedWith, 2) {
		assert.Equal(t, types.VAppACLSubjectSourceLocal, resp.SharedWith[0].SubjectSource)
		assert.Equal(t, testVAppACLUserURN, resp.SharedWith[0].SubjectID)
		assert.Equal(t, types.VAppACLSubjectSourceExternal, resp.SharedWith[1].SubjectSource)
		assert.Equal(t, testVAppACLExternalGroupID, resp.SharedWith[1].SubjectID)
	}
}

func TestGetVAppACLExternalSubjectReadPreservation(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetVAppACL())

	ms.CleanResponse(endpoints.GetVAppACL())
	ms.SetResponseFunc(endpoints.GetVAppACL(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.ControlAccessParams{
			Xmlns: vcloudXMLNS,
			AccessSettings: &itypes.AccessSettingList{AccessSetting: []*itypes.AccessSetting{
				{
					ExternalSubject: &itypes.ExternalSubject{IsUser: true, SubjectID: testVAppACLExternalUserID},
					AccessLevel:     testVAppACLAccess,
				},
			}},
		})
	})

	resp, err := client.GetVAppACL(t.Context(), types.ParamsGetVAppACL{ID: testVAppID})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.SharedWith, 1) {
		assert.Equal(t, types.VAppACLSubjectSourceExternal, resp.SharedWith[0].SubjectSource)
		assert.Equal(t, "user", resp.SharedWith[0].SubjectKind)
		assert.Equal(t, testVAppACLExternalUserID, resp.SharedWith[0].SubjectID)
	}
}

func TestRemoveVAppACL(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.SetVAppACL())

	ms.CleanResponse(endpoints.SetVAppACL())
	ms.SetResponseFunc(endpoints.SetVAppACL(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.ControlAccessParams
		assert.NoError(t, xml.NewDecoder(r.Body).Decode(&body))
		assert.False(t, body.IsSharedToEveryone)
		assert.Nil(t, body.EveryoneAccessLevel)
		assert.Nil(t, body.AccessSettings)

		xmlResponse(w, body)
	})

	resp, err := client.RemoveVAppACL(t.Context(), types.ParamsRemoveVAppACL{ID: testVAppID})
	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.False(t, resp.IsSharedToEveryone)
		assert.Nil(t, resp.EveryoneAccessLevel)
		assert.Empty(t, resp.SharedWith)
	}
}

func TestVAppACLValidation(t *testing.T) {
	client, _ := newClient(t)

	_, err := client.GetVAppACL(t.Context(), types.ParamsGetVAppACL{})
	assert.Error(t, err)

	_, err = client.SetVAppACL(t.Context(), types.ParamsSetVAppACL{ID: testVAppID, SharedWith: []types.ParamsSetVAppACLSubject{{SubjectKind: "service-account", SubjectID: "x", AccessLevel: "ReadOnly"}}})
	assert.Error(t, err)

	_, err = client.SetVAppACL(t.Context(), types.ParamsSetVAppACL{ID: testVAppID, IsSharedToEveryone: true})
	assert.Error(t, err)

	readOnly := "ReadOnly"
	_, err = client.SetVAppACL(t.Context(), types.ParamsSetVAppACL{ID: testVAppID, IsSharedToEveryone: false, EveryoneAccessLevel: &readOnly})
	assert.Error(t, err)

	_, err = client.SetVAppACL(t.Context(), types.ParamsSetVAppACL{ID: testVAppID, SharedWith: []types.ParamsSetVAppACLSubject{{SubjectKind: "user", SubjectSource: types.VAppACLSubjectSourceLocal, SubjectID: testVAppACLExternalUserID, AccessLevel: "ReadOnly"}}})
	assert.Error(t, err)

	_, err = client.SetVAppACL(t.Context(), types.ParamsSetVAppACL{ID: testVAppID, SharedWith: []types.ParamsSetVAppACLSubject{{SubjectKind: "user", SubjectSource: types.VAppACLSubjectSourceExternal, SubjectID: testVAppACLExternalUserID, AccessLevel: "Owner"}}})
	assert.Error(t, err)

	_, err = client.SetVAppACL(t.Context(), types.ParamsSetVAppACL{ID: testVAppID, SharedWith: []types.ParamsSetVAppACLSubject{{SubjectKind: "group", SubjectSource: types.VAppACLSubjectSourceLocal, SubjectID: testVAppACLUserURN, AccessLevel: "ReadOnly"}}})
	assert.Error(t, err)

	_, err = client.RemoveVAppACL(t.Context(), types.ParamsRemoveVAppACL{})
	assert.Error(t, err)
}
