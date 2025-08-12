/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

func (c *Client) retrieveEdgeGatewayIDByName(ctx context.Context, name string) (string, error) {
	epQuery := endpoints.QueryEdgeGateway()

	respQuery, err := c.c.Do(
		ctx,
		epQuery,
		cav.WithQueryParam(epQuery.QueryParams[1], "name=="+name),
	)
	if err != nil {
		return "", err
	}

	// Record is already checked in the middleware.
	return respQuery.Result().(*apiResponseQueryEdgeGateway).Record[0].ID, nil
}
