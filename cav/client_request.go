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

// NewRawRequest creates raw request for backend.
func (c *client) NewRawRequest(ctx context.Context, backend BackendTarget) (req *resty.Request, err error) {
	sc, err := c.identifyClient(ctx, backend)
	if err != nil {
		return nil, err
	}

	clientName, _ := backendToSubClientName(backend)
	ctxv := context.WithValue(ctx, contextKeyClientName, clientName)
	hC, err := sc.newHTTPClient(ctxv)
	if err != nil {
		return nil, err
	}

	hR := hC.NewRequest().
		SetContext(ctxv)

	return hR, nil
}
