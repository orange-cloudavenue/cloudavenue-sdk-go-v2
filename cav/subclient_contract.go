/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import "context"

// SubClient exposes backend-specific context data needed by request construction.
type SubClient interface {
	// ContextData returns backend metadata needed to build requests.
	ContextData(ctx context.Context) ContextData
}
