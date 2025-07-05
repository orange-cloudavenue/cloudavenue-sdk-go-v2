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

func init() {
	Endpoint{
		Name:             "SessionVmware",
		Method:           MethodPOST,
		SubClient:        ClientVmware,
		PathTemplate:     "/cloudapi/1.0.0/sessions",
		PathParams:       []PathParam{},
		QueryParams:      []QueryParam{},
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/1.0.0/sessions/post/",
		RequestFunc:      nil,
		requestInternalFunc: func(ctx context.Context, client *resty.Client, endpoint *Endpoint, opts ...EndpointRequestOption) (*resty.Response, error) {
			r := client.R().
				SetContext(ctx).
				SetHeader("Accept", "application/json;version="+vmwareVCDVersion)

			for _, opt := range opts {
				if err := opt(endpoint, r); err != nil {
					return nil, err
				}
			}

			return r.Post(endpoint.PathTemplate)
		},
	}.Register()
}
