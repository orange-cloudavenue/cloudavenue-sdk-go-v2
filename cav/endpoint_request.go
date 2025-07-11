/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"context"

	"resty.dev/v3"
)

var (
	defaultRequestFunc = func(ctx context.Context, client Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error) {
		req, err := client.NewRequest(ctx, endpoint)
		if err != nil {
			return nil, err
		}

		for _, opt := range opts {
			if err := opt(endpoint, req); err != nil {
				return nil, err
			}
		}

		return req.SetResult(endpoint.BodyResponseType).Execute(endpoint.Method.String(), endpoint.PathTemplate)
	}

	defaultRequestFuncWithJob = func(ctx context.Context, client Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error) {
		if endpoint.JobOptions == nil {
			endpoint.JobOptions = &JobOptions{}
		}

		return defaultRequestFunc(ctx, client, endpoint, opts...)
	}
)
