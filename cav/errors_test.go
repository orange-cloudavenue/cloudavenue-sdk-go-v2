/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL-2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"testing"

	"github.com/stretchr/testify/require"

	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

func TestClassifyStatusCode(t *testing.T) {
	require.ErrorIs(t, classifyStatusCode(400), pkgerrors.ErrBadRequest)
	require.ErrorIs(t, classifyStatusCode(401), pkgerrors.ErrUnauthorized)
	require.ErrorIs(t, classifyStatusCode(403), pkgerrors.ErrForbidden)
	require.ErrorIs(t, classifyStatusCode(404), pkgerrors.ErrNotFound)
	require.ErrorIs(t, classifyStatusCode(405), pkgerrors.ErrMethodNotAllowed)
	require.ErrorIs(t, classifyStatusCode(408), pkgerrors.ErrRequestTimeout)
	require.ErrorIs(t, classifyStatusCode(409), pkgerrors.ErrConflict)
	require.ErrorIs(t, classifyStatusCode(429), pkgerrors.ErrTooManyRequests)
	require.ErrorIs(t, classifyStatusCode(500), pkgerrors.ErrInternalServerError)
	require.ErrorIs(t, classifyStatusCode(502), pkgerrors.ErrBadGateway)
	require.ErrorIs(t, classifyStatusCode(503), pkgerrors.ErrServiceUnavailable)
	require.ErrorIs(t, classifyStatusCode(504), pkgerrors.ErrGatewayTimeout)
	require.NoError(t, classifyStatusCode(200))
}
