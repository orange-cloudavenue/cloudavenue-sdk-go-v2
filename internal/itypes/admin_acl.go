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
	"encoding/xml"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	localSubjectTypeOrgUser      = "urn:vcloud:orguser"
	localSubjectTypeOrgGroup     = "urn:vcloud:orggroup"
	vappACLSubjectKindUser       = "user"
	vappACLSubjectKindGroup      = "group"
	vappACLSubjectSourceLocal    = "local"
	vappACLSubjectSourceExternal = "external"
)

// ControlAccessParams represents the VMware vCD ControlAccessParams XML element.
type ControlAccessParams struct {
	XMLName             xml.Name           `xml:"ControlAccessParams"`
	Xmlns               string             `xml:"xmlns,attr"`
	IsSharedToEveryone  bool               `xml:"IsSharedToEveryone"`
	EveryoneAccessLevel *string            `xml:"EveryoneAccessLevel,omitempty"`
	AccessSettings      *AccessSettingList `xml:"AccessSettings,omitempty"`
}

// AccessSettingList represents the wrapper for a list of access settings.
type AccessSettingList struct {
	AccessSetting []*AccessSetting `xml:"AccessSetting"`
}

// LocalSubject represents a local subject (user or group) referenced in an access setting.
type LocalSubject struct {
	HREF string `xml:"href,attr"`
	Name string `xml:"name,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

// AccessSetting represents a single access setting entry.
type AccessSetting struct {
	XMLName         xml.Name         `xml:"AccessSetting"`
	Subject         *LocalSubject    `xml:"Subject,omitempty"`
	ExternalSubject *ExternalSubject `xml:"ExternalSubject,omitempty"`
	AccessLevel     string           `xml:"AccessLevel"`
}

// ExternalSubject represents an external (IdP) subject referenced in an access setting.
type ExternalSubject struct {
	IdpType   string `xml:"IdpType"`
	IsUser    bool   `xml:"IsUser"`
	SubjectID string `xml:"SubjectId"`
}

// ToModel converts the VMware vCD ControlAccessParams XML to the public model.
func (r *ControlAccessParams) ToModel() *types.ModelAdminCatalogACL {
	model := &types.ModelAdminCatalogACL{
		IsSharedToEveryone:  r.IsSharedToEveryone,
		EveryoneAccessLevel: r.EveryoneAccessLevel,
	}

	if r.AccessSettings == nil {
		return model
	}

	for _, as := range r.AccessSettings.AccessSetting {
		if as == nil {
			continue
		}

		subject := types.ModelAdminCatalogACLSubject{
			AccessLevel: as.AccessLevel,
		}

		if as.Subject != nil {
			subject.SubjectName = as.Subject.Name
		}

		if as.ExternalSubject != nil {
			if as.ExternalSubject.IsUser {
				subject.UserID = as.ExternalSubject.SubjectID
			} else {
				subject.GroupID = as.ExternalSubject.SubjectID
			}
		}

		model.SharedWith = append(model.SharedWith, subject)
	}

	return model
}

// ToVAppACLModel converts the VMware vCD ControlAccessParams XML to the vApp public model.
func (r *ControlAccessParams) ToVAppACLModel() *types.ModelVAppACL {
	model := &types.ModelVAppACL{
		IsSharedToEveryone:  r.IsSharedToEveryone,
		EveryoneAccessLevel: r.EveryoneAccessLevel,
	}

	if r.AccessSettings == nil {
		return model
	}

	for _, as := range r.AccessSettings.AccessSetting {
		if as == nil {
			continue
		}

		subject := types.ModelVAppACLSubject{
			AccessLevel: as.AccessLevel,
		}

		if as.Subject != nil {
			subject.SubjectName = as.Subject.Name
			subject.SubjectKind = localSubjectKind(as.Subject.Type)
			subject.SubjectID = localSubjectID(as.Subject.HREF)
			subject.SubjectSource = vappACLSubjectSourceLocal
		}

		if as.ExternalSubject != nil {
			if as.ExternalSubject.IsUser {
				subject.SubjectKind = vappACLSubjectKindUser
			} else {
				subject.SubjectKind = vappACLSubjectKindGroup
			}
			subject.SubjectID = as.ExternalSubject.SubjectID
			subject.SubjectSource = vappACLSubjectSourceExternal
		}

		model.SharedWith = append(model.SharedWith, subject)
	}

	return model
}

func localSubjectKind(subjectType string) string {
	switch subjectType {
	case localSubjectTypeOrgUser:
		return vappACLSubjectKindUser
	case localSubjectTypeOrgGroup:
		return vappACLSubjectKindGroup
	default:
		return ""
	}
}

func localSubjectID(href string) string {
	if href == "" {
		return ""
	}
	if id, ok := localSubjectURNFromHREF(href, "/api/admin/user/", "urn:vcloud:orguser:"); ok {
		return id
	}
	if id, ok := localSubjectURNFromHREF(href, "/api/admin/group/", "urn:vcloud:orggroup:"); ok {
		return id
	}
	return href
}

func localSubjectURNFromHREF(href, pathPrefix, urnPrefix string) (string, bool) {
	idx := len(href) - len(pathPrefix)
	if idx < 0 {
		return "", false
	}
	pos := -1
	for i := 0; i <= len(href)-len(pathPrefix); i++ {
		if href[i:i+len(pathPrefix)] == pathPrefix {
			pos = i
		}
	}
	if pos < 0 {
		return "", false
	}
	identifier := href[pos+len(pathPrefix):]
	if identifier == "" {
		return "", false
	}
	return urnPrefix + identifier, true
}
