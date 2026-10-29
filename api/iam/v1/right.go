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
	"context"
	"fmt"
	"strconv"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListRights = "IAM.ListRights"
	opGetRight   = "IAM.GetRight"
)

// ListRights lists all rights, fetching every page.
func (c *Client) ListRights(ctx context.Context, params types.ParamsListRights) ([]*types.ModelRight, error) {
	var result []*types.ModelRight

	for page := 1; ; page++ {
		ep := endpoints.ListRights()

		opts := []cav.EndpointRequestOption{
			cav.WithQueryParam(endpointQueryParam(ep, "page"), strconv.Itoa(page)),
			cav.OverrideSetResult(new(itypes.APIResponseListRights)),
		}

		if params.Filter != "" {
			opts = append(opts, cav.WithQueryParam(endpointQueryParam(ep, "filter"), params.Filter))
		}

		resp, err := c.c.Do(ctx, ep, opts...)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", opListRights, err)
		}

		rights, ok := resp.Result().(*itypes.APIResponseListRights)
		if !ok {
			return nil, fmt.Errorf("%s: unexpected response type %T", opListRights, resp.Result())
		}

		result = append(result, rights.ToModel()...)
		if len(rights.Rights) == 0 || page >= rights.PageCount {
			break
		}
	}

	return result, nil
}

// GetRight retrieves a right by ID.
func (c *Client) GetRight(ctx context.Context, params types.ParamsGetRight) (*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetRight, err)
	}

	ep := endpoints.GetRight()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.OverrideSetResult(new(itypes.APIResponseRight)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetRight, err)
	}

	right, ok := resp.Result().(*itypes.APIResponseRight)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetRight, resp.Result())
	}

	return right.ToModel(), nil
}
