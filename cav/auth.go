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

// auth defines credential lifecycle required by subclients.
type auth interface {
	// Headers returns headers required for authenticated requests.
	Headers() map[string]string

	// Refresh refreshes credential state.
	Refresh(context.Context) error

	// IsInitialized reports whether credentials are ready for use.
	IsInitialized() bool

	// getSession returns session data suitable for secure cache storage.
	getSession() map[string]string

	// restoreSession reloads session data from secure cache.
	restoreSession(data map[string]string) error

	// getExtraData returns backend metadata derived from session state.
	getExtraData() map[string]string
}
