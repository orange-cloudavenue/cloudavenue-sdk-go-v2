/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalSubjectKind(t *testing.T) {
	tests := []struct {
		subjectType string
		want        string
	}{
		{"application/vnd.vmware.admin.user+xml", "user"},
		{"application/vnd.vmware.admin.group+xml", "group"},
		{"urn:vcloud:orguser", "user"},
		{"urn:vcloud:orggroup", "group"},
		{"", ""},
		{"application/other+xml", ""},
	}

	for _, tt := range tests {
		t.Run(tt.subjectType, func(t *testing.T) {
			assert.Equal(t, tt.want, localSubjectKind(tt.subjectType))
		})
	}
}

func TestLocalSubjectKindFromHREF(t *testing.T) {
	const uuid = "da4903bb-2f30-4dcd-af5b-45314b2aa193"

	tests := []struct {
		name string
		href string
		want string
	}{
		{"legacy user", "https://console.example.org/api/admin/user/" + uuid, "user"},
		{"legacy group", "https://console.example.org/api/admin/group/" + uuid, "group"},
		{"cloudapi user", "https://console.example.org/cloudapi/1.0.0/users/urn:vcloud:user:" + uuid, "user"},
		{"cloudapi group", "https://console.example.org/cloudapi/1.0.0/groups/urn:vcloud:group:" + uuid, "group"},
		{"unknown", "https://console.example.org/api/other/" + uuid, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, localSubjectKindFromHREF(tt.href))
		})
	}
}

func TestLocalSubjectID(t *testing.T) {
	const uuid = "da4903bb-2f30-4dcd-af5b-45314b2aa193"

	assert.Equal(t, "urn:vcloud:orguser:"+uuid, localSubjectID("https://console.example.org/api/admin/user/"+uuid))
	assert.Equal(t, "urn:vcloud:orggroup:"+uuid, localSubjectID("https://console.example.org/api/admin/group/"+uuid))
	assert.Equal(t, "urn:vcloud:orguser:"+uuid, localSubjectID("https://console.example.org/cloudapi/1.0.0/users/urn:vcloud:user:"+uuid))
	assert.Equal(t, "urn:vcloud:orggroup:"+uuid, localSubjectID("https://console.example.org/cloudapi/1.0.0/groups/urn:vcloud:group:"+uuid))
	assert.Empty(t, localSubjectID(""))
}

func TestToVAppACLModel_CloudapiSubjectWithUnknownType(t *testing.T) {
	const uuid = "da4903bb-2f30-4dcd-af5b-45314b2aa193"

	acl := ControlAccessParams{AccessSettings: &AccessSettingList{AccessSetting: []*AccessSetting{{
		Subject:     &LocalSubject{HREF: "https://console.example.org/cloudapi/1.0.0/users/urn:vcloud:user:" + uuid, Name: "example", Type: "application/json"},
		AccessLevel: "ReadOnly",
	}}}}

	model := acl.ToVAppACLModel()
	if assert.Len(t, model.SharedWith, 1) {
		assert.Equal(t, "user", model.SharedWith[0].SubjectKind)
		assert.Equal(t, "urn:vcloud:orguser:"+uuid, model.SharedWith[0].SubjectID)
		assert.Equal(t, "local", model.SharedWith[0].SubjectSource)
		assert.Equal(t, "example", model.SharedWith[0].SubjectName)
	}
}
