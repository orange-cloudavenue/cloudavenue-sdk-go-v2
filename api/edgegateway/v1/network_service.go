/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// UpdateNetworkService updates CAV standard service properties through Cerberus.
func (c *Client) UpdateNetworkService(ctx context.Context, params types.ParamsUpdateNetworkService) (*cav.CerberusJobCreatedAPIResponse, error) {
	if params.ServiceID == "" {
		return nil, fmt.Errorf("service ID is required")
	}
	body := itypes.APIRequestNetworkServiceUpdate{Adm: params.Adm, Os: params.Os, Tf: params.Tf}
	resp, err := c.c.Do(ctx, endpoints.UpdateNetworkService(), cav.WithPathParam(endpoints.UpdateNetworkService().PathParams[0], params.ServiceID), cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("failed to update network service: %w", err)
	}
	result, ok := resp.Result().(*cav.CerberusJobCreatedAPIResponse)
	if !ok || result == nil {
		return nil, fmt.Errorf("unexpected network service response type %T", resp.Result())
	}
	return result, nil
}
