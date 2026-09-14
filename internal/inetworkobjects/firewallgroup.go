/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package inetworkobjects

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

type FindFirewallGroupFunc func(ctx context.Context, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error)

func ListFirewallGroupsByType(ctx context.Context, c cav.Client, vdcGroupID, vdcGroupName, typeValue string, resolve ResolveVdcGroupFunc) (*types.ModelListFirewallGroup, error) {
	if vdcGroupID == "" {
		ref, err := resolve(ctx, "", vdcGroupName)
		if err != nil {
			return nil, err
		}
		vdcGroupID = ref.ID
	}

	ep := endpoints.ListFirewallGroup()
	resp, err := c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("typeValue==%s;ownerRef.id==%s", typeValue, vdcGroupID)),
	)
	if err != nil {
		return nil, err
	}

	return resp.Result().(*itypes.ApiResponseListFirewallGroup).ToModel(), nil
}

func GetFirewallGroupModel(ctx context.Context, id, name, typeValue string, find FindFirewallGroupFunc) (*types.ModelGetFirewallGroup, error) {
	idOrName := id
	if idOrName == "" {
		idOrName = name
	}

	fwGroup, err := find(ctx, idOrName, typeValue)
	if err != nil {
		return nil, err
	}

	model := fwGroup.ToModel()
	return &model, nil
}

func ResolveFirewallGroupTarget(ctx context.Context, idOrName, typeValue string, find FindFirewallGroupFunc) (*itypes.ApiResponseFirewallGroup, error) {
	if idOrName == "" {
		return nil, &pkgerrors.APIError{Operation: "ResolveFirewallGroup", StatusCode: 400, Message: "id or name is required"}
	}

	return find(ctx, idOrName, typeValue)
}

func PutFirewallGroup(ctx context.Context, c cav.Client, body itypes.ApiRequestFirewallGroup) error {
	ep := endpoints.UpdateFirewallGroup()
	_, err := c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], body.ID),
		cav.SetBody(body),
	)
	return err
}

func DeleteFirewallGroup(ctx context.Context, c cav.Client, id string) error {
	ep := endpoints.DeleteFirewallGroup()
	_, err := c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], id),
	)
	return err
}
