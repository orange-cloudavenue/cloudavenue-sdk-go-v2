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
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

func newClient(t *testing.T) (*Client, *mock.Server) {
	t.Helper()

	mC, ms, err := mock.NewClient(
		mock.WithLogger(
			slog.New(
				slog.NewTextHandler(
					os.Stdout,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			),
		),
	)
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")
	return eC, ms
}

// setEdgeGatewayListResponse configures both the current Infrapi list endpoint
// and its compatibility alias. They share a route in the mock server.
func setEdgeGatewayListResponse(ms *mock.Server, data any, status *int) {
	ms.CleanResponse(endpoints.ListEdgeGateway())
	ms.SetResponse(endpoints.ListEdgeGateway(), data, status)
	ms.CleanResponse(endpoints.QueryEdgeGateway())
	ms.SetResponse(endpoints.QueryEdgeGateway(), data, status)
}
