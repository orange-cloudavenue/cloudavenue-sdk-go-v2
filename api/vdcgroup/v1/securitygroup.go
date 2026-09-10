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
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/inetworkobjects"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListSecurityGroup   = "SecurityGroup.List"
	opGetSecurityGroup    = "SecurityGroup.Get"
	opCreateSecurityGroup = "SecurityGroup.Create"
	opUpdateSecurityGroup = "SecurityGroup.Update"
	opDeleteSecurityGroup = "SecurityGroup.Delete"
)

// getFirewallGroupWithRetry resolves a Firewall Group by ID or name, retrying up to 5
// times on "not found" errors to compensate for NSX-T eventual consistency shortly
// after a Firewall Group is created.
func getFirewallGroupWithRetry(ctx context.Context, c cav.Client, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error) {
	const maxAttempts = 5

	var lastErr error

	for attempt := range maxAttempts {
		if attempt > 0 {
			time.Sleep(2 * time.Second)
		}

		fwGroup, err := getFirewallGroup(ctx, c, idOrName, typeValue)
		if err == nil {
			return fwGroup, nil
		}

		var apiErr *errors.APIError
		if !stderrors.As(err, &apiErr) || !apiErr.IsNotFound() {
			return nil, err
		}

		lastErr = err
	}

	return nil, lastErr
}

// getFirewallGroup resolves a Firewall Group by ID (if idOrName is a firewallGroup URN)
// or by name+typeValue+ownerRef filtering against the List endpoint.
func getFirewallGroup(ctx context.Context, c cav.Client, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error) {
	if urn.IsSecurityGroup(idOrName) {
		ep := endpoints.GetFirewallGroup()

		resp, err := c.Do(
			ctx,
			ep,
			cav.WithPathParam(ep.PathParams[0], idOrName),
		)
		if err != nil {
			return nil, err
		}

		return resp.Result().(*itypes.ApiResponseFirewallGroup), nil
	}

	ep := endpoints.ListFirewallGroup()

	resp, err := c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("typeValue==%s;name==%s", typeValue, idOrName)),
	)
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.ApiResponseListFirewallGroup)
	if len(list.Values) == 0 {
		return nil, &errors.APIError{Operation: "GetFirewallGroup", StatusCode: 404, Message: fmt.Sprintf("firewall group %q not found", idOrName)}
	}
	if len(list.Values) > 1 {
		return nil, errors.Newf("multiple firewall groups found for %q", idOrName)
	}

	return &list.Values[0], nil
}

func resolveVdcGroupRef(ctx context.Context, c cav.Client, id, name string) (string, string, error) {
	ref, err := inetworkobjects.ResolveVdcGroupRef(ctx, id, name, func(ctx context.Context, id, name string) (inetworkobjects.VdcGroupRef, error) {
		ep := endpoints.ListVdcGroup()
		filter := "id==" + id
		if id == "" {
			filter = "name==" + name
		}

		resp, err := c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], filter))
		if err != nil {
			return inetworkobjects.VdcGroupRef{}, err
		}

		list := resp.Result().(*itypes.ApiResponseListVdcGroup)
		if len(list.Values) == 0 {
			return inetworkobjects.VdcGroupRef{}, errors.Newf("vdc group not found")
		}

		return inetworkobjects.VdcGroupRef{ID: list.Values[0].ID, Name: list.Values[0].Name}, nil
	})
	if err != nil {
		return "", "", err
	}

	return ref.ID, ref.Name, nil
}

func listFirewallGroupsByType(ctx context.Context, c cav.Client, vdcGroupID, vdcGroupName, typeValue string) (*types.ModelListFirewallGroup, error) {
	return inetworkobjects.ListFirewallGroupsByType(ctx, c, vdcGroupID, vdcGroupName, typeValue, func(ctx context.Context, id, name string) (inetworkobjects.VdcGroupRef, error) {
		ep := endpoints.ListVdcGroup()
		filter := "id==" + id
		if id == "" {
			filter = "name==" + name
		}

		resp, err := c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], filter))
		if err != nil {
			return inetworkobjects.VdcGroupRef{}, err
		}

		list := resp.Result().(*itypes.ApiResponseListVdcGroup)
		if len(list.Values) == 0 {
			return inetworkobjects.VdcGroupRef{}, errors.Newf("vdc group not found")
		}

		return inetworkobjects.VdcGroupRef{ID: list.Values[0].ID, Name: list.Values[0].Name}, nil
	})
}

