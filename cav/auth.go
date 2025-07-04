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
)

// Auth implements methods required for authentication.
type auth interface {
	// Headers returns headers that must be included in the http request.
	Headers() map[string]string

	// Refresh refreshes the authentication token.
	Refresh(context.Context) error

	// IsInitialized checks if the authentication is initialized.
	IsInitialized() bool
}
