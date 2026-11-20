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
	"net/url"
	"slices"
	"strings"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	localSubjectTypeOrgUser  = "urn:vcloud:orguser"
	localSubjectTypeOrgGroup = "urn:vcloud:orggroup"
	// Media types of the local subjects: the API rejects any other value of the
	// Subject type attribute with HTTP 400 (go-vcloud-director sends the same).
	LocalSubjectMediaTypeUser    = "application/vnd.vmware.admin.user+xml"
	LocalSubjectMediaTypeGroup   = "application/vnd.vmware.admin.group+xml"
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
			if subject.SubjectKind == "" {
				// The API may answer with another type attribute: the href tells the kind.
				subject.SubjectKind = localSubjectKindFromHREF(as.Subject.HREF)
			}
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
	case LocalSubjectMediaTypeUser, localSubjectTypeOrgUser:
		return vappACLSubjectKindUser
	case LocalSubjectMediaTypeGroup, localSubjectTypeOrgGroup:
		return vappACLSubjectKindGroup
	default:
		return ""
	}
}

// localSubjectPath is a subject href path (a legacy admin one or a cloudapi one, the
// API answers with the latter) followed by the subject identifier.
type localSubjectPath struct {
	segments []string
	kind     string
	cloudapi bool
}

var localSubjectPaths = []localSubjectPath{
	{segments: []string{"api", "admin", "user"}, kind: vappACLSubjectKindUser},
	{segments: []string{"api", "admin", "group"}, kind: vappACLSubjectKindGroup},
	{segments: []string{"cloudapi", "1.0.0", "users"}, kind: vappACLSubjectKindUser, cloudapi: true},
	{segments: []string{"cloudapi", "1.0.0", "groups"}, kind: vappACLSubjectKindGroup, cloudapi: true},
}

// parseLocalSubjectHREF returns the kind and the orguser / orggroup URN of a subject href.
// Both are empty when the path of the href is not a known subject path.
func parseLocalSubjectHREF(href string) (kind, id string) {
	u, err := url.Parse(href)
	if err != nil {
		return "", ""
	}
	path := strings.Split(strings.Trim(u.Path, "/"), "/")

	for i := range path {
		for _, p := range localSubjectPaths {
			end := i + len(p.segments)
			if end >= len(path) || path[end] == "" || !slices.Equal(path[i:end], p.segments) {
				continue
			}

			prefix := localSubjectTypeOrgUser + ":"
			if p.kind == vappACLSubjectKindGroup {
				prefix = localSubjectTypeOrgGroup + ":"
			}

			// cloudapi hrefs end with the subject URN (urn:vcloud:user:<uuid>)
			identifier := path[end]
			if p.cloudapi {
				identifier = urn.ExtractUUID(path[end])
			}
			if identifier == "" {
				continue
			}
			return p.kind, prefix + identifier
		}
	}

	return "", ""
}

// localSubjectKindFromHREF returns the subject kind from the subject href.
func localSubjectKindFromHREF(href string) string {
	kind, _ := parseLocalSubjectHREF(href)
	return kind
}

// localSubjectID returns the orguser / orggroup URN of the subject href, or the href itself
// when it is not a known subject path.
func localSubjectID(href string) string {
	if _, id := parseLocalSubjectHREF(href); id != "" {
		return id
	}
	return href
}