// resolveFirewallGroupTarget resolves a firewall group by ID or name within a type.
func resolveFirewallGroupTarget(ctx context.Context, c cav.Client, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error) {
	return inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, typeValue, func(ctx context.Context, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error) {
		return getFirewallGroupWithRetry(ctx, c, idOrName, typeValue)
	})
}

func securityGroupMembersToRefs(members []types.ParamsFirewallGroupMember) []itypes.ApiObjectReference {
	refs := make([]itypes.ApiObjectReference, 0, len(members))
	for _, member := range members {
		refs = append(refs, itypes.ApiObjectReference{ID: member.ID, Name: member.Name})
	}
	return refs
}

func createSecurityGroupBody(ctx context.Context, c cav.Client, params types.ParamsCreateSecurityGroup) (itypes.ApiRequestFirewallGroup, error) {
	vdcGroupID, vdcGroupName, err := resolveVdcGroupRef(ctx, c, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return itypes.ApiRequestFirewallGroup{}, err
	}

	return itypes.ApiRequestFirewallGroup{
		Name:        params.Name,
		Description: params.Description,
		TypeValue:   itypes.FirewallGroupTypeSecurityGroup,
		Members:     securityGroupMembersToRefs(params.Members),
		OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil
}

// CreateSecurityGroup creates a security group for a VDC group.
func (c *Client) CreateSecurityGroup(ctx context.Context, params types.ParamsCreateSecurityGroup) (*types.ModelGetFirewallGroup, error) {
	body, err := createSecurityGroupBody(ctx, c.c, params)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateSecurityGroup, err)
	}

	ep := endpoints.CreateFirewallGroup()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateSecurityGroup, err)
	}

	created, ok := resp.Result().(*itypes.ApiResponseFirewallGroup)
	if !ok || created == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateSecurityGroup, resp.Result())
	}

	fwGroup, err := getFirewallGroupWithRetry(ctx, c.c, created.ID, itypes.FirewallGroupTypeSecurityGroup)
	if err != nil {
		return nil, fmt.Errorf("%s: extract: %w", opCreateSecurityGroup, err)
	}

	model := fwGroup.ToModel()
	return &model, nil
}

// ListSecurityGroup lists security groups owned by a VDC group.
func (c *Client) ListSecurityGroup(ctx context.Context, params types.ParamsListSecurityGroup) (*types.ModelListFirewallGroup, error) {
	model, err := listFirewallGroupsByType(ctx, c.c, params.VdcGroupID, params.VdcGroupName, itypes.FirewallGroupTypeSecurityGroup)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListSecurityGroup, err)
	}

	return model, nil
}

// GetSecurityGroup returns a security group by ID or name for a VDC group.
func (c *Client) GetSecurityGroup(ctx context.Context, params types.ParamsGetSecurityGroup) (*types.ModelGetFirewallGroup, error) {
	model, err := inetworkobjects.GetFirewallGroupModel(ctx, params.ID, params.Name, itypes.FirewallGroupTypeSecurityGroup, func(ctx context.Context, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error) {
		return getFirewallGroupWithRetry(ctx, c.c, idOrName, typeValue)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetSecurityGroup, err)
	}

	return model, nil
}

// UpdateSecurityGroup replaces security group fields for a VDC group.
func (c *Client) UpdateSecurityGroup(ctx context.Context, params types.ParamsUpdateSecurityGroup) (*types.ModelGetFirewallGroup, error) {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := resolveFirewallGroupTarget(ctx, c.c, idOrName, itypes.FirewallGroupTypeSecurityGroup)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateSecurityGroup, err)
	}

	description := current.Description
	if params.Description != "" {
		description = params.Description
	}

	members := current.Members
	if len(params.Members) != 0 {
		members = securityGroupMembersToRefs(params.Members)
	}

	body := itypes.ApiRequestFirewallGroup{
		ID:          current.ID,
		Name:        current.Name,
		Description: description,
		TypeValue:   itypes.FirewallGroupTypeSecurityGroup,
		Members:     members,
		OwnerRef:    current.OwnerRef,
	}

	if err := inetworkobjects.PutFirewallGroup(ctx, c.c, body); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateSecurityGroup, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteSecurityGroup deletes a security group from a VDC group.
func (c *Client) DeleteSecurityGroup(ctx context.Context, params types.ParamsDeleteSecurityGroup) error {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := resolveFirewallGroupTarget(ctx, c.c, idOrName, itypes.FirewallGroupTypeSecurityGroup)
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteSecurityGroup, err)
	}

	if err := inetworkobjects.DeleteFirewallGroup(ctx, c.c, current.ID); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteSecurityGroup, err)
	}

	return nil
}
