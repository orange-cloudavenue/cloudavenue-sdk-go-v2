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
	"context"
	"fmt"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetVAppACL                 = "VApp.ACL.Get"
	opSetVAppACL                 = "VApp.ACL.Set"
	opRemoveVAppACL              = "VApp.ACL.Remove"
	pathAdminUserByID            = "/api/admin/user/%s"
	pathAdminGroupByID           = "/api/admin/group/%s"
	urnOrgUserPrefix             = "urn:vcloud:orguser:"
	urnOrgGroupPrefix            = "urn:vcloud:orggroup:"
	urnOrgUserType               = "urn:vcloud:orguser"
	urnOrgGroupType              = "urn:vcloud:orggroup"
	vappACLSubjectKindUser       = "user"
	vappACLSubjectKindGroup      = "group"
	vappACLSubjectSourceLocal    = "local"
	vappACLSubjectSourceExternal = "external"
)

var (
	getVAppACLOp = cav.Operation[getVAppByIDParams, *itypes.ControlAccessParams]{
		Name:     opGetVAppACL,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.GetVAppACL(),
		Validate: func(p getVAppByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p getVAppByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.GetVAppACL()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getVAppByIDParams) (*itypes.ControlAccessParams, error) {
			acl, ok := resp.Result().(*itypes.ControlAccessParams)
			if !ok || acl == nil {
				return nil, fmt.Errorf("unexpected get response type %T", resp.Result())
			}
			return acl, nil
		},
	}
	setVAppACLOp = cav.Operation[setVAppACLParams, *itypes.ControlAccessParams]{
		Name:     opSetVAppACL,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.SetVAppACL(),
		Validate: func(p setVAppACLParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			return nil
		},
		RequestOptions: func(p setVAppACLParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.SetVAppACL()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p setVAppACLParams) (any, error) {
			return p.Body, nil
		},
		Extract: func(resp *cav.Response, _ setVAppACLParams) (*itypes.ControlAccessParams, error) {
			acl, ok := resp.Result().(*itypes.ControlAccessParams)
			if !ok || acl == nil {
				return nil, fmt.Errorf("unexpected set response type %T", resp.Result())
			}
			return acl, nil
		},
	}
)

// GetVAppACL returns the full ACL state for a vApp.
func (c *Client) GetVAppACL(ctx context.Context, params types.ParamsGetVAppACL) (*types.ModelVAppACL, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetVAppACL, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppACL, err)
	}

	resp, err := cav.Execute(ctx, c.c, getVAppACLOp, getVAppByIDParams{ID: vappID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppACL, err)
	}

	return resp.ToVAppACLModel(), nil
}

// SetVAppACL replaces the full ACL state for a vApp.
func (c *Client) SetVAppACL(ctx context.Context, params types.ParamsSetVAppACL) (*types.ModelVAppACL, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opSetVAppACL, err)
	}

	vappID, err := c.resolveVAppID(ctx, params.ID, params.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opSetVAppACL, err)
	}

	resp, err := cav.Execute(ctx, c.c, setVAppACLOp, setVAppACLParams{ID: vappID, Body: toVAppControlAccessParams(c.c.GetConsole().GetAPIVCDEndpoint(), params)})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opSetVAppACL, err)
	}

	return resp.ToVAppACLModel(), nil
}

// RemoveVAppACL clears all explicit ACL grants on a vApp.
// It replaces the controlAccess state with IsSharedToEveryone=false
// and an empty subject list, matching govcd RemoveAccessControl semantics.
func (c *Client) RemoveVAppACL(ctx context.Context, params types.ParamsRemoveVAppACL) (*types.ModelVAppACL, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opRemoveVAppACL, err)
	}

	return c.SetVAppACL(ctx, types.ParamsSetVAppACL{ID: params.ID, Name: params.Name})
}

func toVAppControlAccessParams(apiEndpoint string, params types.ParamsSetVAppACL) itypes.ControlAccessParams {
	result := itypes.ControlAccessParams{
		Xmlns:               vcloudXMLNS,
		IsSharedToEveryone:  params.IsSharedToEveryone,
		EveryoneAccessLevel: params.EveryoneAccessLevel,
	}

	if len(params.SharedWith) == 0 {
		return result
	}

	settings := make([]*itypes.AccessSetting, 0, len(params.SharedWith))
	for _, subject := range params.SharedWith {
		setting := &itypes.AccessSetting{AccessLevel: subject.AccessLevel}
		switch subject.SubjectSource {
		case vappACLSubjectSourceLocal:
			setting.Subject = &itypes.LocalSubject{
				HREF: toVAppACLSubjectHREF(apiEndpoint, subject),
				Type: toVAppACLSubjectType(subject.SubjectKind),
			}
		case vappACLSubjectSourceExternal:
			setting.ExternalSubject = &itypes.ExternalSubject{
				IsUser:    subject.SubjectKind == vappACLSubjectKindUser,
				SubjectID: subject.SubjectID,
			}
		}
		settings = append(settings, setting)
	}

	result.AccessSettings = &itypes.AccessSettingList{AccessSetting: settings}
	return result
}

func toVAppACLSubjectHREF(apiEndpoint string, subject types.ParamsSetVAppACLSubject) string {
	if subject.SubjectSource != vappACLSubjectSourceLocal {
		return ""
	}
	switch subject.SubjectKind {
	case vappACLSubjectKindUser:
		return localSubjectHREFFromURN(apiEndpoint, subject.SubjectID, urnOrgUserPrefix, pathAdminUserByID)
	case vappACLSubjectKindGroup:
		return localSubjectHREFFromURN(apiEndpoint, subject.SubjectID, urnOrgGroupPrefix, pathAdminGroupByID)
	default:
		return subject.SubjectID
	}
}

func toVAppACLSubjectType(subjectKind string) string {
	switch subjectKind {
	case vappACLSubjectKindUser:
		return urnOrgUserType
	case vappACLSubjectKindGroup:
		return urnOrgGroupType
	default:
		return ""
	}
}

func localSubjectHREFFromURN(apiEndpoint, subjectID, urnPrefix, pathTemplate string) string {
	if subjectID == "" {
		return ""
	}
	if after, ok := strings.CutPrefix(subjectID, urnPrefix); ok {
		return apiEndpoint + fmt.Sprintf(pathTemplate, after)
	}
	return subjectID
}
